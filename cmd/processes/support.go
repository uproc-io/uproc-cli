package processes

import "github.com/spf13/cobra"

func newSupportCmd() *cobra.Command {
	supportCmd := &cobra.Command{
		Use:   "support",
		Short: "Business verbs for customer care workflows",
	}

	supportCmd.AddCommand(newSupportCreateTicketCmd())
	supportCmd.AddCommand(newSupportAssignTicketCmd())
	supportCmd.AddCommand(newSupportReplyTicketCmd())
	supportCmd.AddCommand(newSupportMarkResolvedCmd())
	supportCmd.AddCommand(newSupportCloseTicketCmd())
	supportCmd.AddCommand(newSupportReopenTicketCmd())

	return supportCmd
}

func newSupportCreateTicketCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-ticket <item_json>",
		Short: "Create a customer-care ticket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			item, err := parseJSONObjectArg("item_json", args[0])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "customer-care", "create_ticket", map[string]any{"item": item})
		},
	}
}

func newSupportAssignTicketCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "assign-ticket <ticket_id> <assignee>",
		Short: "Assign a customer-care ticket",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ticketID, err := parsePositiveIntArg("ticket_id", args[0])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "customer-care", "assign_ticket", map[string]any{
				"ticket_id": ticketID,
				"assignee":  args[1],
			})
		},
	}
}

func newSupportReplyTicketCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reply-ticket <ticket_id> <message>",
		Short: "Reply to a customer-care ticket",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ticketID, err := parsePositiveIntArg("ticket_id", args[0])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "customer-care", "reply_to_ticket", map[string]any{
				"ticket_id": ticketID,
				"message":   args[1],
			})
		},
	}
}

func newSupportMarkResolvedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-resolved <ticket_id>",
		Short: "Mark a customer-care ticket as resolved",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "customer-care", "mark_resolved", "ticket_id", args[0])
		},
	}
}

func newSupportCloseTicketCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "close-ticket <ticket_id>",
		Short: "Close a customer-care ticket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "customer-care", "close_ticket", "ticket_id", args[0])
		},
	}
}

func newSupportReopenTicketCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reopen-ticket <ticket_id>",
		Short: "Reopen a customer-care ticket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "customer-care", "reopen_ticket", "ticket_id", args[0])
		},
	}
}
