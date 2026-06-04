package processes

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	processesCmd := &cobra.Command{
		Use:   "processes",
		Short: "Commands for Uproc Processes API",
	}

	processesCmd.AddCommand(newLoginCmd())
	processesCmd.AddCommand(newRequestCmd())
	processesCmd.AddCommand(newModuleCmd())
	processesCmd.AddCommand(newFormsCmd())
	processesCmd.AddCommand(newCandidateCmd())
	processesCmd.AddCommand(newSupportCmd())
	processesCmd.AddCommand(newApprovalCmd())
	processesCmd.AddCommand(newAdminCmd())
	processesCmd.AddCommand(newInstallCmd())
	processesCmd.AddCommand(newUpdateCmd())
	processesCmd.AddCommand(newInteractiveCmd())

	return processesCmd
}
