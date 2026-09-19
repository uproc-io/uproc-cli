package admin

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"bizzmod-cli/internal/api"
	"bizzmod-cli/internal/config"
)

// newLicenseCmd creates the `uproc admin license` command group.
func NewLicenseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "license",
		Short: "Manage and verify instance license",
	}
	cmd.AddCommand(newLicenseStatusCmd())
	cmd.AddCommand(newLicenseCheckCmd())
	cmd.AddCommand(newLicenseForceCmd())
	cmd.AddCommand(newLicenseUnforceCmd())
	return cmd
}

// ── license status ────────────────────────────────────────────────

func newLicenseStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show license status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config load error: %w", err)
			}
			client := api.NewClient(cfg)

			licenseKey := os.Getenv("UPROC_LICENSE_KEY")
			if licenseKey == "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "No UPROC_LICENSE_KEY set in environment.\n")
				return nil
			}

			body, _ := json.Marshal(map[string]any{"license_key": licenseKey})
			response, status, err := client.Do("POST", "/api/v1/platform/license/validate", body)
			if err != nil || status != 200 {
				return printResponse(cmd, response, status, err)
			}

			var result struct {
				Valid        bool    `json:"valid"`
				Forced       bool    `json:"forced"`
				ForcedReason string  `json:"forced_reason,omitempty"`
				Reason       string  `json:"reason,omitempty"`
				CustomerName string  `json:"customer_name"`
				ExpiresAt    *string `json:"expires_at,omitempty"`
				ModulesActive int   `json:"modules_active"`
			}
			if err := json.Unmarshal(response, &result); err != nil {
				return fmt.Errorf("json decode error: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "License Status for %s:\n", result.CustomerName)
			fmt.Fprintf(cmd.OutOrStdout(), "  Valid:        %v\n", result.Valid)
			fmt.Fprintf(cmd.OutOrStdout(), "  Forced:       %v\n", result.Forced)
			fmt.Fprintf(cmd.OutOrStdout(), "  Reason:       %s\n", result.Reason)
			if result.ExpiresAt != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "  Expires At:   %s\n", *result.ExpiresAt)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  Modules:      %d\n", result.ModulesActive)
			return nil
		},
	}
}

// ── license check (force immediate) ───────────────────────────────

func newLicenseCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Force immediate license check against cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config load error: %w", err)
			}
			client := api.NewClient(cfg)

			licenseKey := os.Getenv("UPROC_LICENSE_KEY")
			if licenseKey == "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "No UPROC_LICENSE_KEY set in environment.\n")
				return nil
			}

			body, _ := json.Marshal(map[string]any{"license_key": licenseKey})
			response, status, err := client.Do("POST", "/api/v1/platform/license/validate", body)
			if err != nil || status != 200 {
				return printResponse(cmd, response, status, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "License check completed (response logged above).\n")
			return printResponse(cmd, response, status, nil)
		},
	}
}

// ── license force ─────────────────────────────────────────────────

func newLicenseForceCmd() *cobra.Command {
	var customerID int
	var reason string

	cmd := &cobra.Command{
		Use:   "force",
		Short: "Force license validation (superadmin only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config load error: %w", err)
			}
			client := api.NewClient(cfg)

			body, _ := json.Marshal(map[string]any{
				"customer_id": customerID,
				"force":       true,
				"reason":      reason,
			})
			response, status, err := client.Do("PATCH", "/api/v1/external/admin/customers/"+fmt.Sprintf("%d", customerID)+"/license-force", body)
			if err != nil || status != 200 {
				return printResponse(cmd, response, status, err)
			}

			var result struct {
				Success bool   `json:"success"`
				Forced  bool   `json:"forced"`
				Reason  string `json:"reason"`
			}
			if err := json.Unmarshal(response, &result); err == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "License forced for customer %d: valid=%v reason=%s\n", customerID, result.Forced, result.Reason)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&customerID, "customer-id", 0, "Customer ID to force")
	cmd.Flags().StringVar(&reason, "reason", "", "Reason for force (required)")
	_ = cmd.MarkFlagRequired("customer-id")
	_ = cmd.MarkFlagRequired("reason")
	return cmd
}

// ── license unforce ───────────────────────────────────────────────

func newLicenseUnforceCmd() *cobra.Command {
	var customerID int
	var reason string

	cmd := &cobra.Command{
		Use:   "unforce",
		Short: "Unforce license (stop forcing validation, trust Stripe)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config load error: %w", err)
			}
			client := api.NewClient(cfg)

			body, _ := json.Marshal(map[string]any{
				"customer_id": customerID,
				"force":       false,
				"reason":      reason,
			})
			response, status, err := client.Do("PATCH", "/api/v1/external/admin/customers/"+fmt.Sprintf("%d", customerID)+"/license-force", body)
			if err != nil || status != 200 {
				return printResponse(cmd, response, status, err)
			}

			var result struct {
				Success bool   `json:"success"`
				Forced  bool   `json:"forced"`
				Reason  string `json:"reason"`
			}
			if err := json.Unmarshal(response, &result); err == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "License unforced for customer %d: valid=%v reason=%s\n", customerID, result.Forced, result.Reason)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&customerID, "customer-id", 0, "Customer ID to unforce")
	cmd.Flags().StringVar(&reason, "reason", "", "Reason for unforce (required)")
	_ = cmd.MarkFlagRequired("customer-id")
	_ = cmd.MarkFlagRequired("reason")
	return cmd
}

// ── helpers ───────────────────────────────────────────────────────

func printResponse(cmd *cobra.Command, response []byte, status int, err error) error {
	if status != 0 {
		fmt.Fprintf(cmd.ErrOrStderr(), "http %d", status)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), ": %v", err)
		}
		fmt.Fprintln(cmd.ErrOrStderr())
		if response != nil && len(response) > 0 {
			fmt.Fprintf(cmd.ErrOrStderr(), "%s\n", string(response))
		}
		return fmt.Errorf("http %d", status)
	}
	if response != nil && len(response) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(response))
	}
	return nil
}