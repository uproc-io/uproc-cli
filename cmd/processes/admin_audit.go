package processes

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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
	auditCmd.AddCommand(newAdminAuditCaptureCmd())
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

// newAdminAuditCaptureCmd walks the audit trail, detects data-process tool calls
// that returned no result and writes one markdown occurrence per detection under
// <out>/<tool>/occurrences/.
func newAdminAuditCaptureCmd() *cobra.Command {
	var (
		outDir     string
		moduleSlug string
		entityType string
		fromDate   string
		toDate     string
		pageSize   int
		maxPages   int
		dryRun     bool
		includeRaw bool
	)

	cmd := &cobra.Command{
		Use:   "capture",
		Short: "Detect no-result data-process tool traces and write them to docs/plans/issues/data-process",
		Long: `Walk the admin audit trail, detect data-process tool calls with no result
(audit after.matched = false, run free_reason = no_result, or reason = no_result)
and write one markdown occurrence per detection under:

  <out>/<tool>/occurrences/<date>-<trace>.md

Emails and IP addresses are redacted by default. Existing occurrence files are
never overwritten (dedupe by file name / trace). Use --dry-run to preview.`,
		Example: `  uproc applications admin audit capture --out developments/backend/docs/plans/issues/data-process
  uproc applications admin audit capture --out docs/plans/issues/data-process --dry-run
  uproc applications admin audit capture --out docs/plans/issues/data-process --from-date 2026-01-01 --page-size 100`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			detected := 0
			written := 0
			skipped := 0
			page := 1
			for page <= maxPages {
				params := url.Values{}
				params.Set("module_slug", moduleSlug)
				params.Set("entity_type", entityType)
				params.Set("page", fmt.Sprintf("%d", page))
				params.Set("page_size", fmt.Sprintf("%d", pageSize))
				if strings.TrimSpace(fromDate) != "" {
					params.Set("date_from", fromDate)
				}
				if strings.TrimSpace(toDate) != "" {
					params.Set("date_to", toDate)
				}

				body, status, reqErr := client.Do(
					"GET",
					"/api/v1/external/admin/audit?"+params.Encode(),
					nil,
				)
				if reqErr != nil {
					return printResponse(cmd, body, status, reqErr)
				}

				var listResp struct {
					Success bool `json:"success"`
					Data    struct {
						Rows     []map[string]any `json:"rows"`
						MaxPages int              `json:"max_pages"`
					} `json:"data"`
				}
				if err := json.Unmarshal(body, &listResp); err != nil {
					return fmt.Errorf("cannot parse audit list: %w", err)
				}

				for _, row := range listResp.Data.Rows {
					recordID := auditInt(row["id"])
					if recordID <= 0 {
						continue
					}
					detail, detailErr := fetchAuditDetail(client, recordID)
					if detailErr != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "skip record %d: %v\n", recordID, detailErr)
						continue
					}
					if !auditNoResult(detail) {
						continue
					}
					detected++

					path, existed, writeErr := writeAuditOccurrence(outDir, detail, !includeRaw, dryRun)
					if writeErr != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "skip record %d: %v\n", recordID, writeErr)
						continue
					}
					if existed {
						skipped++
						continue
					}
					if dryRun {
						fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] would write %s\n", path)
					} else {
						fmt.Fprintf(cmd.OutOrStdout(), "wrote %s\n", path)
					}
					written++
				}

				if listResp.Data.MaxPages <= page || len(listResp.Data.Rows) == 0 {
					break
				}
				page++
			}

			fmt.Fprintf(
				cmd.OutOrStdout(),
				"capture: detected=%d written=%d skipped=%d (out=%s)\n",
				detected,
				written,
				skipped,
				outDir,
			)
			return nil
		},
	}

	cmd.Flags().StringVar(&outDir, "out", "docs/plans/issues/data-process", "output folder (backend docs/plans/issues/data-process)")
	cmd.Flags().StringVar(&moduleSlug, "module-slug", "data-process", "audit module slug to scan")
	cmd.Flags().StringVar(&entityType, "entity-type", "tool", "audit entity type to scan")
	cmd.Flags().StringVar(&fromDate, "from-date", "", "start date (YYYY-MM-DD or RFC3339)")
	cmd.Flags().StringVar(&toDate, "to-date", "", "end date (YYYY-MM-DD or RFC3339)")
	cmd.Flags().IntVar(&pageSize, "page-size", 100, "items per page")
	cmd.Flags().IntVar(&maxPages, "max-pages", 20, "maximum pages to scan")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview without writing files")
	cmd.Flags().BoolVar(&includeRaw, "include-raw", false, "do NOT redact emails/IPs (local only; never commit)")
	return cmd
}

func fetchAuditDetail(client interface {
	Do(method, path string, body []byte) ([]byte, int, error)
}, recordID int) (map[string]any, error) {
	path := fmt.Sprintf("/api/v1/external/admin/audit/%d", recordID)
	body, _, err := client.Do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var detailResp struct {
		Success bool           `json:"success"`
		Data    map[string]any `json:"data"`
	}
	if err := json.Unmarshal(body, &detailResp); err != nil {
		return nil, fmt.Errorf("cannot parse audit detail: %w", err)
	}
	if detailResp.Data == nil {
		return nil, fmt.Errorf("empty audit detail")
	}
	return detailResp.Data, nil
}

// auditNoResult reports whether a detail row is a no-result tool call.
func auditNoResult(detail map[string]any) bool {
	if after, ok := detail["after"].(map[string]any); ok {
		if matched, ok := after["matched"].(bool); ok && !matched {
			return true
		}
	}
	run, _ := detail["run"].(map[string]any)
	if run == nil {
		return false
	}
	if reason, _ := run["free_reason"].(string); strings.EqualFold(strings.TrimSpace(reason), "no_result") {
		return true
	}
	if raw, ok := run["result"].(string); ok && strings.TrimSpace(raw) != "" {
		var parsed map[string]any
		if json.Unmarshal([]byte(raw), &parsed) == nil {
			if reason, _ := parsed["reason"].(string); strings.EqualFold(strings.TrimSpace(reason), "no_result") {
				return true
			}
		}
	}
	return false
}

// writeAuditOccurrence writes <out>/<tool>/occurrences/<date>-<trace>.md.
// Returns the path, whether it already existed, and any write error.
func writeAuditOccurrence(outDir string, detail map[string]any, redact bool, dryRun bool) (string, bool, error) {
	tool := auditString(detail["entity_id"])
	if tool == "" {
		tool = "unknown-tool"
	}
	dir := filepath.Join(outDir, sanitizePathSegment(tool), "occurrences")
	path := filepath.Join(dir, occurrenceFileName(detail))
	if _, err := os.Stat(path); err == nil {
		return path, true, nil
	}

	content := renderAuditOccurrence(detail)
	if redact {
		content = redactPII(content)
	}
	if dryRun {
		return path, false, nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return path, false, err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return path, false, err
	}
	return path, false, nil
}

func occurrenceFileName(detail map[string]any) string {
	date := auditDate(detail["created_at"])
	trace := sanitizePathSegment(auditString(detail["trace_id"]))
	if len(trace) > 12 {
		trace = trace[:12]
	}
	if trace == "" {
		trace = fmt.Sprintf("record-%d", auditInt(detail["id"]))
	}
	return fmt.Sprintf("%s-%s.md", date, trace)
}

func renderAuditOccurrence(detail map[string]any) string {
	tool := auditString(detail["entity_id"])
	run, _ := detail["run"].(map[string]any)
	after, _ := detail["after"].(map[string]any)

	var b strings.Builder
	fmt.Fprintf(&b, "# Occurrence — %s — %s\n\n", tool, auditDate(detail["created_at"]))
	fmt.Fprintf(&b, "- **Detected:** %s\n", auditDate(detail["created_at"]))
	fmt.Fprintf(&b, "- **Detection source:** `uproc applications admin audit capture`\n")
	fmt.Fprintf(&b, "- **Tool:** `%s`\n", tool)
	fmt.Fprintf(&b, "- **Reason:** `no_result`\n")
	fmt.Fprintf(&b, "- **Audit record id:** %d\n", auditInt(detail["id"]))
	fmt.Fprintf(&b, "- **Trace id:** %s\n", valueOrDash(detail["trace_id"]))
	fmt.Fprintf(&b, "- **Customer id:** %s\n", valueOrDash(detail["customer_id"]))
	fmt.Fprintf(&b, "- **Actor:** %s (%s)\n", valueOrDash(detail["actor_name"]), valueOrDash(detail["actor_id"]))
	fmt.Fprintf(&b, "- **Action:** %s\n", valueOrDash(detail["action"]))
	fmt.Fprintf(&b, "- **Result:** %s\n", valueOrDash(detail["result"]))
	if run != nil {
		fmt.Fprintf(&b, "- **Execution mode:** %s\n", valueOrDash(run["execution_mode"]))
		fmt.Fprintf(&b, "- **Free reason:** %s\n", valueOrDash(run["free_reason"]))
		fmt.Fprintf(&b, "- **Status:** %s\n", valueOrDash(run["status"]))
		fmt.Fprintf(&b, "- **Cost:** %s\n", valueOrDash(run["cost"]))
	}
	fmt.Fprintf(&b, "\n## Params\n\n```json\n%s\n```\n", prettyJSON(run["params"]))
	fmt.Fprintf(&b, "\n## Output (run result)\n\n```json\n%s\n```\n", prettyJSON(run["result"]))
	fmt.Fprintf(&b, "\n## Audit `after`\n\n```json\n%s\n```\n", prettyJSON(after))
	fmt.Fprintf(&b, "\n## Audit record\n\n```json\n%s\n```\n", prettyJSON(detail))
	fmt.Fprintf(&b, "\n## Notes\n\n- Expected vs actual:\n- Provider / fallback suspected:\n- Related issue / plan:\n- Reproduction command:\n")
	return b.String()
}

func prettyJSON(value any) string {
	if value == nil {
		return "{}"
	}
	switch typed := value.(type) {
	case string:
		raw := strings.TrimSpace(typed)
		if raw == "" {
			return "{}"
		}
		var parsed any
		if json.Unmarshal([]byte(raw), &parsed) == nil {
			if encoded, err := json.MarshalIndent(parsed, "", "  "); err == nil {
				return string(encoded)
			}
		}
		return raw
	default:
		if encoded, err := json.MarshalIndent(value, "", "  "); err == nil {
			return string(encoded)
		}
		return fmt.Sprintf("%v", value)
	}
}

var (
	auditEmailRe = regexp.MustCompile(`[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}`)
	auditIPv4Re  = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
)

func redactPII(text string) string {
	text = auditEmailRe.ReplaceAllString(text, "***@***")
	text = auditIPv4Re.ReplaceAllString(text, "x.x.x.x")
	return text
}

func auditDate(value any) string {
	raw := auditString(value)
	if raw == "" {
		return time.Now().UTC().Format("2006-01-02")
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05.999999", "2006-01-02 15:04:05"} {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return parsed.UTC().Format("2006-01-02")
		}
	}
	if len(raw) >= 10 {
		return raw[:10]
	}
	return raw
}

func auditString(value any) string {
	if value == nil {
		return ""
	}
	if text, ok := value.(string); ok {
		return strings.TrimSpace(text)
	}
	return strings.TrimSpace(fmt.Sprintf("%v", value))
}

func auditInt(value any) int {
	switch typed := value.(type) {
	case float64:
		return int(typed)
	case int:
		return typed
	case json.Number:
		parsed, _ := typed.Int64()
		return int(parsed)
	default:
		var parsed int
		_, _ = fmt.Sscanf(fmt.Sprintf("%v", value), "%d", &parsed)
		return parsed
	}
}

func valueOrDash(value any) string {
	text := auditString(value)
	if text == "" {
		return "-"
	}
	return text
}

func sanitizePathSegment(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	replacer := regexp.MustCompile(`[^A-Za-z0-9._-]+`)
	return strings.Trim(replacer.ReplaceAllString(value, "-"), "-")
}
