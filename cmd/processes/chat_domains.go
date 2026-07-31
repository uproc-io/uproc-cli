package processes

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type chatSampleQuery struct {
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	Type        string `json:"type"`
	Text        string `json:"text"`
	En          string `json:"en"`
}

type chatDomain struct {
	ID            int               `json:"id"`
	Name          string            `json:"name"`
	SampleQueries []chatSampleQuery `json:"sample_queries"`
}

func fetchChatDomains(cmd *cobra.Command) ([]chatDomain, error) {
	client, err := mustClient()
	if err != nil {
		return nil, err
	}
	body, status, reqErr := client.Do("GET", "/api/v1/external/modules/data-chatbot/data_domains", nil)
	if reqErr != nil {
		return nil, printResponse(cmd, body, status, reqErr)
	}
	var envelope struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			DataDomains []chatDomain `json:"data_domains"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("cannot parse data_domains response: %w", err)
	}
	if !envelope.Success {
		return nil, fmt.Errorf("%s", envelope.Message)
	}
	return envelope.Data.DataDomains, nil
}

func callChatAction(cmd *cobra.Command, action string, payload map[string]any) (map[string]any, error) {
	client, err := mustClient()
	if err != nil {
		return nil, err
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot encode %s payload: %w", action, err)
	}
	path := fmt.Sprintf("/api/v1/external/modules/data-chatbot/actions/%s", action)
	respBody, status, reqErr := client.Do("POST", path, bodyBytes)
	if reqErr != nil {
		return nil, printResponse(cmd, respBody, status, reqErr)
	}
	var parsed map[string]any
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, printResponse(cmd, respBody, status, nil)
	}
	success, _ := parsed["success"].(bool)
	if !success {
		return nil, printResponse(cmd, respBody, status, nil)
	}
	return parsed, nil
}

func chatSend(cmd *cobra.Command, action string, payload map[string]any) (string, string) {
	parsed, err := callChatAction(cmd, action, payload)
	if err != nil {
		fmt.Fprintln(cmd.OutOrStdout(), err)
		return "", ""
	}
	data, _ := parsed["data"].(map[string]any)
	response, _ := data["response"].(map[string]any)
	message, _ := response["message"].(string)
	sessionID, _ := response["session_id"].(string)
	if message != "" {
		fmt.Fprintln(cmd.OutOrStdout(), message)
	}
	if sessionID != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "session: %s\n", sessionID)
	}
	return message, sessionID
}

func chatPrompt(reader *bufio.Reader, cmd *cobra.Command, prompt string) (string, error) {
	fmt.Fprintf(cmd.OutOrStdout(), "%s", prompt)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func runChatInteractive(cmd *cobra.Command) error {
	reader := bufio.NewReader(cmd.InOrStdin())

	domains, err := fetchChatDomains(cmd)
	if err != nil {
		return err
	}
	if len(domains) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No data domains available. Configure them in the data-chatbot setup first.")
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Available data domains (channels):")
	for i, domain := range domains {
		fmt.Fprintf(cmd.OutOrStdout(), "  [%d] %s\n", i+1, domain.Name)
	}

	selection, err := chatPrompt(reader, cmd, fmt.Sprintf("Select domain (1-%d), or 'q' to quit: ", len(domains)))
	if err != nil {
		return err
	}
	if strings.EqualFold(selection, "q") {
		return nil
	}
	domainNum, convErr := strconv.Atoi(selection)
	if convErr != nil || domainNum < 1 || domainNum > len(domains) {
		return fmt.Errorf("invalid domain selection")
	}
	domain := domains[domainNum-1]

	channel := "web"
	originUserID := ""
	sessionID := ""
	lastQuestion := ""

	for {
		fmt.Fprintf(cmd.OutOrStdout(), "\n[%s] Available questions:\n", domain.Name)
		if len(domain.SampleQueries) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "  (none — type your own question)")
		}
		for i, sample := range domain.SampleQueries {
			label := sample.Text
			if label == "" {
				label = sample.En
			}
			tags := []string{}
			if sample.Category != "" {
				tags = append(tags, sample.Category)
			}
			if sample.Type != "" {
				tags = append(tags, sample.Type)
			}
			suffix := ""
			if len(tags) > 0 {
				suffix = " (" + strings.Join(tags, " · ") + ")"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  [%d] %s%s\n", i+1, label, suffix)
		}

		selection, err := chatPrompt(reader, cmd, "Select a question, type your own, 'c' to change channel/domain, or 'q' to quit: ")
		if err != nil {
			return err
		}
		if strings.EqualFold(selection, "q") {
			return nil
		}
		if strings.EqualFold(selection, "c") {
			if err := runChatChannelSwitch(reader, cmd, &channel, &originUserID); err != nil {
				return err
			}
			continue
		}

		question := selection
		if num, numErr := strconv.Atoi(selection); numErr == nil && num >= 1 && num <= len(domain.SampleQueries) {
			question = domain.SampleQueries[num-1].Text
			if question == "" {
				question = domain.SampleQueries[num-1].En
			}
		}
		if question == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "Question cannot be empty.")
			continue
		}

		edited, err := promptWithDefault(reader, cmd, "Edit question (Enter to keep)", question)
		if err != nil {
			return err
		}
		question = edited

		action := "send_chat_query"
		payload := map[string]any{
			"domain":   domain.Name,
			"question": question,
			"channel":  channel,
		}
		if sessionID != "" {
			action = "follow_up"
			payload["origin_session_id"] = sessionID
			payload["context"] = lastQuestion
		}
		if originUserID != "" {
			payload["origin_user_id"] = originUserID
		}

		_, newSession := chatSend(cmd, action, payload)
		if newSession != "" {
			sessionID = newSession
		}
		lastQuestion = question
	}
}

func runChatChannelSwitch(reader *bufio.Reader, cmd *cobra.Command, channel *string, originUserID *string) error {
	newChannel, err := chatPrompt(reader, cmd, "New channel (e.g. web, telegram) or 'q': ")
	if err != nil {
		return err
	}
	if strings.EqualFold(newChannel, "q") || newChannel == "" {
		return nil
	}
	*channel = strings.ToLower(strings.TrimSpace(newChannel))
	if *channel != "web" && *originUserID == "" {
		oid, err := chatPrompt(reader, cmd, "Origin user ID for this channel: ")
		if err != nil {
			return err
		}
		*originUserID = strings.TrimSpace(oid)
	}
	return nil
}
