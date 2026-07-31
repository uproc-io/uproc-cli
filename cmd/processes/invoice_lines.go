package processes

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newInvoiceLinesCmd() *cobra.Command {
	invoiceLinesCmd := &cobra.Command{
		Use:   "invoice-lines",
		Short: "Business verbs for invoice line workflows",
	}

	invoiceLinesCmd.AddCommand(newInvoiceLinesAddCmd())
	invoiceLinesCmd.AddCommand(newCollectionListCmd("list", "List invoice lines", "sales-invoices", "invoice_lines"))
	invoiceLinesCmd.AddCommand(newInvoiceLinesUpdateCmd())
	invoiceLinesCmd.AddCommand(newInvoiceLinesDeleteCmd())

	return invoiceLinesCmd
}

func newInvoiceLinesAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <invoice_id> <concept> [quantity] [unit_price] [tax_rate] [sort_order]",
		Short: "Add a line to an existing invoice",
		Args:  cobra.RangeArgs(2, 6),
		RunE: func(cmd *cobra.Command, args []string) error {
			invoiceID, err := parsePositiveIntArg("invoice_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{
				"invoice_id": invoiceID,
				"concept":    args[1],
			}
			if len(args) >= 3 {
				quantity, err := strconv.ParseFloat(args[2], 64)
				if err != nil {
					return err
				}
				payload["quantity"] = quantity
			}
			if len(args) >= 4 {
				unitPrice, err := strconv.ParseFloat(args[3], 64)
				if err != nil {
					return err
				}
				payload["unit_price"] = unitPrice
			}
			if len(args) >= 5 {
				taxRate, err := strconv.ParseFloat(args[4], 64)
				if err != nil {
					return err
				}
				payload["tax_rate"] = taxRate
			}
			if len(args) == 6 {
				sortOrder, err := parseNonNegativeIntArg("sort_order", args[5])
				if err != nil {
					return err
				}
				payload["sort_order"] = sortOrder
			}
			return runModuleAction(cmd, "sales-invoices", "add_invoice_line", payload)
		},
	}
}

func newInvoiceLinesUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update <invoice_id> <line_id> [concept] [quantity] [unit_price] [tax_rate] [sort_order]",
		Short: "Update a line in an existing invoice",
		Args:  cobra.RangeArgs(2, 7),
		RunE: func(cmd *cobra.Command, args []string) error {
			invoiceID, err := parsePositiveIntArg("invoice_id", args[0])
			if err != nil {
				return err
			}
			lineID, err := parsePositiveIntArg("line_id", args[1])
			if err != nil {
				return err
			}
			payload := map[string]any{
				"invoice_id": invoiceID,
				"line_id":    lineID,
			}
			if len(args) >= 3 {
				payload["concept"] = args[2]
			}
			if len(args) >= 4 {
				quantity, err := strconv.ParseFloat(args[3], 64)
				if err != nil {
					return err
				}
				payload["quantity"] = quantity
			}
			if len(args) >= 5 {
				unitPrice, err := strconv.ParseFloat(args[4], 64)
				if err != nil {
					return err
				}
				payload["unit_price"] = unitPrice
			}
			if len(args) >= 6 {
				taxRate, err := strconv.ParseFloat(args[5], 64)
				if err != nil {
					return err
				}
				payload["tax_rate"] = taxRate
			}
			if len(args) == 7 {
				sortOrder, err := parseNonNegativeIntArg("sort_order", args[6])
				if err != nil {
					return err
				}
				payload["sort_order"] = sortOrder
			}
			return runModuleAction(cmd, "sales-invoices", "update_invoice_line", payload)
		},
	}
}

func parseNonNegativeIntArg(name, raw string) (int, error) {
	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer", name)
	}
	return value, nil
}

func newInvoiceLinesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <invoice_id> <line_id>",
		Short: "Delete a line from an existing invoice",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			invoiceID, err := parsePositiveIntArg("invoice_id", args[0])
			if err != nil {
				return err
			}
			lineID, err := parsePositiveIntArg("line_id", args[1])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "sales-invoices", "delete_invoice_line", map[string]any{
				"invoice_id": invoiceID,
				"line_id":    lineID,
			})
		},
	}
}
