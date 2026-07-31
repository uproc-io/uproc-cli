package processes

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newErpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "erp",
		Short: "Write records to ERP providers (Holded, Factura Directa, Invoice Ninja)",
	}
	cmd.AddCommand(newErpCreateCmd())
	cmd.AddCommand(newErpUpdateCmd())
	cmd.AddCommand(newErpDeleteCmd())
	return cmd
}

func newErpCreateCmd() *cobra.Command {
	var credentialID int
	var resource string
	var payloadJSON string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a record on an ERP provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runErpWrite(cmd, "create", credentialID, resource, "", payloadJSON)
		},
	}
	cmd.Flags().IntVar(&credentialID, "credential", 0, "ERP credential ID")
	cmd.Flags().StringVar(&resource, "resource", "", "resource (invoices, contacts, clients, payments...)")
	cmd.Flags().StringVar(&payloadJSON, "json", "{}", "provider payload as JSON")
	_ = cmd.MarkFlagRequired("credential")
	_ = cmd.MarkFlagRequired("resource")
	return cmd
}

func newErpUpdateCmd() *cobra.Command {
	var credentialID int
	var resource string
	var externalID string
	var payloadJSON string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a record on an ERP provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runErpWrite(cmd, "update", credentialID, resource, externalID, payloadJSON)
		},
	}
	cmd.Flags().IntVar(&credentialID, "credential", 0, "ERP credential ID")
	cmd.Flags().StringVar(&resource, "resource", "", "resource")
	cmd.Flags().StringVar(&externalID, "external-id", "", "provider record ID")
	cmd.Flags().StringVar(&payloadJSON, "json", "{}", "provider payload as JSON")
	_ = cmd.MarkFlagRequired("credential")
	_ = cmd.MarkFlagRequired("resource")
	_ = cmd.MarkFlagRequired("external-id")
	return cmd
}

func newErpDeleteCmd() *cobra.Command {
	var credentialID int
	var resource string
	var externalID string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a record on an ERP provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runErpWrite(cmd, "delete", credentialID, resource, externalID, "{}")
		},
	}
	cmd.Flags().IntVar(&credentialID, "credential", 0, "ERP credential ID")
	cmd.Flags().StringVar(&resource, "resource", "", "resource")
	cmd.Flags().StringVar(&externalID, "external-id", "", "provider record ID")
	_ = cmd.MarkFlagRequired("credential")
	_ = cmd.MarkFlagRequired("resource")
	_ = cmd.MarkFlagRequired("external-id")
	return cmd
}

func runErpWrite(cmd *cobra.Command, operation string, credentialID int, resource, externalID, payloadJSON string) error {
	if credentialID <= 0 {
		return fmt.Errorf("--credential is required")
	}
	if resource == "" {
		return fmt.Errorf("--resource is required")
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		return fmt.Errorf("invalid --json payload: %w", err)
	}
	body := map[string]any{
		"resource":    resource,
		"operation":   operation,
		"external_id": externalID,
		"item":        payload,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("cannot encode erp payload: %w", err)
	}
	client, err := mustClient()
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/api/v1/external/erp/credentials/%d/write", credentialID)
	respBody, status, reqErr := client.Do("POST", path, bodyBytes)
	return printResponse(cmd, respBody, status, reqErr)
}
