package processes

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	applicationsCmd := &cobra.Command{
		Use:   "applications",
		Short: "Commands for Uproc Applications API",
	}
	applicationsCmd.AddCommand(buildProcessChildren()...)
	return applicationsCmd
}

// NewProcessesAliasCmd returns the deprecated hidden `processes` alias.
func NewProcessesAliasCmd() *cobra.Command {
	processesAlias := &cobra.Command{
		Use:    "processes",
		Short:  "DEPRECATED: use 'applications' instead",
		Hidden: true,
	}
	processesAlias.AddCommand(buildProcessChildren()...)
	return processesAlias
}

func buildProcessChildren() []*cobra.Command {
	children := []*cobra.Command{
		// === Generic commands (keep) ===
		newLoginCmd(),
		newRequestCmd(),
		newModuleCmd(),

		// === Real slug commands (primary) ===
		newCampaignAutomationCmd(),
		newMarketSignalsCmd(),
		newEditorialEngineCmd(),
		newLeadManagementCmd(),
		newLeadProspectingCmd(),
		newFinancialReconciliationCmd(),
		newDataSyncCmd(),
		newCaseLifecycleCmd(),
		newContractLifecycleCmd(),
		newInventoryPlanningCmd(),
		newOrderIngestCmd(),
		newTaxReportingCmd(),
		newEmailAssistantCmd(),
		newOrderTrackCmd(),
		newFileSyncCmd(),
		newSalesInvoicesCmd(),
		newPurchaseInvoicesCmd(),
		newDataProcessCmd(),
		newCustomerCareCmd(),
		newApprovalManagementCmd(),
		newCandidateEvaluationCmd(),
		newDataChatbotCmd(),
		newProcessVisibilityCmd(),
		newDocumentGeneratorCmd(),
		newDocumentSigningCmd(),

		// === New module commands (were missing) ===
		newDebtTrackCmd(),
		newLeadIntelligenceCmd(),
		newDataManagementCmd(),
		newFormGeneratorCmd(),
		newDocumentIngestCmd(),
		newInvoiceIngestCmd(),

		// === Admin & system ===
		newAdminCmd(),
		newInstallCmd(),
		newUpdateCmd(),
		newInteractiveCmd(),

		// === Hidden aliases (backwards compatibility) ===
		hide(newCampaignCmd()),
		hide(newSignalsCmd()),
		hide(newEditorialCmd()),
		hide(newLeadsCmd()),
		hide(newProspectingCmd()),
		hide(newReconciliationCmd()),
		hide(newSyncCmd()),
		hide(newCasesCmd()),
		hide(newContractCmd()),
		hide(newInventoryCmd()),
		hide(newOrdersIngestCmd()),
		hide(newTaxCmd()),
		hide(newEmailCmd()),
		hide(newOrderCmd()),
		hide(newInvoiceCmd()),
		hide(newSupportCmd()),
		hide(newApprovalCmd()),
		hide(newCandidateCmd()),
		hide(newChatCmd()),
		hide(newProcessCmd()),
		hide(newDocumentsCmd()),
		hide(newSigningCmd()),
		hide(newFormsCmd()),
		hide(newInvoiceLinesCmd()),
	}

	return children
}

func hide(cmd *cobra.Command) *cobra.Command {
	cmd.Hidden = true
	return cmd
}
