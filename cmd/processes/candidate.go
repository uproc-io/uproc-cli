package processes

import "github.com/spf13/cobra"

func newCandidateCmd() *cobra.Command {
	candidateCmd := &cobra.Command{
		Use:   "candidate",
		Short: "Business verbs for candidate evaluation workflows",
	}

	candidateCmd.AddCommand(newCandidateCreateProfileCmd())
	candidateCmd.AddCommand(newCandidateCreateJobOpeningCmd())
	candidateCmd.AddCommand(newCandidateCreateApplicationCmd())
	candidateCmd.AddCommand(newCandidateMoveStageCmd())
	candidateCmd.AddCommand(newCandidateUpdateStatusCmd())
	candidateCmd.AddCommand(newCandidateCreateEvaluationCmd())

	return candidateCmd
}

func newCandidateCreateProfileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-profile <item_json>",
		Short: "Create a candidate-evaluation profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			item, err := parseJSONObjectArg("item_json", args[0])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "candidate-evaluation", "create_profile", map[string]any{"item": item})
		},
	}
}

func newCandidateCreateJobOpeningCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-job-opening <item_json>",
		Short: "Create a candidate-evaluation job opening",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			item, err := parseJSONObjectArg("item_json", args[0])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "candidate-evaluation", "create_job_opening", map[string]any{"item": item})
		},
	}
}

func newCandidateCreateApplicationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-application <item_json>",
		Short: "Create a candidate-evaluation application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			item, err := parseJSONObjectArg("item_json", args[0])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "candidate-evaluation", "create_application", map[string]any{"item": item})
		},
	}
}

func newCandidateMoveStageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "move-stage <application_id> <stage>",
		Short: "Move a candidate application to a new stage",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			applicationID, err := parsePositiveIntArg("application_id", args[0])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "candidate-evaluation", "move_application_stage", map[string]any{
				"application_id": applicationID,
				"stage":          args[1],
			})
		},
	}
}

func newCandidateUpdateStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-status <application_id> <status>",
		Short: "Update a candidate application status",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			applicationID, err := parsePositiveIntArg("application_id", args[0])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "candidate-evaluation", "update_application_status", map[string]any{
				"application_id": applicationID,
				"status":         args[1],
			})
		},
	}
}

func newCandidateCreateEvaluationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-evaluation <item_json>",
		Short: "Create a candidate-evaluation evaluation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			item, err := parseJSONObjectArg("item_json", args[0])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "candidate-evaluation", "create_evaluation", map[string]any{"item": item})
		},
	}
}
