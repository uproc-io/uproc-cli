package processes

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	processesCmd := &cobra.Command{
		Use:   "processes",
		Short: "Commands for Uproc Processes API",
	}

	// === Generic commands (keep) ===
	processesCmd.AddCommand(newLoginCmd())
	processesCmd.AddCommand(newRequestCmd())
	processesCmd.AddCommand(newModuleCmd())

	// === Real slug commands (primary) ===
	processesCmd.AddCommand(newCampaignAutomationCmd())
	processesCmd.AddCommand(newMarketSignalsCmd())
	processesCmd.AddCommand(newEditorialEngineCmd())
	processesCmd.AddCommand(newLeadManagementCmd())
	processesCmd.AddCommand(newLeadProspectingCmd())
	processesCmd.AddCommand(newFinancialReconciliationCmd())
	processesCmd.AddCommand(newDataSyncCmd())
	processesCmd.AddCommand(newCaseLifecycleCmd())
	processesCmd.AddCommand(newContractLifecycleCmd())
	processesCmd.AddCommand(newInventoryPlanningCmd())
	processesCmd.AddCommand(newOrderIngestCmd())
	processesCmd.AddCommand(newTaxReportingCmd())
	processesCmd.AddCommand(newEmailAssistantCmd())
	processesCmd.AddCommand(newOrderTrackCmd())
	processesCmd.AddCommand(newSalesInvoicesCmd())
	processesCmd.AddCommand(newCustomerCareCmd())
	processesCmd.AddCommand(newApprovalManagementCmd())
	processesCmd.AddCommand(newCandidateEvaluationCmd())
	processesCmd.AddCommand(newDataChatbotCmd())
	processesCmd.AddCommand(newProcessVisibilityCmd())
	processesCmd.AddCommand(newDocumentGeneratorCmd())
	processesCmd.AddCommand(newDocumentSigningCmd())

	// === New module commands (were missing) ===
	processesCmd.AddCommand(newDebtTrackCmd())
	processesCmd.AddCommand(newLeadIntelligenceCmd())
	processesCmd.AddCommand(newDataManagementCmd())
	processesCmd.AddCommand(newFormGeneratorCmd())
	processesCmd.AddCommand(newDocumentIngestCmd())
	processesCmd.AddCommand(newInvoiceIngestCmd())

	// === Admin & system ===
	processesCmd.AddCommand(newAdminCmd())
	processesCmd.AddCommand(newInstallCmd())
	processesCmd.AddCommand(newUpdateCmd())
	processesCmd.AddCommand(newInteractiveCmd())

	// === Hidden aliases (backwards compatibility) ===
	hide := func(cmd *cobra.Command) *cobra.Command { cmd.Hidden = true; return cmd }
	processesCmd.AddCommand(hide(newCampaignCmd()))
	processesCmd.AddCommand(hide(newSignalsCmd()))
	processesCmd.AddCommand(hide(newEditorialCmd()))
	processesCmd.AddCommand(hide(newLeadsCmd()))
	processesCmd.AddCommand(hide(newProspectingCmd()))
	processesCmd.AddCommand(hide(newReconciliationCmd()))
	processesCmd.AddCommand(hide(newSyncCmd()))
	processesCmd.AddCommand(hide(newCasesCmd()))
	processesCmd.AddCommand(hide(newContractCmd()))
	processesCmd.AddCommand(hide(newInventoryCmd()))
	processesCmd.AddCommand(hide(newOrdersIngestCmd()))
	processesCmd.AddCommand(hide(newTaxCmd()))
	processesCmd.AddCommand(hide(newEmailCmd()))
	processesCmd.AddCommand(hide(newOrderCmd()))
	processesCmd.AddCommand(hide(newInvoiceCmd()))
	processesCmd.AddCommand(hide(newSupportCmd()))
	processesCmd.AddCommand(hide(newApprovalCmd()))
	processesCmd.AddCommand(hide(newCandidateCmd()))
	processesCmd.AddCommand(hide(newChatCmd()))
	processesCmd.AddCommand(hide(newProcessCmd()))
	processesCmd.AddCommand(hide(newDocumentsCmd()))
	processesCmd.AddCommand(hide(newSigningCmd()))
	processesCmd.AddCommand(hide(newFormsCmd()))
	processesCmd.AddCommand(hide(newInvoiceLinesCmd()))

	return processesCmd
}
