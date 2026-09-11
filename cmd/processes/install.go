package processes

import (
	"fmt"

	"bizzmod-cli/internal/config"
	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "install <user_api_key>",
		Short: "Show installation plan for a customer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			cfg.UserAPIKey = args[0]
			if err := config.Validate(cfg); err != nil {
				return err
			}

			client, err := mustClientFromConfig(cfg)
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/api/v1/external/install?dry_run=%t", dryRun)
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", true, "show every installation step without executing")

	return cmd
}
