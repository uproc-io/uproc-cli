package admin

import "github.com/spf13/cobra"

// NewAdminCmd groups administrative operations that require an elevated role.
func NewAdminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "Administrative operations (license management)",
	}
	cmd.AddCommand(NewLicenseCmd())
	return cmd
}
