package processes

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewDataColumnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "column",
		Short: "Data-management column commands",
	}

	cmd.AddCommand(newDataColumnUpdateCmd())
	return cmd
}

func newDataColumnUpdateCmd() *cobra.Command {
	var entitySlug, columnKey, colType, relationEntitySlug, label string
	var entityID, columnID int
	var visible, editable string

	c := &cobra.Command{
		Use:   "update",
		Short: "Update a data-management column",
		Long: `Update column properties in a data-management table.
Supports entity slug + column key OR entity ID + column ID.

Examples:
  uproc data column update --entity-slug dm-pacientes --column-key centro_id --type relation --relation-entity-slug dm-centros
  uproc data column update --entity-slug dm-pacientes --column-key nombre_paciente --editable false
  uproc data column update --entity-id 123 --column-id 456 --editable false --visible false`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			item := make(map[string]any)

			if colType != "" {
				item["type"] = colType
			}
			if label != "" {
				item["label"] = label
			}
			if visible != "" {
				switch strings.ToLower(visible) {
				case "true", "1", "yes":
					item["visible"] = true
				case "false", "0", "no":
					item["visible"] = false
				default:
					return fmt.Errorf("invalid visible value: %q (use true/false)", visible)
				}
			}
			if editable != "" {
				switch strings.ToLower(editable) {
				case "true", "1", "yes":
					item["editable"] = true
				case "false", "0", "no":
					item["editable"] = false
				default:
					return fmt.Errorf("invalid editable value: %q (use true/false)", editable)
				}
			}

			arguments := map[string]any{
				"item": item,
			}

			if entityID > 0 {
				arguments["entity_id"] = entityID
			} else if entitySlug != "" {
				arguments["entity_slug"] = entitySlug
			} else {
				return fmt.Errorf("need --entity-id or --entity-slug")
			}

			if columnID > 0 {
				arguments["column_id"] = columnID
			} else if columnKey != "" {
				arguments["column_key"] = columnKey
			} else {
				return fmt.Errorf("need --column-id or --column-key")
			}

			if relationEntitySlug != "" {
				if colType == "relation" {
					arguments["relation_entity_slug"] = relationEntitySlug
				} else {
					cmd.PrintErrln("warning: --relation-entity-slug requires --type relation")
				}
			}

			body, _ := json.Marshal(map[string]any{
				"name":      "data-management.update_column",
				"arguments": arguments,
			})

			respBody, status, reqErr := client.Do("POST", "/api/v1/external/mcp/call", body)
			return printResponse(cmd, respBody, status, reqErr)
		},
	}

	c.Flags().StringVarP(&entitySlug, "entity-slug", "e", "", "Entity slug (e.g. dm-pacientes)")
	c.Flags().StringVarP(&columnKey, "column-key", "c", "", "Column key (e.g. centro_id)")
	c.Flags().IntVar(&entityID, "entity-id", 0, "Entity ID (alternative to --entity-slug)")
	c.Flags().IntVar(&columnID, "column-id", 0, "Column ID (alternative to --column-key)")
	c.Flags().StringVarP(&colType, "type", "t", "", "Column type (text|number|date|boolean|relation)")
	c.Flags().StringVar(&relationEntitySlug, "relation-entity-slug", "", "Target entity slug for relation type")
	c.Flags().StringVarP(&label, "label", "l", "", "Column label")
	c.Flags().StringVar(&editable, "editable", "", "Whether column is editable (true/false)")
	c.Flags().StringVar(&visible, "visible", "", "Whether column is visible (true/false)")

	return c
}
