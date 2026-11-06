package processes

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

// newAdminAuditCmd mirrors Admin > System > Audit from the web app.
// Superadmins can query any customer; admins are scoped to their own customer.
func newAdminAuditCmd() *cobra.Command {
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Read the admin audit trail (traces) by module, entity, actor or result",
		Long: `Read the admin audit trail (the same data as Admin > System > Audit).

Admins only see their own customer; superadmins can query any customer or all
of them. Useful to recover traces of data-process tool calls with no results
(filter --module-slug data-process --result failed) and inspect the exact
request/response payload to improve the affected tool.`,
	}

	auditCmd.AddCommand(newAdminAuditListCmd())
	auditCmd.AddCommand(newAdminAuditShowCmd())
	return auditCmd
}

func newAdminAuditListCmd() *cobra.Command {
	var (
		customerID int
		moduleSlug string
		action     string
		actorType  string
		entityType string
		entityID   string
		traceID    string
		result     string
		policy     string
		fromDate   string
		toDate     string
		page       int
		pageSize   int
		sortBy     string
		sortOrder  string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List audit trail entries",
		Example: `  uproc applications admin audit list --module-slug data-process --result failed
  uproc applications admin audit list --module-slug approval-management --action approval_decision
  uproc applications admin audit list --actor-type agent --page-size 50
  uproc applications admin audit list --customer-id 51 --trace-id abc123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			params := url.Values{}
			if customerID > 0 {
				params.Set("customer_id", fmt.Sprintf("%d", customerID))
			}
			for key, value := range map[string]string{
				"module_slug":   moduleSlug,
				"action":        action,
				"actor_type":    actorType,
				"entity_type":   entityType,
				"entity_id":     entityID,
				"trace_id":      traceID,
				"result":        result,
				"policy_result": policy,
				"date_from":     fromDate,
				"date_to":       toDate,
				"sort_by":       sortBy,
				"sort_order":    sortOrder,
			} {
				if strings.TrimSpace(value) != "" {
					params.Set(key, value)
				}
			}
			params.Set("page", fmt.Sprintf("%d", page))
			params.Set("page_size", fmt.Sprintf("%d", pageSize))

			body, status, reqErr := client.Do(
				"GET",
				"/api/v1/external/admin/audit?"+params.Encode(),
				nil,
			)
			return printResponse(cmd, body, status, reqErr)
		},
	}

	cmd.Flags().IntVar(&customerID, "customer-id", 0, "optional customer id (superadmin only)")
	cmd.Flags().StringVar(&moduleSlug, "module-slug", "", "filter by module slug (e.g. data-process)")
	cmd.Flags().StringVar(&action, "action", "", "filter by action (create, update, delete, change_status, tool_call, approval_decision, ...)")
	cmd.Flags().StringVar(&actorType, "actor-type", "", "filter by actor type (user, system, agent)")
	cmd.Flags().StringVar(&entityType, "entity-type", "", "filter by entity type (e.g. tool, approval)")
	cmd.Flags().StringVar(&entityID, "entity-id", "", "filter by entity id")
	cmd.Flags().StringVar(&traceID, "trace-id", "", "filter by trace id")
	cmd.Flags().StringVar(&result, "result", "", "filter by result (success, failed, blocked)")
	cmd.Flags().StringVar(&policy, "policy-result", "", "filter by policy result")
	cmd.Flags().StringVar(&fromDate, "from-date", "", "start date (YYYY-MM-DD or RFC3339)")
	cmd.Flags().StringVar(&toDate, "to-date", "", "end date (YYYY-MM-DD or RFC3339)")
	cmd.Flags().IntVar(&page, "page", 1, "page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 25, "items per page")
	cmd.Flags().StringVar(&sortBy, "sort-field", "", "sort field (e.g. created_at)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort order (asc, desc)")
	return cmd
}

func newAdminAuditShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <record_id>",
		Short: "Show one audit entry with before/after diff and data-process run payload",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/api/v1/external/admin/audit/%s", args[0])
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}
