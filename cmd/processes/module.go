package processes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func newModuleCmd() *cobra.Command {
	moduleCmd := &cobra.Command{
		Use:   "module",
		Short: "Module commands for external API",
	}

	moduleCmd.AddCommand(newModuleListCmd())
	moduleCmd.AddCommand(newModuleGetCmd())
	moduleCmd.AddCommand(newModuleOverviewCmd())
	moduleCmd.AddCommand(newModuleCollectionsCmd())
	moduleCmd.AddCommand(newModuleCollectionCmd())
	moduleCmd.AddCommand(newModuleDataCmd())
	moduleCmd.AddCommand(newModuleSettingsTabsCmd())
	moduleCmd.AddCommand(newModuleSettingsTabCmd())
	moduleCmd.AddCommand(newModuleUploadCmd())
	moduleCmd.AddCommand(newModuleWebhookCmd())
	moduleCmd.AddCommand(newModuleSubmitPublicFormCmd())

	return moduleCmd
}

func newModuleSubmitPublicFormCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-public-form <customer_domain> <form_slug> <payload_json>",
		Short: "Submit a public form payload for form-generator",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPublicFormSubmission(cmd, args[0], args[1], args[2])
		},
	}
	cmd.Deprecated = "use 'uproc processes forms submit-public' instead"
	return cmd
}

func runPublicFormSubmission(cmd *cobra.Command, customerDomain, formSlug, payloadJSON string) error {
	client, err := mustClient()
	if err != nil {
		return err
	}

	payloadBody, err := parsePublicFormSubmissionPayload(payloadJSON)
	if err != nil {
		return err
	}

	path := buildPublicFormSubmissionPath("form-generator", customerDomain, formSlug)
	body, status, reqErr := client.Do("POST", path, payloadBody)
	return printResponse(cmd, body, status, reqErr)
}

func parsePublicFormSubmissionPayload(input string) ([]byte, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil, fmt.Errorf("payload_json is required")
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(trimmed), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload_json: %w", err)
	}

	requestBody := map[string]any{"payload": payload}
	encoded, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("cannot encode payload_json: %w", err)
	}

	return encoded, nil
}

func buildPublicFormSubmissionPath(moduleSlug, customerDomain, formSlug string) string {
	return fmt.Sprintf(
		"/api/v1/external/public/modules/%s/forms/%s/%s/submit",
		url.PathEscape(strings.TrimSpace(moduleSlug)),
		url.PathEscape(strings.TrimSpace(customerDomain)),
		url.PathEscape(strings.TrimSpace(formSlug)),
	)
}

func newModuleSettingsTabsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "settings-tabs <module_slug>",
		Short: "List available settings tabs for a module (read-only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/api/v1/external/modules/%s/settings/tabs", args[0])
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newModuleSettingsTabCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "settings-tab <module_slug> <tab_key>",
		Short: "Get values for one settings tab (read-only)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf(
				"/api/v1/external/modules/%s/settings/tab/%s",
				args[0],
				args[1],
			)
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newModuleListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}
			body, status, reqErr := client.Do("GET", "/api/v1/external/modules", nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newModuleGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <module_slug>",
		Short: "Get module details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/api/v1/external/modules/%s", args[0])
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newModuleOverviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "overview <module_slug> [section]",
		Short: "Get module overview (kpis, charts, tables)",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			section := "all"
			if len(args) == 2 {
				section = args[1]
			}
			if !isValidOverviewSection(section) {
				return fmt.Errorf("invalid section %q. allowed: kpis, charts, tables", section)
			}

			client, err := mustClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/api/v1/external/modules/%s/kpis", args[0])
			body, status, reqErr := client.Do("GET", path, nil)
			if reqErr != nil {
				return printResponse(cmd, body, status, reqErr)
			}

			rendered, renderErr := formatModuleOverviewOutput(body, args[0], section)
			if renderErr != nil {
				return printResponse(cmd, body, status, nil)
			}

			if rendered == "" {
				return printResponse(cmd, body, status, nil)
			}

			fmt.Fprintln(cmd.OutOrStdout(), rendered)
			return nil
		},
	}
}

func newModuleCollectionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "collections <module_slug>",
		Short: "Get module collections",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/api/v1/external/modules/%s/collections", args[0])
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newModuleCollectionCmd() *cobra.Command {
	var page int
	var sortField string
	var sortOrder string
	var filterField string
	var filterValue string

	cmd := &cobra.Command{
		Use:   "collection <module_slug> <collection_name>",
		Short: "Get module collection table",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			query := url.Values{}
			query.Set("page", fmt.Sprintf("%d", page))
			if sortField != "" {
				query.Set("sort_field", sortField)
			}
			if sortOrder != "" {
				query.Set("sort_order", sortOrder)
			}
			if filterField != "" {
				query.Set("filter_field", filterField)
			}
			if filterValue != "" {
				query.Set("filter_value", filterValue)
			}

			path := fmt.Sprintf(
				"/api/v1/external/modules/%s/collections/%s?%s",
				args[0],
				args[1],
				query.Encode(),
			)

			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}

	cmd.Flags().IntVar(&page, "page", 1, "page number")
	cmd.Flags().StringVar(&sortField, "sort-field", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort order: asc|desc")
	cmd.Flags().StringVar(&filterField, "filter-field", "", "filter field")
	cmd.Flags().StringVar(&filterValue, "filter-value", "", "filter value")

	return cmd
}

func newModuleDataCmd() *cobra.Command {
	var page int
	var sortField string
	var sortOrder string
	var filterField string
	var filterValue string

	cmd := &cobra.Command{
		Use:   "data <module_slug> <collection_name>",
		Short: "Get module data from /data endpoint",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			query := url.Values{}
			query.Set("page", fmt.Sprintf("%d", page))
			if sortField != "" {
				query.Set("sort_field", sortField)
			}
			if sortOrder != "" {
				query.Set("sort_order", sortOrder)
			}
			if filterField != "" {
				query.Set("filter_field", filterField)
			}
			if filterValue != "" {
				query.Set("filter_value", filterValue)
			}

			path := fmt.Sprintf(
				"/api/v1/external/modules/%s/data/%s?%s",
				args[0],
				args[1],
				query.Encode(),
			)

			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}

	cmd.Flags().IntVar(&page, "page", 1, "page number")
	cmd.Flags().StringVar(&sortField, "sort-field", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort order: asc|desc")
	cmd.Flags().StringVar(&filterField, "filter-field", "", "filter field")
	cmd.Flags().StringVar(&filterValue, "filter-value", "", "filter value")

	return cmd
}

func newModuleUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload <module_slug> <collection_name> <file_path_or_pattern...>",
		Short: "Upload one or more files to a module collection",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			uploadPaths, err := resolveUploadPaths(args[2:])
			if err != nil {
				return err
			}

			total := len(uploadPaths)
			success := 0
			failed := make([]string, 0)

			path := fmt.Sprintf(
				"/api/v1/external/modules/%s/collections/%s/inputs/file",
				args[0],
				args[1],
			)

			for index, filePath := range uploadPaths {
				fmt.Fprintf(cmd.OutOrStdout(), "[%d/%d] uploading: %s\n", index+1, total, filePath)

				content, readErr := os.ReadFile(filePath)
				if readErr != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "FAIL %s -> %v\n", filePath, readErr)
					failed = append(failed, filePath)
					continue
				}

				payload := map[string]string{
					"file_name": filepath.Base(filePath),
					"content":   base64.StdEncoding.EncodeToString(content),
				}

				b, marshalErr := json.Marshal(payload)
				if marshalErr != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "FAIL %s -> %v\n", filePath, marshalErr)
					failed = append(failed, filePath)
					continue
				}

				_, _, reqErr := client.Do("POST", path, b)
				if reqErr != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "FAIL %s -> %v\n", filePath, reqErr)
					failed = append(failed, filePath)
					continue
				}

				success++
				fmt.Fprintf(cmd.OutOrStdout(), "OK   %s\n", filePath)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Uploaded: %d/%d\n", success, total)
			if len(failed) > 0 {
				return fmt.Errorf("%d upload(s) failed", len(failed))
			}

			return nil
		},
	}
}

func hasGlobPattern(value string) bool {
	return strings.ContainsAny(value, "*?[")
}

func resolveUploadPaths(inputs []string) ([]string, error) {
	resolved := make([]string, 0)
	seen := make(map[string]struct{})

	for _, input := range inputs {
		if hasGlobPattern(input) {
			matches, err := filepath.Glob(input)
			if err != nil {
				return nil, fmt.Errorf("invalid glob pattern %q: %w", input, err)
			}
			if len(matches) == 0 {
				return nil, fmt.Errorf("no files matched pattern %q", input)
			}

			for _, match := range matches {
				info, err := os.Stat(match)
				if err != nil {
					return nil, fmt.Errorf("cannot access %q: %w", match, err)
				}
				if info.IsDir() {
					continue
				}
				if _, exists := seen[match]; exists {
					continue
				}
				seen[match] = struct{}{}
				resolved = append(resolved, match)
			}
			continue
		}

		info, err := os.Stat(input)
		if err != nil {
			return nil, fmt.Errorf("cannot access %q: %w", input, err)
		}
		if info.IsDir() {
			return nil, fmt.Errorf("path %q is a directory; provide files or file patterns", input)
		}
		if _, exists := seen[input]; exists {
			continue
		}
		seen[input] = struct{}{}
		resolved = append(resolved, input)
	}

	if len(resolved) == 0 {
		return nil, fmt.Errorf("no files to upload")
	}

	sort.Strings(resolved)
	return resolved, nil
}

func newModuleWebhookCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "webhook <module_slug> <collection_name> <payload_json>",
		Short: "Send webhook payload to a module collection",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf(
				"/api/v1/external/modules/%s/collections/%s/inputs/webhook",
				args[0],
				args[1],
			)
			body, status, reqErr := client.Do("POST", path, []byte(args[2]))
			return printResponse(cmd, body, status, reqErr)
		},
	}
}
