package processes

import "github.com/spf13/cobra"

func newFileSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file-sync",
		Short: "Copy files between storage providers (SFTP / S3 / Dropbox / Google Drive)",
	}
	cmd.AddCommand(newCollectionListCmd("list-transfers", "List file sync transfers", "file-sync", "transfers"))
	cmd.AddCommand(newCollectionListCmd("list-runs", "List file sync runs", "file-sync", "runs"))
	cmd.AddCommand(&cobra.Command{
		Use:   "run <transfer_id>",
		Short: "Run a file sync transfer now (copy files with overwrite)",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id, err := parsePositiveIntArg("transfer_id", args[0])
			if err != nil {
				return err
			}
			return runModuleAction(c, "file-sync", "run_transfer", map[string]any{"transfer_id": id})
		},
	})
	return cmd
}
