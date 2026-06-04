package processes

import "github.com/spf13/cobra"

func newApprovalCmd() *cobra.Command {
	approvalCmd := &cobra.Command{
		Use:   "approval",
		Short: "Business verbs for approval workflows",
	}

	approvalCmd.AddCommand(newApprovalApproveCmd())
	approvalCmd.AddCommand(newApprovalRejectCmd())
	approvalCmd.AddCommand(newApprovalReassignCmd())
	approvalCmd.AddCommand(newApprovalCancelCmd())

	return approvalCmd
}

func newApprovalApproveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "approve <request_id>",
		Short: "Approve an approval-management request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "approval-management", "approve_request", "request_id", args[0])
		},
	}
}

func newApprovalRejectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reject <request_id>",
		Short: "Reject an approval-management request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "approval-management", "reject_request", "request_id", args[0])
		},
	}
}

func newApprovalReassignCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reassign <request_id> <approver> [note]",
		Short: "Reassign an approval-management request",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			requestID, err := parsePositiveIntArg("request_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{
				"request_id": requestID,
				"approver":   args[1],
			}
			if len(args) == 3 {
				payload["note"] = args[2]
			}
			return runModuleAction(cmd, "approval-management", "reassign_request", payload)
		},
	}
}

func newApprovalCancelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel <request_id>",
		Short: "Cancel an approval-management request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "approval-management", "cancel_request", "request_id", args[0])
		},
	}
}
