package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// profileField is one row in the friendly `uproc me get` output.
type profileField struct {
	label string
	value string
}

func printFriendlyProfile(cmd *cobra.Command, body []byte) error {
	var root struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &root); err != nil {
		// Not JSON: show it as-is rather than hiding the server response.
		fmt.Fprintln(cmd.OutOrStdout(), strings.TrimSpace(string(body)))
		return nil
	}
	if !root.Success || root.Data == nil {
		message := strings.TrimSpace(root.Message)
		if message == "" {
			message = "Could not load the profile."
		}
		return fmt.Errorf("%s", message)
	}

	data := root.Data
	fields := []profileField{
		{"Name", strings.TrimSpace(asString(data["first_name"]) + " " + asString(data["last_name"]))},
		{"Email", asString(data["email"])},
		{"Role", asString(data["role"])},
		{"Language", asString(data["language"])},
	}

	if customer := asString(data["customer"]); customer != "" {
		label := customer
		if id := asString(data["customer_id"]); id != "" {
			label = fmt.Sprintf("%s (#%s)", customer, id)
		}
		fields = append(fields, profileField{"Customer", label})
	}
	if app := asString(data["default_app"]); app != "" {
		fields = append(fields, profileField{"Default app", app})
	}
	if modules, ok := data["modules"].([]interface{}); ok && len(modules) > 0 {
		names := make([]string, 0, len(modules))
		for _, module := range modules {
			if name := strings.TrimSpace(asString(module)); name != "" {
				names = append(names, name)
			}
		}
		sort.Strings(names)
		fields = append(fields, profileField{
			fmt.Sprintf("Modules (%d)", len(names)),
			strings.Join(names, ", "),
		})
	}
	if signature := strings.TrimSpace(asString(data["signature"])); signature != "" {
		fields = append(fields, profileField{"Signature", signature})
	}

	out := cmd.OutOrStdout()
	fmt.Fprintln(out, "Connected profile")
	fmt.Fprintln(out, "-----------------")
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	for _, field := range fields {
		fmt.Fprintf(tw, "%s:\t%s\n", field.label, field.value)
	}
	_ = tw.Flush()
	return nil
}

func asString(value interface{}) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case bool:
		if typed {
			return "yes"
		}
		return "no"
	case float64:
		if typed == float64(int64(typed)) {
			return fmt.Sprintf("%d", int64(typed))
		}
		return fmt.Sprintf("%v", typed)
	default:
		return fmt.Sprintf("%v", typed)
	}
}
