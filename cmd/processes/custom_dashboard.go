package processes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// newCustomDashboardCmd groups the portable custom-dashboard setup commands.
// The setup bundle contains layout_json (pages/blocks/summary.queries),
// dashboard_sections and entity_mappings; it never contains users or secrets.
func newCustomDashboardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom-dashboard",
		Short: "Export/import a customer's custom dashboard setup",
	}
	cmd.AddCommand(newCDExportSetupCmd(), newCDImportSetupCmd())
	return cmd
}

func newCDExportSetupCmd() *cobra.Command {
	var customerID int
	var output string

	cmd := &cobra.Command{
		Use:   "export-setup",
		Short: "Export the custom dashboard setup (layout, sections, mappings) to JSON",
		Long: `Export the custom dashboard setup for a customer as portable JSON.

Examples:
  uproc applications custom-dashboard export-setup --customer-id 51 -o hospicare_dashboard.json
  uproc applications custom-dashboard export-setup --customer-id 51`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID <= 0 {
				return fmt.Errorf("--customer-id is required")
			}
			client, err := mustClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf(
				"/api/v1/external/admin/customers/%d/custom-dashboard/export-setup",
				customerID,
			)
			response, status, reqErr := client.Do("GET", path, nil)
			if reqErr != nil || status != 200 {
				return printResponse(cmd, response, status, reqErr)
			}
			var result struct {
				Success bool            `json:"success"`
				Data    json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(response, &result); err != nil || !result.Success || len(result.Data) == 0 {
				return printResponse(cmd, response, status, nil)
			}
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, result.Data, "", "  "); err != nil {
				pretty.Write(result.Data)
			}
			if output == "" {
				_, err := cmd.OutOrStdout().Write(append(pretty.Bytes(), '\n'))
				return err
			}
			if err := os.WriteFile(output, append(pretty.Bytes(), '\n'), 0600); err != nil {
				return fmt.Errorf("cannot write %s: %w", output, err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Custom dashboard setup saved to %s\n", output)
			return nil
		},
	}
	cmd.Flags().IntVar(&customerID, "customer-id", 0, "Customer ID")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output JSON path (default stdout)")
	return cmd
}

func newCDImportSetupCmd() *cobra.Command {
	var customerID int
	var file string
	var confirm bool

	cmd := &cobra.Command{
		Use:   "import-setup",
		Short: "Import a custom dashboard setup JSON (preview by default; --confirm applies)",
		Long: `Import a custom dashboard setup previously exported with export-setup.

Without --confirm the API returns a preview (dry_run). With --confirm it
updates the active layout, replaces the registered sections and upserts
entity mappings. Users, passwords and DM rows are never touched.

Examples:
  uproc applications custom-dashboard import-setup --customer-id 51 --file hospicare_dashboard.json
  uproc applications custom-dashboard import-setup --customer-id 51 --file hospicare_dashboard.json --confirm`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID <= 0 {
				return fmt.Errorf("--customer-id is required")
			}
			if file == "" {
				return fmt.Errorf("--file is required")
			}
			raw, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("cannot read %s: %w", file, err)
			}
			var payload map[string]any
			if err := json.Unmarshal(raw, &payload); err != nil {
				return fmt.Errorf("invalid setup JSON: %w", err)
			}
			payload["dry_run"] = !confirm
			body, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			client, err := mustClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf(
				"/api/v1/external/admin/customers/%d/custom-dashboard/import-setup",
				customerID,
			)
			response, status, reqErr := client.Do("POST", path, body)
			return printResponse(cmd, response, status, reqErr)
		},
	}
	cmd.Flags().IntVar(&customerID, "customer-id", 0, "Customer ID")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Setup JSON file")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "Apply the setup (without it, preview only)")
	return cmd
}
