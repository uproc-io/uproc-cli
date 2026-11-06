package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newOperationsCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "operations",
		Short:  "Operations workflows (under construction)",
		Hidden: true,
		Long:   "Operations command group is under construction. Dedicated operations workflows will be added here. See https://uproc.io for updates.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "operations is under construction")
			fmt.Fprintln(cmd.OutOrStdout(), "Use `uproc operations --help` to check updates as new commands are added.")
			fmt.Fprintln(cmd.OutOrStdout(), "More info: https://uproc.io")
		},
	}
}
