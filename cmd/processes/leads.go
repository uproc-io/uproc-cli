package processes

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

func newLeadManagementCmd() *cobra.Command {
	leadsCmd := &cobra.Command{
		Use:   "lead-management",
		Short: "Business verbs for lead-management",
	}

	leadsCmd.AddCommand(newLeadsGenerateProposalCmd())
	leadsCmd.AddCommand(newLeadsListCmd())
	leadsCmd.AddCommand(newLeadsSendProposalCmd())
	leadsCmd.AddCommand(newLeadsRerunIntelligenceCmd())
	leadsCmd.AddCommand(newLeadsUploadProposalCmd())
	leadsCmd.AddCommand(newLeadsSendProposalVersionCmd())

	return leadsCmd
}

func newLeadsCmd() *cobra.Command {
	cmd := newLeadManagementCmd()
	cmd.Use = "leads"
	cmd.Hidden = true
	cmd.Long = "DEPRECATED: Use \"lead-management\" instead.\n\n" + cmd.Long
	return cmd
}

func newLeadsListCmd() *cobra.Command {
	var page int
	var sortField string
	var sortOrder string
	var filterField string
	var filterValue string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List lead management leads",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			query := url.Values{}
			query.Set("page", fmt.Sprintf("%d", page))
			if sortField != "" {
				query.Set("sort_field", sortField)
			}
			if sortOrder != "" {
				query.Set("sort_order", sortOrder)
			}
			if filterField != "" {
				query.Set("filter_field", filterField)
			}
			if filterValue != "" {
				query.Set("filter_value", filterValue)
			}

			path := fmt.Sprintf(
				"/api/v1/external/modules/%s/collections/%s?%s",
				"lead-management",
				"leads",
				query.Encode(),
			)

			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}

	cmd.Flags().IntVar(&page, "page", 1, "page number")
	cmd.Flags().StringVar(&sortField, "sort-field", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort order: asc|desc")
	cmd.Flags().StringVar(&filterField, "filter-field", "", "filter field")
	cmd.Flags().StringVar(&filterValue, "filter-value", "", "filter value")

	return cmd
}

func newLeadsSendProposalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "send-proposal <lead_id> <mailbox_id> <to_email> <subject> <body> [proposal_url]",
		Short: "Send a generated proposal to a lead by email",
		Args:  cobra.RangeArgs(5, 6),
		RunE: func(cmd *cobra.Command, args []string) error {
			leadID, err := parsePositiveIntArg("lead_id", args[0])
			if err != nil {
				return err
			}
			mailboxID, err := parsePositiveIntArg("mailbox_id", args[1])
			if err != nil {
				return err
			}
			payload := map[string]any{
				"lead_id":    leadID,
				"mailbox_id": mailboxID,
				"to_email":   args[2],
				"subject":    args[3],
				"body":       args[4],
			}
			if len(args) == 6 {
				payload["proposal_url"] = args[5]
			}
			return runModuleAction(cmd, "lead-management", "send_lead_proposal", payload)
		},
	}
}

func newLeadsGenerateProposalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate-proposal <lead_id> [template_id] [title] [description] [output_format]",
		Short: "Generate a commercial proposal for a lead",
		Args:  cobra.RangeArgs(1, 5),
		RunE: func(cmd *cobra.Command, args []string) error {
			leadID, err := parsePositiveIntArg("lead_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{"lead_id": leadID}
			if len(args) >= 2 {
				templateID, err := parsePositiveIntArg("template_id", args[1])
				if err != nil {
					return err
				}
				payload["template_id"] = templateID
			}
			if len(args) >= 3 {
				payload["title"] = args[2]
			}
			if len(args) >= 4 {
				payload["description"] = args[3]
			}
			if len(args) == 5 {
				payload["output_format"] = args[4]
			}
			return runModuleAction(cmd, "lead-management", "generate_lead_proposal", payload)
		},
	}
}

func newLeadsRerunIntelligenceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rerun-intelligence <lead_id>",
		Short: "Rerun lead intelligence enrichment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "lead-management", "rerun_lead_intelligence", "lead_id", args[0])
		},
	}
}

func newLeadsUploadProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-proposal <lead_id> <pdf_file>",
		Short: "Upload a proposal PDF to a lead and create a version (does NOT send)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			leadID, err := parsePositiveIntArg("lead_id", args[0])
			if err != nil {
				return err
			}
			pdfPath := args[1]

			pdfBytes, err := os.ReadFile(pdfPath)
			if err != nil {
				return fmt.Errorf("reading PDF file: %w", err)
			}
			pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)

			payload := map[string]any{
				"lead_id":     leadID,
				"pdf_content": pdfBase64,
				"file_name":   pdfPath,
			}
			return runModuleAction(cmd, "lead-management", "upload_proposal", payload)
		},
	}
	return cmd
}

func newLeadsSendProposalVersionCmd() *cobra.Command {
	var toEmail string
	var subject string
	var body string

	cmd := &cobra.Command{
		Use:   "send-proposal-version <version_id>",
		Short: "Send an existing proposal version to the client by email",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			versionID, err := parsePositiveIntArg("version_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{
				"version_id": versionID,
				"to_email":   toEmail,
				"subject":    subject,
				"body":       body,
			}
			return runModuleAction(cmd, "lead-management", "send_proposal_version", payload)
		},
	}

	cmd.Flags().StringVar(&toEmail, "to-email", "", "Recipient email (default: lead's email)")
	cmd.Flags().StringVar(&subject, "subject", "", "Email subject")
	cmd.Flags().StringVar(&body, "body", "", "Email body")

	return cmd
}
