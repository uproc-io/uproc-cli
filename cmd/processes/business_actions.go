package processes

import "github.com/spf13/cobra"

func runSingleIDAction(cmd *cobra.Command, moduleSlug, action, idField, rawID string) error {
	id, err := parsePositiveIntArg(idField, rawID)
	if err != nil {
		return err
	}
	return runModuleAction(cmd, moduleSlug, action, map[string]any{idField: id})
}
