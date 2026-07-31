package processes

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newPurchaseInvoicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purchase-invoices",
		Short: "Business verbs for purchase invoices workflows",
	}

	cmd.AddCommand(newCollectionListCmd("list", "List purchase invoices", "purchase-invoices", "purchase_invoices"))
	cmd.AddCommand(newCollectionListCmd("list-lines", "List purchase invoice lines", "purchase-invoices", "purchase_invoice_lines"))
	cmd.AddCommand(newCollectionListCmd("list-payments", "List purchase invoice payments", "purchase-invoices", "purchase_invoice_payments"))
	cmd.AddCommand(newPurchaseInvoicesValidateCmd())
	cmd.AddCommand(newPurchaseInvoicesPayCmd())
	cmd.AddCommand(newPurchaseInvoicesAssignFromIngestCmd())

	return cmd
}

func newPurchaseInvoicesValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <id> [reason]",
		Short: "Mark a purchase invoice as validated",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID, err := parsePositiveIntArg("id", args[0])
			if err != nil {
				return err
			}
			reason := ""
			if len(args) == 2 {
				reason = args[1]
			}
			return runModuleStatusChange(cmd, "purchase-invoices", "purchase_invoices", itemID, "validated", reason)
		},
	}
}

func newPurchaseInvoicesPayCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pay <id> [reason]",
		Short: "Mark a purchase invoice as paid",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID, err := parsePositiveIntArg("id", args[0])
			if err != nil {
				return err
			}
			reason := ""
			if len(args) == 2 {
				reason = args[1]
			}
			return runModuleStatusChange(cmd, "purchase-invoices", "purchase_invoices", itemID, "paid", reason)
		},
	}
}

func newPurchaseInvoicesAssignFromIngestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "assign-from-ingest <invoice_id> [payload_json]",
		Short: "Assign a parsed purchase invoice from invoice-ingest",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			parsedInvoiceID, err := parsePositiveIntArg("invoice_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{"invoice_id": parsedInvoiceID}
			if len(args) == 2 && args[1] != "" {
				var extra map[string]any
				if err := json.Unmarshal([]byte(args[1]), &extra); err != nil {
					return fmt.Errorf("invalid payload_json: %w", err)
				}
				payload["payload"] = extra
			}
			return runModuleAction(cmd, "purchase-invoices", "assign_from_ingest", payload)
		},
	}
}

func runModuleStatusChange(cmd *cobra.Command, moduleSlug, collectionName string, itemID int, status, reason string) error {
	client, err := mustClient()
	if err != nil {
		return err
	}

	payload := map[string]any{"status": status}
	if reason != "" {
		payload["reason"] = reason
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cannot encode %s payload: %w", status, err)
	}

	path := fmt.Sprintf(
		"/api/v1/external/modules/%s/data/%s/%d/status",
		moduleSlug,
		collectionName,
		itemID,
	)
	respBody, statusCode, reqErr := client.Do("PUT", path, body)
	return printResponse(cmd, respBody, statusCode, reqErr)
}
