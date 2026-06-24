package cmd

import (
	"bizzmod-cli/cmd/processes"
	"github.com/spf13/cobra"
)

func newDataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data",
		Short: "Data workflows",
	}

	cmd.AddCommand(processes.NewDataColumnCmd())
	return cmd
}
