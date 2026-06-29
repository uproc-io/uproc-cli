package processes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func newAdminCmd() *cobra.Command {
	adminCmd := &cobra.Command{
		Use:   "admin",
		Short: "Admin commands for external API",
	}

	adminCmd.AddCommand(newAdminUsersCmd())
	adminCmd.AddCommand(newAdminCustomersCmd())
	adminCmd.AddCommand(newAdminCredentialsCmd())
	adminCmd.AddCommand(newAdminModulesCmd())
	adminCmd.AddCommand(newAdminTicketsCmd())
	adminCmd.AddCommand(newAdminLogsCmd())
	adminCmd.AddCommand(newAdminAIRequestsCmd())
	adminCmd.AddCommand(newAdminUsageCmd())
	adminCmd.AddCommand(newAdminChangelogCmd())

	return adminCmd
}

func newAdminUsersCmd() *cobra.Command {
	usersCmd := &cobra.Command{
		Use:   "users",
		Short: "Manage admin users",
	}

	usersCmd.AddCommand(newAdminUsersListCmd())
	usersCmd.AddCommand(newAdminUsersGetCmd())
	usersCmd.AddCommand(newAdminUsersCreateCmd())
	usersCmd.AddCommand(newAdminUsersUpdateCmd())
	usersCmd.AddCommand(newAdminUsersTemplateCmd())
	usersCmd.AddCommand(newAdminUsersBulkCreateCmd())
	usersCmd.AddCommand(newAdminUsersResetPasswordCmd())

	return usersCmd
}

func newAdminUsersListCmd() *cobra.Command {
	var customerID int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List admin users",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/api/v1/external/admin/users"
			if customerID > 0 {
				path = fmt.Sprintf("%s?customer_id=%d", path, customerID)
			}
			return runAdminListWithContract(cmd, "users", "GET", path)
		},
	}

	cmd.Flags().IntVar(&customerID, "customer-id", 0, "optional customer id filter")
	return cmd
}

func newAdminUsersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <user_id>",
		Short: "Get admin user details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/v1/external/admin/users/%s", args[0])
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newAdminUsersCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "create",
		Short:  "Create an admin user",
		Hidden: true,
		Args:   cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdminInteractiveContractAction(cmd, "users", "create")
		},
	}
}

func newAdminUsersUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "update",
		Short:  "Update an admin user",
		Hidden: true,
		Args:   cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdminInteractiveContractAction(cmd, "users", "update")
		},
	}
}

func newAdminCustomersCmd() *cobra.Command {
	customersCmd := &cobra.Command{
		Use:   "customers",
		Short: "Manage admin customers",
	}

	customersCmd.AddCommand(newAdminCustomersListCmd())
	customersCmd.AddCommand(newAdminCustomersGetCmd())
	customersCmd.AddCommand(newAdminCustomersCreateCmd())
	customersCmd.AddCommand(newAdminCustomersUpdateCmd())
	customersCmd.AddCommand(newAdminCustomersExportCmd())
	customersCmd.AddCommand(newAdminCustomersImportCmd())

	return customersCmd
}

func newAdminCustomersListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List admin customers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdminListWithContract(cmd, "customers", "GET", "/api/v1/external/admin/customers")
		},
	}
}

func newAdminCustomersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <customer_id>",
		Short: "Get admin customer details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/v1/external/admin/customers/%s", args[0])
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newAdminCustomersCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "create",
		Short:  "Create an admin customer",
		Hidden: true,
		Args:   cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdminInteractiveContractAction(cmd, "customers", "create")
		},
	}
}

func newAdminCustomersUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "update",
		Short:  "Update an admin customer",
		Hidden: true,
		Args:   cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdminInteractiveContractAction(cmd, "customers", "update")
		},
	}
}

func newAdminCredentialsCmd() *cobra.Command {
	credentialsCmd := &cobra.Command{
		Use:   "credentials",
		Short: "Manage admin credentials",
	}

	credentialsCmd.AddCommand(newAdminCredentialsListCmd())
	credentialsCmd.AddCommand(newAdminCredentialsGetCmd())
	credentialsCmd.AddCommand(newAdminCredentialsCreateCmd())
	credentialsCmd.AddCommand(newAdminCredentialsUpdateCmd())

	return credentialsCmd
}

func newAdminCredentialsListCmd() *cobra.Command {
	var customerID int
	var category string
	var credentialType string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List admin credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			queryParts := []string{}
			if customerID > 0 {
				queryParts = append(queryParts, fmt.Sprintf("customer_id=%d", customerID))
			}
			if strings.TrimSpace(category) != "" {
				queryParts = append(queryParts, fmt.Sprintf("category=%s", category))
			}
			if strings.TrimSpace(credentialType) != "" {
				queryParts = append(queryParts, fmt.Sprintf("type=%s", credentialType))
			}

			path := "/api/v1/external/admin/credentials"
			if len(queryParts) > 0 {
				path = fmt.Sprintf("%s?%s", path, strings.Join(queryParts, "&"))
			}
			return runAdminListWithContract(cmd, "credentials", "GET", path)
		},
	}

	cmd.Flags().IntVar(&customerID, "customer-id", 0, "optional customer id filter")
	cmd.Flags().StringVar(&category, "category", "", "optional credential category")
	cmd.Flags().StringVar(&credentialType, "type", "", "optional credential type")
	return cmd
}

func newAdminCredentialsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <credential_id>",
		Short: "Get admin credential details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/v1/external/admin/credentials/%s", args[0])
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newAdminCredentialsCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "create",
		Short:  "Create an admin credential",
		Hidden: true,
		Args:   cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdminInteractiveContractAction(cmd, "credentials", "create")
		},
	}
}

func newAdminCredentialsUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "update",
		Short:  "Update an admin credential",
		Hidden: true,
		Args:   cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdminInteractiveContractAction(cmd, "credentials", "update")
		},
	}
}

func newAdminModulesCmd() *cobra.Command {
	modulesCmd := &cobra.Command{
		Use:   "modules",
		Short: "Read admin modules",
	}

	modulesCmd.AddCommand(newAdminModulesListCmd())
	modulesCmd.AddCommand(newAdminModulesGetCmd())

	return modulesCmd
}

func newAdminModulesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List admin modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			body, status, reqErr := client.Do("GET", "/api/v1/external/admin/modules", nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newAdminModulesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <module_slug>",
		Short: "Get admin module details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/v1/external/admin/modules/%s", args[0])
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newAdminTicketsCmd() *cobra.Command {
	ticketsCmd := &cobra.Command{
		Use:   "tickets",
		Short: "Read and manage support tickets",
	}

	ticketsCmd.AddCommand(newAdminTicketsListCmd())
	ticketsCmd.AddCommand(newAdminTicketsGetCmd())
	ticketsCmd.AddCommand(newAdminTicketsCreateCmd())
	ticketsCmd.AddCommand(newAdminTicketsUpdateCmd())

	return ticketsCmd
}

func newAdminTicketsCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create support ticket",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdminInteractiveContractAction(cmd, "tickets", "create")
		},
	}
}

func newAdminTicketsUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update support ticket",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdminInteractiveContractAction(cmd, "tickets", "update")
		},
	}
}

func newAdminTicketsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List support tickets",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdminListWithContract(cmd, "tickets", "GET", "/api/v1/external/tickets/all")
		},
	}
}

func newAdminTicketsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <ticket_id>",
		Short: "Get support ticket details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/api/v1/external/tickets/%s/detail", args[0])
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newAdminLogsCmd() *cobra.Command {
	var moduleSlug string
	var level string
	var page int

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "List admin logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			if strings.TrimSpace(moduleSlug) == "" {
				return fmt.Errorf("--module-slug is required")
			}

			path := fmt.Sprintf(
				"/api/v1/external/admin/logs?module_slug=%s&level=%s&page=%d",
				moduleSlug,
				level,
				page,
			)
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}

	cmd.Flags().StringVar(&moduleSlug, "module-slug", "", "module slug to filter logs")
	cmd.Flags().StringVar(&level, "level", "all", "log level filter")
	cmd.Flags().IntVar(&page, "page", 1, "page number")
	return cmd
}

func newAdminAIRequestsCmd() *cobra.Command {
	var customerID int
	var moduleSlug string
	var page int
	var limit int

	cmd := &cobra.Command{
		Use:   "ai-requests",
		Short: "List AI request logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			queryParts := []string{fmt.Sprintf("page=%d", page), fmt.Sprintf("limit=%d", limit)}
			if customerID > 0 {
				queryParts = append(queryParts, fmt.Sprintf("customer_id=%d", customerID))
			}
			if strings.TrimSpace(moduleSlug) != "" {
				queryParts = append(queryParts, fmt.Sprintf("module_slug=%s", moduleSlug))
			}

			path := fmt.Sprintf("/api/v1/external/admin/ai-requests?%s", strings.Join(queryParts, "&"))
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}

	cmd.Flags().IntVar(&customerID, "customer-id", 0, "optional customer id filter")
	cmd.Flags().StringVar(&moduleSlug, "module-slug", "", "optional module slug filter")
	cmd.Flags().IntVar(&page, "page", 1, "page number")
	cmd.Flags().IntVar(&limit, "limit", 25, "items per page")
	return cmd
}

func newAdminUsageCmd() *cobra.Command {
	usageCmd := &cobra.Command{
		Use:   "usage",
		Short: "Read usage statistics",
	}

	usageCmd.AddCommand(newAdminUsageListCmd())
	usageCmd.AddCommand(newAdminUsageSummaryCmd())

	return usageCmd
}

func newAdminUsageListCmd() *cobra.Command {
	var customerID int
	var moduleSlug string
	var source string
	var fromDate string
	var toDate string
	var page int
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List usage events",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			queryParts := []string{fmt.Sprintf("page=%d", page), fmt.Sprintf("limit=%d", limit)}
			if customerID > 0 {
				queryParts = append(queryParts, fmt.Sprintf("customer_id=%d", customerID))
			}
			if strings.TrimSpace(moduleSlug) != "" {
				queryParts = append(queryParts, fmt.Sprintf("module_slug=%s", moduleSlug))
			}
			if strings.TrimSpace(source) != "" && source != "all" {
				queryParts = append(queryParts, fmt.Sprintf("source=%s", source))
			}
			if strings.TrimSpace(fromDate) != "" {
				queryParts = append(queryParts, fmt.Sprintf("from_date=%s", fromDate))
			}
			if strings.TrimSpace(toDate) != "" {
				queryParts = append(queryParts, fmt.Sprintf("to_date=%s", toDate))
			}

			path := fmt.Sprintf("/api/v1/external/admin/usage?%s", strings.Join(queryParts, "&"))
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}

	cmd.Flags().IntVar(&customerID, "customer-id", 0, "optional customer id filter")
	cmd.Flags().StringVar(&moduleSlug, "module-slug", "", "optional module slug filter")
	cmd.Flags().StringVar(&source, "source", "all", "optional source filter: all, api, cli, mcp")
	cmd.Flags().StringVar(&fromDate, "from-date", "", "optional from date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&toDate, "to-date", "", "optional to date (YYYY-MM-DD)")
	cmd.Flags().IntVar(&page, "page", 1, "page number")
	cmd.Flags().IntVar(&limit, "limit", 25, "items per page")
	return cmd
}

func newAdminUsageSummaryCmd() *cobra.Command {
	var customerID int
	var moduleSlug string
	var source string
	var fromDate string
	var toDate string

	cmd := &cobra.Command{
		Use:   "summary",
		Short: "Show usage summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			queryParts := []string{}
			if customerID > 0 {
				queryParts = append(queryParts, fmt.Sprintf("customer_id=%d", customerID))
			}
			if strings.TrimSpace(moduleSlug) != "" {
				queryParts = append(queryParts, fmt.Sprintf("module_slug=%s", moduleSlug))
			}
			if strings.TrimSpace(source) != "" && source != "all" {
				queryParts = append(queryParts, fmt.Sprintf("source=%s", source))
			}
			if strings.TrimSpace(fromDate) != "" {
				queryParts = append(queryParts, fmt.Sprintf("from_date=%s", fromDate))
			}
			if strings.TrimSpace(toDate) != "" {
				queryParts = append(queryParts, fmt.Sprintf("to_date=%s", toDate))
			}

			path := "/api/v1/external/admin/usage/summary"
			if len(queryParts) > 0 {
				path = fmt.Sprintf("%s?%s", path, strings.Join(queryParts, "&"))
			}
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}

	cmd.Flags().IntVar(&customerID, "customer-id", 0, "optional customer id filter")
	cmd.Flags().StringVar(&moduleSlug, "module-slug", "", "optional module slug filter")
	cmd.Flags().StringVar(&source, "source", "all", "optional source filter: all, api, cli, mcp")
	cmd.Flags().StringVar(&fromDate, "from-date", "", "optional from date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&toDate, "to-date", "", "optional to date (YYYY-MM-DD)")
	return cmd
}

func newAdminChangelogCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "changelog",
		Short: "Show platform changelog",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			body, status, reqErr := client.Do("GET", "/api/v1/external/admin/changelog", nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newAdminCustomersExportCmd() *cobra.Command {
	var customerID int
	var includeFiles, includeLogs bool

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export all customer data to a ZIP",
		Long: `Export ALL customer data (160+ tables, DM rows, optional files) as a ZIP file.
All data is exported synchronously and returned as base64 in the response.

Examples:
  uproc admin customers export --customer-id 51
  uproc admin customers export --customer-id 51 --include-files
  uproc admin customers export --customer-id 51 --include-logs`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID <= 0 {
				return fmt.Errorf("--customer-id is required")
			}

			fmt.Fprintf(cmd.ErrOrStderr(), "📦 Exporting customer %d (this may take a while)...\n", customerID)
			client, err := mustClient()
			if err != nil {
				return err
			}

			body, _ := json.Marshal(map[string]any{
				"name": "admin.customer.export_all",
				"arguments": map[string]any{
					"customer_id":   customerID,
					"include_files": includeFiles,
					"include_logs":  includeLogs,
				},
			})

			respBody, status, reqErr := client.Do("POST", "/api/v1/external/mcp/call", body)
			if reqErr != nil {
				return printResponse(cmd, respBody, status, reqErr)
			}
			if status != 200 {
				return printResponse(cmd, respBody, status, nil)
			}
			var resp struct {
				Success bool `json:"success"`
				Data    *struct {
					ZipBase64 string `json:"zip_base64"`
				} `json:"data"`
			}
			if err := json.Unmarshal(respBody, &resp); err != nil || !resp.Success || resp.Data == nil || resp.Data.ZipBase64 == "" {
				return printResponse(cmd, respBody, status, nil)
			}
			decoded, err := base64.StdEncoding.DecodeString(resp.Data.ZipBase64)
			if err != nil {
				return fmt.Errorf("failed to decode zip: %w", err)
			}
			filename := fmt.Sprintf("export-customer-%d-%s.zip", customerID, time.Now().Format("20060102-150405"))
			if err := os.WriteFile(filename, decoded, 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "✅ Export saved to %s (%d bytes)\n", filename, len(decoded))
			return nil
		},
	}

	cmd.Flags().IntVar(&customerID, "customer-id", 0, "Customer ID to export")
	cmd.Flags().BoolVar(&includeFiles, "include-files", false, "Include MinIO file blobs")
	cmd.Flags().BoolVar(&includeLogs, "include-logs", false, "Include log tables")
	return cmd
}

func newAdminCustomersImportCmd() *cobra.Command {
	var zipBase64 string
	var mode string
	var confirm bool

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import customer data from a ZIP (base64 or file path)",
		Long: `Import ALL customer data from a previously exported ZIP file.
The ZIP file content must be provided as a base64 string or via stdin.
Default mode is "upsert" (insert or update existing records).

Examples:
  uproc admin customers import --zip-base64 "$(cat export.zip | base64)"
  uproc admin customers import --zip-base64 "$(cat export.zip | base64)" --mode replace
  cat export.zip | base64 | xargs -0 uproc admin customers import --mode upsert`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if zipBase64 == "" {
				return fmt.Errorf("--zip-base64 is required (use file path with --file or pipe base64 content)")
			}
			if mode == "" {
				mode = "upsert"
			}

			client, err := mustClient()
			if err != nil {
				return err
			}

			body, _ := json.Marshal(map[string]any{
				"name": "admin.customer.import_all",
				"arguments": map[string]any{
					"zip_base64": zipBase64,
					"mode":       mode,
					"confirm":    confirm,
				},
			})

			respBody, status, reqErr := client.Do("POST", "/api/v1/external/mcp/call", body)
			return printResponse(cmd, respBody, status, reqErr)
		},
	}

	cmd.Flags().StringVar(&zipBase64, "zip-base64", "", "ZIP file content as base64 string")
	cmd.Flags().StringVar(&mode, "mode", "upsert", "Import mode: append, replace, or upsert")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm import execution (default false = preview only)")
	return cmd
}

func newAdminUsersTemplateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "template",
		Short: "Download CSV import template for users",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}
			body, _ := json.Marshal(map[string]any{
				"name":      "admin.users.import_template",
				"arguments": map[string]any{},
			})
			respBody, status, reqErr := client.Do("POST", "/api/v1/external/mcp/call", body)
			return printResponse(cmd, respBody, status, reqErr)
		},
	}
}

func newAdminUsersBulkCreateCmd() *cobra.Command {
	var csvFile string
	var confirm bool

	cmd := &cobra.Command{
		Use:   "bulk-create",
		Short: "Bulk create users from a CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := os.ReadFile(csvFile)
			if err != nil {
				return fmt.Errorf("cannot read file: %w", err)
			}
			client, err := mustClient()
			if err != nil {
				return err
			}
			body, _ := json.Marshal(map[string]any{
				"name": "admin.users.bulk_create",
				"arguments": map[string]any{
					"csv_base64": base64.StdEncoding.EncodeToString(content),
					"confirm":    confirm,
				},
			})
			respBody, status, reqErr := client.Do("POST", "/api/v1/external/mcp/call", body)
			return printResponse(cmd, respBody, status, reqErr)
		},
	}

	cmd.Flags().StringVarP(&csvFile, "file", "f", "", "CSV file path")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm import")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newAdminUsersResetPasswordCmd() *cobra.Command {
	var userID int

	cmd := &cobra.Command{
		Use:   "reset-password",
		Short: "Send password reset email to a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			if userID <= 0 {
				return fmt.Errorf("--user-id is required")
			}
			client, err := mustClient()
			if err != nil {
				return err
			}
			body, _ := json.Marshal(map[string]any{
				"name": "admin.users.send_reset_password",
				"arguments": map[string]any{
					"user_id": userID,
				},
			})
			respBody, status, reqErr := client.Do("POST", "/api/v1/external/mcp/call", body)
			return printResponse(cmd, respBody, status, reqErr)
		},
	}

	cmd.Flags().IntVar(&userID, "user-id", 0, "User ID to reset password for")
	return cmd
}
