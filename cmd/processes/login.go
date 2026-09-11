package processes

import (
	"bufio"
	"fmt"
	"strings"

	"bizzmod-cli/internal/config"

	"github.com/spf13/cobra"
)

func newLoginCmd() *cobra.Command {
	var profileName string
	var useProfile bool

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Configure credentials with step-by-step profile wizard",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(profileName) == "" {
				active, err := config.GetActiveProfileName()
				if err == nil && strings.TrimSpace(active) != "" {
					profileName = active
				} else {
					profileName = "default"
				}
			}

			currentCfg, err := config.LoadProfile(profileName)
			if err != nil {
				return err
			}
			cfg := currentCfg

			if len(args) > 0 {
				if len(args) != 2 {
					return fmt.Errorf("when using arguments, provide API URL and user API key")
				}
				cfg = config.Config{
					APIURL:     args[0],
					UserAPIKey: args[1],
				}
			} else {
				interactiveCfg, err := promptReviewFields(cfg, cmd)
				if err != nil {
					return err
				}
				cfg = interactiveCfg
			}

			changed := hasConfigChanged(currentCfg, cfg)
			if changed {
				if err := validateExternalCredentials(cfg); err != nil {
					return err
				}
			}

			if err := config.SaveProfile(profileName, cfg, useProfile); err != nil {
				return err
			}

			if changed {
				fmt.Fprintf(cmd.OutOrStdout(), "ok: profile %q saved and validated\n", profileName)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "ok: profile %q unchanged\n", profileName)
			}

			if useProfile {
				fmt.Fprintf(cmd.OutOrStdout(), "ok: active profile set to %q\n", profileName)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&profileName, "profile", "", "profile name to save")
	cmd.Flags().BoolVar(&useProfile, "use", false, "set profile as active after saving")

	return cmd
}

func promptReviewFields(cfg config.Config, cmd *cobra.Command) (config.Config, error) {
	reader := bufio.NewReader(cmd.InOrStdin())

	apiURL, err := promptWithDefault(reader, cmd, "UPROC_PROCESSES_API_URL", cfg.APIURL)
	if err != nil {
		return cfg, err
	}
	cfg.APIURL = apiURL

	apiKey, err := promptWithDefault(reader, cmd, "USER_API_KEY", cfg.UserAPIKey)
	if err != nil {
		return cfg, err
	}
	cfg.UserAPIKey = apiKey

	return cfg, nil
}

func promptWithDefault(reader *bufio.Reader, cmd *cobra.Command, field string, current string) (string, error) {
	current = strings.TrimSpace(current)
	for {
		if current == "" {
			fmt.Fprintf(cmd.OutOrStdout(), "%s: ", field)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "%s [%s]: ", field, current)
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		value := strings.TrimSpace(line)
		if value == "" {
			if current != "" {
				return current, nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), "value is required")
			continue
		}
		return value, nil
	}
}

func hasConfigChanged(before config.Config, after config.Config) bool {
	return strings.TrimRight(strings.TrimSpace(before.APIURL), "/") != strings.TrimRight(strings.TrimSpace(after.APIURL), "/") ||
		strings.TrimSpace(before.UserAPIKey) != strings.TrimSpace(after.UserAPIKey)
}

func validateExternalCredentials(cfg config.Config) error {
	client, err := mustClientFromConfig(cfg)
	if err != nil {
		return err
	}

	body, status, reqErr := client.Do("GET", "/api/v1/external/modules", nil)
	if reqErr != nil {
		details := strings.TrimSpace(formatStructuredOutput(body))
		if details != "" {
			return fmt.Errorf("credentials validation failed: %w\n%s", reqErr, details)
		}
		return fmt.Errorf("credentials validation failed: %w", reqErr)
	}

	if status < 200 || status >= 300 {
		trimmed := strings.TrimSpace(string(body))
		if trimmed != "" {
			return fmt.Errorf("credentials validation failed with http %d: %s", status, trimmed)
		}
		return fmt.Errorf("credentials validation failed with http %d", status)
	}

	return nil
}
