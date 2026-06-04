package processes

import "github.com/spf13/cobra"

func newFormsCmd() *cobra.Command {
	formsCmd := &cobra.Command{
		Use:   "forms",
		Short: "Business verbs for forms workflows",
	}

	formsCmd.AddCommand(newFormsSubmitPublicCmd())
	formsCmd.AddCommand(newFormsPublishCmd())
	formsCmd.AddCommand(newFormsArchiveCmd())
	formsCmd.AddCommand(newFormsRestoreCmd())
	formsCmd.AddCommand(newFormsMarkSubmissionProcessedCmd())
	formsCmd.AddCommand(newFormsArchiveSubmissionCmd())

	return formsCmd
}

func newFormsSubmitPublicCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "submit-public <customer_domain> <form_slug> <payload_json>",
		Short: "Submit a public form payload for form-generator",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPublicFormSubmission(cmd, args[0], args[1], args[2])
		},
	}
}

func newFormsPublishCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "publish <form_id>",
		Short: "Publish a form-generator form",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFormsAction(cmd, "publish", "form_id", args[0])
		},
	}
}

func newFormsArchiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "archive <form_id>",
		Short: "Archive a form-generator form",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFormsAction(cmd, "archive", "form_id", args[0])
		},
	}
}

func newFormsRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore <form_id>",
		Short: "Restore an archived form-generator form",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFormsAction(cmd, "restore", "form_id", args[0])
		},
	}
}

func newFormsMarkSubmissionProcessedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-submission-processed <submission_id>",
		Short: "Mark a form-generator submission as processed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFormsAction(cmd, "mark_submission_processed", "submission_id", args[0])
		},
	}
}

func newFormsArchiveSubmissionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "archive-submission <submission_id>",
		Short: "Archive a form-generator submission",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFormsAction(cmd, "archive_submission", "submission_id", args[0])
		},
	}
}

func runFormsAction(cmd *cobra.Command, action, idField, rawID string) error {
	id, err := parsePositiveIntArg(idField, rawID)
	if err != nil {
		return err
	}
	return runModuleAction(cmd, "form-generator", action, map[string]any{idField: id})
}
