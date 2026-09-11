package cmd

import (
	"bufio"
	"fmt"
	"strings"

	"bizzmod-cli/internal/config"
	"github.com/spf13/cobra"
)

func newProfileCmd() *cobra.Command {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage local login profiles",
	}

	profileCmd.AddCommand(newProfileListCmd())
	profileCmd.AddCommand(newProfileAddCmd())
	profileCmd.AddCommand(newProfileUseCmd())
	profileCmd.AddCommand(newProfileShowCmd())

	return profileCmd
}

func newProfileAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <profile_name>",
		Short: "Interactively add a URL and user API key profile, then activate it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := strings.TrimSpace(args[0])
			current, err := config.LoadProfile(name)
			if err != nil {
				return err
			}
			reader := bufio.NewReader(cmd.InOrStdin())
			apiURL, err := promptProfileValue(reader, cmd, "API URL", current.APIURL)
			if err != nil {
				return err
			}
			apiKey, err := promptProfileValue(reader, cmd, "User API key", current.UserAPIKey)
			if err != nil {
				return err
			}
			cfg := config.Config{APIURL: apiURL, UserAPIKey: apiKey}
			if err := config.Validate(cfg); err != nil {
				return err
			}
			if err := config.SaveProfile(name, cfg, true); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "ok: profile %q saved and set as active\n", name)
			return nil
		},
	}
}

func promptProfileValue(reader *bufio.Reader, cmd *cobra.Command, label, current string) (string, error) {
	current = strings.TrimSpace(current)
	for {
		if current == "" {
			fmt.Fprintf(cmd.OutOrStdout(), "%s: ", label)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "%s [configured]: ", label)
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		value := strings.TrimSpace(line)
		if value != "" {
			return value, nil
		}
		if current != "" {
			return current, nil
		}
		fmt.Fprintln(cmd.OutOrStdout(), "value is required")
	}
}

func newProfileListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List profiles with API URL and API key status",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, active, err := config.ListProfiles()
			if err != nil {
				return err
			}

			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no profiles)")
				return nil
			}

			for _, name := range names {
				marker := " "
				if name == active {
					marker = "*"
				}
				cfg, err := config.LoadProfile(name)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s\t%s\tuser API key: %t\n", marker, name, cfg.APIURL, cfg.UserAPIKey != "")
			}

			return nil
		},
	}
}

func newProfileUseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use <profile_name>",
		Short: "Set active profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.SetActiveProfile(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "ok: active profile set to %q\n", args[0])
			return nil
		},
	}
}

func newProfileShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show active profile URL and user API key status",
		RunE: func(cmd *cobra.Command, args []string) error {
			active, err := config.GetActiveProfileName()
			if err != nil {
				return err
			}
			cfg, err := config.LoadProfile(active)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\napi_url: %s\nuser_api_key: configured=%t\n", active, cfg.APIURL, cfg.UserAPIKey != "")
			return nil
		},
	}
}
