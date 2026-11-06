package admin

import (
	"github.com/spf13/cobra"

	"bizzmod-cli/cmd/processes"
)

// NewAdminCmd groups administrative operations that require an elevated role.
//
// It mirrors the web `/admin` section as a top-level command group:
// `uproc admin users|customers|credentials|modules|tickets|logs|ai-requests|usage|audit|changelog`
// plus license management. The same resources remain available under
// `uproc applications admin ...` for backwards compatibility.
func NewAdminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "Administrative operations (users, customers, license, ...)",
	}

	cmd.AddCommand(processes.NewAdminChildren()...)
	cmd.AddCommand(processes.NewInstallCmd())
	cmd.AddCommand(processes.NewUpdateCmd())
	cmd.AddCommand(NewLicenseCmd())
	return cmd
}
