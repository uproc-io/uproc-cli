package cmd

import (
	"bizzmod-cli/cmd/processes"
	"bizzmod-cli/internal/config"
	"fmt"
	"github.com/spf13/cobra"
)

func NewRootCmd(version string) *cobra.Command {
	var configPath string

	rootCmd := &cobra.Command{
		Use:     "uproc",
		Short:   "Uproc CLI for managing and interacting with Uproc services (https://uproc.io)",
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			config.SetConfigPath(configPath)
			return nil
		},
	}

	rootCmd.SetVersionTemplate("{{printf \"uproc version %s\\n\" .Version}}")

	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd == rootCmd {
			config.SetConfigPath(configPath)
			if path, err := config.ResolvedConfigPath(); err == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Config file: %s\n\n", path)
			}
		}
		defaultHelp(cmd, args)
	})

	rootCmd.AddCommand(processes.NewCmd())
	rootCmd.AddCommand(processes.NewProcessesAliasCmd())
	rootCmd.AddCommand(newOperationsCmd())
	rootCmd.AddCommand(newDataCmd())
	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newProfileCmd())
	rootCmd.AddCommand(newMeCmd())
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = cmd.OutOrStdout().Write([]byte("uproc version " + version + "\n"))
		},
	})
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "path to config file (defaults to ./config.yml)")

	return rootCmd
}
