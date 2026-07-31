package processes

import "github.com/spf13/cobra"

func newSalesInvoicesCmd() *cobra.Command {
	invoiceCmd := &cobra.Command{
		Use:   "sales-invoices",
		Short: "Business verbs for invoice generator workflows",
	}

	invoiceCmd.AddCommand(newInvoiceIssueCmd())
	invoiceCmd.AddCommand(newCollectionListCmd("list", "List invoices", "sales-invoices", "invoices"))
	invoiceCmd.AddCommand(newInvoiceRectifyCmd())
	invoiceCmd.AddCommand(newInvoiceSendCmd())
	invoiceCmd.AddCommand(newInvoiceGetPDFCmd())
	invoiceCmd.AddCommand(newInvoiceLinesAddCmd())
	invoiceCmd.AddCommand(newCollectionListCmd("list-lines", "List invoice lines", "sales-invoices", "invoice_lines"))
	invoiceCmd.AddCommand(newInvoiceLinesUpdateCmd())
	invoiceCmd.AddCommand(newInvoiceLinesDeleteCmd())

	return invoiceCmd
}

func newInvoiceCmd() *cobra.Command {
	cmd := newSalesInvoicesCmd()
	cmd.Use = "invoice"
	cmd.Hidden = true
	cmd.Long = "DEPRECATED: use 'sales-invoices' instead"
	return cmd
}

func newInvoiceGetPDFCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-pdf <invoice_id>",
		Short: "Get the preview URL for an invoice PDF",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "sales-invoices", "get_invoice_pdf", "invoice_id", args[0])
		},
	}
}

func newInvoiceIssueCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "issue <invoice_id>",
		Short: "Issue an invoice",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "sales-invoices", "issue_invoice", "invoice_id", args[0])
		},
	}
}

func newInvoiceRectifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rectify <invoice_id> [reason]",
		Short: "Create a rectificative invoice",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			invoiceID, err := parsePositiveIntArg("invoice_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{"invoice_id": invoiceID}
			if len(args) == 2 {
				payload["reason"] = args[1]
			}
			return runModuleAction(cmd, "sales-invoices", "rectify_invoice", payload)
		},
	}
}

func newInvoiceSendCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "send <invoice_id> [email] [subject] [message]",
		Short: "Send an existing invoice",
		Args:  cobra.RangeArgs(1, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			invoiceID, err := parsePositiveIntArg("invoice_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{"invoice_id": invoiceID}
			if len(args) >= 2 {
				payload["email"] = args[1]
			}
			if len(args) >= 3 {
				payload["subject"] = args[2]
			}
			if len(args) == 4 {
				payload["message"] = args[3]
			}
			return runModuleAction(cmd, "sales-invoices", "send_invoice", payload)
		},
	}
}
