package processes

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newDMQualityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quality",
		Short: "Data-management data quality",
		Long: `Inspect and run the deterministic data-quality analysis of a table.

Examples:
  uproc applications data-management quality get --entity-slug clientes
  uproc applications data-management quality analyze --entity-id 123
  uproc applications data-management quality rules --entity-slug clientes --column-key email --accept regex:email
  uproc applications data-management quality correct --entity-slug clientes --column-key email --apply`,
	}
	cmd.AddCommand(newDMQualityGetCmd())
	cmd.AddCommand(newDMQualityAnalyzeCmd())
	cmd.AddCommand(newDMQualityRulesCmd())
	cmd.AddCommand(newDMQualityCorrectCmd())
	return cmd
}

func qualityEntityQuery(entityID int, entitySlug string) string {
	if entitySlug != "" {
		return "entity_slug=" + entitySlug
	}
	return fmt.Sprintf("entity_id=%d", entityID)
}

func newDMQualityGetCmd() *cobra.Command {
	var entityID int
	var entitySlug string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Show the latest quality run, findings and history for a table",
		RunE: func(cmd *cobra.Command, args []string) error {
			if entityID <= 0 && entitySlug == "" {
				return fmt.Errorf("--entity-id or --entity-slug is required")
			}
			client, err := mustClient()
			if err != nil {
				return err
			}
			response, status, reqErr := client.Do(
				"GET",
				"/api/v1/external/data-management/quality?"+qualityEntityQuery(entityID, entitySlug),
				nil,
			)
			return printResponse(cmd, response, status, reqErr)
		},
	}
	cmd.Flags().IntVar(&entityID, "entity-id", 0, "Table ID")
	cmd.Flags().StringVar(&entitySlug, "entity-slug", "", "Table slug (alternative to --entity-id)")
	return cmd
}

func newDMQualityAnalyzeCmd() *cobra.Command {
	var entityID int
	var entitySlug string
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Queue a quality analysis for a table",
		RunE: func(cmd *cobra.Command, args []string) error {
			if entityID <= 0 && entitySlug == "" {
				return fmt.Errorf("--entity-id or --entity-slug is required")
			}
			client, err := mustClient()
			if err != nil {
				return err
			}
			body, _ := json.Marshal(map[string]any{"entity_id": entityID, "entity_slug": entitySlug})
			response, status, reqErr := client.Do(
				"POST",
				"/api/v1/external/data-management/quality/analyze",
				body,
			)
			return printResponse(cmd, response, status, reqErr)
		},
	}
	cmd.Flags().IntVar(&entityID, "entity-id", 0, "Table ID")
	cmd.Flags().StringVar(&entitySlug, "entity-slug", "", "Table slug (alternative to --entity-id)")
	return cmd
}

func newDMQualityRulesCmd() *cobra.Command {
	var entityID, columnID int
	var entitySlug, columnKey, accept, reject string
	cmd := &cobra.Command{
		Use:   "rules",
		Short: "Accept or reject inferred quality rules for a column",
		RunE: func(cmd *cobra.Command, args []string) error {
			if entityID <= 0 && entitySlug == "" {
				return fmt.Errorf("--entity-id or --entity-slug is required")
			}
			if columnID <= 0 && columnKey == "" {
				return fmt.Errorf("--column-id or --column-key is required")
			}
			client, err := mustClient()
			if err != nil {
				return err
			}
			body, _ := json.Marshal(map[string]any{
				"entity_id":   entityID,
				"entity_slug": entitySlug,
				"column_id":   columnID,
				"column_key":  columnKey,
				"accept":      splitCSV(accept),
				"reject":      splitCSV(reject),
			})
			response, status, reqErr := client.Do(
				"PUT",
				"/api/v1/external/data-management/quality/rules",
				body,
			)
			return printResponse(cmd, response, status, reqErr)
		},
	}
	cmd.Flags().IntVar(&entityID, "entity-id", 0, "Table ID")
	cmd.Flags().StringVar(&entitySlug, "entity-slug", "", "Table slug (alternative to --entity-id)")
	cmd.Flags().IntVar(&columnID, "column-id", 0, "Column ID")
	cmd.Flags().StringVar(&columnKey, "column-key", "", "Column key (alternative to --column-id)")
	cmd.Flags().StringVar(&accept, "accept", "", "Comma-separated rule ids to accept (e.g. regex:email,unique)")
	cmd.Flags().StringVar(&reject, "reject", "", "Comma-separated rule ids to reject")
	return cmd
}

func newDMQualityCorrectCmd() *cobra.Command {
	var entityID int
	var entitySlug, columnKey string
	var apply bool
	var limit int
	cmd := &cobra.Command{
		Use:   "correct",
		Short: "Preview or apply normalization corrections for a column",
		RunE: func(cmd *cobra.Command, args []string) error {
			if entityID <= 0 && entitySlug == "" {
				return fmt.Errorf("--entity-id or --entity-slug is required")
			}
			if columnKey == "" {
				return fmt.Errorf("--column-key is required")
			}
			client, err := mustClient()
			if err != nil {
				return err
			}
			body, _ := json.Marshal(map[string]any{
				"entity_id":   entityID,
				"entity_slug": entitySlug,
				"column_key":  columnKey,
				"dry_run":     !apply,
				"limit":       limit,
			})
			response, status, reqErr := client.Do(
				"POST",
				"/api/v1/external/data-management/quality/corrections",
				body,
			)
			return printResponse(cmd, response, status, reqErr)
		},
	}
	cmd.Flags().IntVar(&entityID, "entity-id", 0, "Table ID")
	cmd.Flags().StringVar(&entitySlug, "entity-slug", "", "Table slug (alternative to --entity-id)")
	cmd.Flags().StringVar(&columnKey, "column-key", "", "Column key")
	cmd.Flags().BoolVar(&apply, "apply", false, "Apply the corrections (default: preview only)")
	cmd.Flags().IntVar(&limit, "limit", 50, "Max rows to preview/apply")
	return cmd
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
