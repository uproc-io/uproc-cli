package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"bizzmod-cli/internal/api"
	"bizzmod-cli/internal/config"
	"github.com/spf13/cobra"
)

func mustMeClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if err := config.Validate(cfg); err != nil {
		return nil, err
	}
	return api.NewClient(cfg), nil
}

func printJSON(cmd *cobra.Command, body []byte) error {
	var parsed any
	if err := json.Unmarshal(body, &parsed); err != nil {
		fmt.Fprintln(cmd.OutOrStdout(), string(body))
		return nil
	}
	pretty, _ := json.MarshalIndent(parsed, "", "  ")
	fmt.Fprintln(cmd.OutOrStdout(), string(pretty))
	return nil
}

func newMeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "me",
		Short: "Read or update the connected user's profile",
	}
	cmd.AddCommand(newMeGetCmd())
	cmd.AddCommand(newMeUpdateCmd())
	return cmd
}

func newMeGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Show the connected user's profile",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustMeClient()
			if err != nil {
				return err
			}
			body, _, err := client.Do("GET", "/api/v1/external/profile", nil)
			if err != nil {
				return err
			}
			return printJSON(cmd, body)
		},
	}
}

func newMeUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [field=value ...]",
		Short: "Update profile fields (first_name, last_name, email, signature, language, default_app)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			payload := map[string]any{}
			for _, kv := range args {
				idx := strings.Index(kv, "=")
				if idx <= 0 {
					return fmt.Errorf("expected field=value, got %q", kv)
				}
				payload[kv[:idx]] = kv[idx+1:]
			}
			client, err := mustMeClient()
			if err != nil {
				return err
			}
			body, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("cannot encode profile payload: %w", err)
			}
			respBody, _, err := client.Do("PUT", "/api/v1/external/profile", body)
			if err != nil {
				return err
			}
			return printJSON(cmd, respBody)
		},
	}
}
