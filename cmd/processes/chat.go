package processes

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDataChatbotCmd() *cobra.Command {
	chatCmd := &cobra.Command{
		Use:   "data-chatbot",
		Short: "Business verbs for data chatbot workflows",
	}

	chatCmd.AddCommand(newChatAskCmd())
	chatCmd.AddCommand(newChatFollowUpCmd())
	chatCmd.AddCommand(newChatDomainsCmd())
	chatCmd.AddCommand(newChatInteractiveCmd())
	chatCmd.AddCommand(newCollectionListCmd("list", "List chatbot queries", "data-chatbot", "queries"))

	return chatCmd
}

func newChatCmd() *cobra.Command {
	cmd := newDataChatbotCmd()
	cmd.Use = "chat"
	cmd.Hidden = true
	cmd.Long = "DEPRECATED: use 'data-chatbot' instead"
	return cmd
}

func newChatAskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ask <domain> <question> [context] [channel] [sender_id] [origin_session_id]",
		Short: "Send a natural-language query to the data chatbot",
		Args:  cobra.RangeArgs(2, 6),
		RunE: func(cmd *cobra.Command, args []string) error {
			payload := map[string]any{
				"domain":   args[0],
				"question": args[1],
			}
			if len(args) >= 3 {
				payload["context"] = args[2]
			}
			if len(args) >= 4 {
				payload["channel"] = args[3]
			}
			if len(args) >= 5 {
				payload["sender_id"] = args[4]
			}
			if len(args) == 6 {
				payload["origin_session_id"] = args[5]
			}
			return runModuleAction(cmd, "data-chatbot", "send_chat_query", payload)
		},
	}
}

func newChatFollowUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "follow-up <origin_session_id> <question> [channel] [domain] [origin_user_id]",
		Short: "Continue an existing data-chatbot conversation",
		Args:  cobra.RangeArgs(2, 5),
		RunE: func(cmd *cobra.Command, args []string) error {
			payload := map[string]any{
				"origin_session_id": args[0],
				"question":          args[1],
			}
			if len(args) >= 3 {
				payload["channel"] = args[2]
			}
			if len(args) >= 4 {
				payload["domain"] = args[3]
			}
			if len(args) == 5 {
				payload["origin_user_id"] = args[4]
			}
			return runModuleAction(cmd, "data-chatbot", "follow_up", payload)
		},
	}
}

func newChatDomainsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "domains",
		Short: "List the available data-chatbot domains (channels) and their sample questions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			domains, err := fetchChatDomains(cmd)
			if err != nil {
				return err
			}
			if len(domains) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No data domains available.")
				return nil
			}
			for _, domain := range domains {
				fmt.Fprintf(cmd.OutOrStdout(), "Domain: %s\n", domain.Name)
				if len(domain.SampleQueries) == 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "  (no sample questions)")
				}
				for _, sample := range domain.SampleQueries {
					label := sample.Text
					if label == "" {
						label = sample.En
					}
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", label)
				}
			}
			return nil
		},
	}
}

func newChatInteractiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "interactive",
		Short: "Guided flow: pick a domain, select/edit a question, send it and follow up",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runChatInteractive(cmd)
		},
	}
}
