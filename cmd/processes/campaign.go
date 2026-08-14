package processes

import "github.com/spf13/cobra"

func newCampaignAutomationCmd() *cobra.Command {
	campaignCmd := &cobra.Command{
		Use:   "campaign-automation",
		Short: "Business verbs for campaign-automation",
	}

	campaignCmd.AddCommand(newCampaignPreviewAudienceCmd())
	campaignCmd.AddCommand(newCollectionListCmd("list", "List campaign-automation campaigns", "campaign-automation", "campaigns"))
	campaignCmd.AddCommand(newCollectionListCmd("list-audiences", "List campaign-automation audiences", "campaign-automation", "audiences"))
	campaignCmd.AddCommand(newCampaignAddAudienceCmd())
	campaignCmd.AddCommand(newCampaignPauseCmd())
	campaignCmd.AddCommand(newCampaignActivateCmd())
	campaignCmd.AddCommand(newCampaignSyncAudienceCmd())
	campaignCmd.AddCommand(newCampaignLaunchPlatformsCmd())

	return campaignCmd
}

func newCampaignCmd() *cobra.Command {
	cmd := newCampaignAutomationCmd()
	cmd.Use = "campaign"
	cmd.Hidden = true
	cmd.Long = "DEPRECATED: Use \"campaign-automation\" instead.\n\n" + cmd.Long
	return cmd
}

func newCampaignPreviewAudienceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "preview-audience <campaign_id> [limit]",
		Short: "Preview a campaign automation audience",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			campaignID, err := parsePositiveIntArg("campaign_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{"campaign_id": campaignID}
			if len(args) == 2 {
				limit, err := parsePositiveIntArg("limit", args[1])
				if err != nil {
					return err
				}
				payload["limit"] = limit
			}
			return runModuleAction(cmd, "campaign-automation", "preview_audience", payload)
		},
	}
}

func newCampaignAddAudienceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-audience <campaign_id> [mode]",
		Short: "Add an audience definition to a campaign",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			campaignID, err := parsePositiveIntArg("campaign_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{"campaign_id": campaignID}
			if len(args) == 2 {
				payload["mode"] = args[1]
			}
			return runModuleAction(cmd, "campaign-automation", "add_audience", payload)
		},
	}
}

func newCampaignPauseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pause <campaign_id>",
		Short: "Pause a campaign automation flow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "campaign-automation", "pause", "campaign_id", args[0])
		},
	}
}

func newCampaignActivateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "activate <campaign_id>",
		Short: "Activate a campaign automation flow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "campaign-automation", "activate", "campaign_id", args[0])
		},
	}
}

func newCampaignSyncAudienceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync-audience <campaign_id> <audience_id>",
		Short: "Sync a campaign audience to configured platforms (e.g. Sendy)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			campaignID, err := parsePositiveIntArg("campaign_id", args[0])
			if err != nil {
				return err
			}
			audienceID, err := parsePositiveIntArg("audience_id", args[1])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "campaign-automation", "sync_audience", map[string]any{
				"campaign_id": campaignID,
				"audience_id": audienceID,
			})
		},
	}
}

func newCampaignLaunchPlatformsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "launch-platforms <campaign_id>",
		Short: "Create and launch a campaign on its configured platforms (e.g. Sendy)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			campaignID, err := parsePositiveIntArg("campaign_id", args[0])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "campaign-automation", "launch_platforms", map[string]any{
				"campaign_id": campaignID,
			})
		},
	}
}
