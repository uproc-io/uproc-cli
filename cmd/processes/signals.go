package processes

import "github.com/spf13/cobra"

func newSignalsCmd() *cobra.Command {
	signalsCmd := &cobra.Command{
		Use:   "signals",
		Short: "Business verbs for market signals workflows",
	}

	signalsCmd.AddCommand(newSignalsApproveCmd())
	signalsCmd.AddCommand(newCollectionListCmd("list", "List market signals", "market-signals", "signals"))
	signalsCmd.AddCommand(newCollectionListCmd("list-executions", "List market signal executions", "market-signals", "executions"))
	signalsCmd.AddCommand(newCollectionListCmd("list-activations", "List market signal activations", "market-signals", "activations"))
	signalsCmd.AddCommand(newSignalsDiscardCmd())
	signalsCmd.AddCommand(newSignalsMarkPendingReviewCmd())
	signalsCmd.AddCommand(newSignalsActivateCmd())
	signalsCmd.AddCommand(newSignalsCloseCmd())

	return signalsCmd
}

func newSignalsApproveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "approve <signal_id>",
		Short: "Approve a market signal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "market-signals", "approve_signal", "signal_id", args[0])
		},
	}
}

func newSignalsDiscardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "discard <signal_id>",
		Short: "Discard a market signal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "market-signals", "discard_signal", "signal_id", args[0])
		},
	}
}

func newSignalsMarkPendingReviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-pending-review <signal_id>",
		Short: "Move a market signal back to pending review",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "market-signals", "mark_pending_review", "signal_id", args[0])
		},
	}
}

func newSignalsActivateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "activate <signal_id>",
		Short: "Activate a reviewed market signal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "market-signals", "activate_signal", "signal_id", args[0])
		},
	}
}

func newSignalsCloseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "close <signal_id>",
		Short: "Close a market signal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "market-signals", "close_signal", "signal_id", args[0])
		},
	}
}
