package processes

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newDataProcessCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data-process",
		Short: "Business verbs for data-process workflows",
	}
	cmd.AddCommand(newDataProcessExecuteCmd())
	cmd.AddCommand(newDataProcessBatchCmd())
	cmd.AddCommand(newCollectionListCmd("runs", "List data-process runs", "data-process", "runs"))
	return cmd
}

func newDataProcessExecuteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "execute <tool> [params_json]",
		Short: "Execute a data-process tool inline (no daily inline-call limit; balance applies)",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			payload := map[string]any{"processor": args[0]}
			if len(args) == 2 && args[1] != "" {
				var params map[string]any
				if err := json.Unmarshal([]byte(args[1]), &params); err != nil {
					return fmt.Errorf("params_json must be valid JSON: %w", err)
				}
				payload["params"] = params
			}
			return runModuleAction(cmd, "data-process", "execute", payload)
		},
	}
}

func newDataProcessBatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch <tool> <source_entity_id> [column_mapping_json] [fixed_values_json] [output_name]",
		Short: "Queue a batch run of a data-process tool over a data-management entity",
		Args:  cobra.RangeArgs(2, 5),
		RunE: func(cmd *cobra.Command, args []string) error {
			var entityID int
			if _, err := fmt.Sscanf(args[1], "%d", &entityID); err != nil || entityID <= 0 {
				return fmt.Errorf("source_entity_id must be a positive integer")
			}
			payload := map[string]any{"tool": args[0], "source_entity_id": entityID}
			if len(args) >= 3 && args[2] != "" {
				var mapping map[string]any
				if err := json.Unmarshal([]byte(args[2]), &mapping); err != nil {
					return fmt.Errorf("column_mapping_json must be valid JSON: %w", err)
				}
				payload["column_mapping"] = mapping
			}
			if len(args) >= 4 && args[3] != "" {
				var fixed map[string]any
				if err := json.Unmarshal([]byte(args[3]), &fixed); err != nil {
					return fmt.Errorf("fixed_values_json must be valid JSON: %w", err)
				}
				payload["fixed_values"] = fixed
			}
			if len(args) == 5 {
				payload["output_name"] = args[4]
			}
			return runModuleAction(cmd, "data-process", "batch", payload)
		},
	}
}
