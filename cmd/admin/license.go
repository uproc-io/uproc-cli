package admin

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"bizzmod-cli/internal/api"
	"bizzmod-cli/internal/config"
)

// NewLicenseCmd creates the `uproc admin license` command group.
func NewLicenseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "license",
		Short: "Inspect and manage the instance license",
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
		Short: "Show the current license status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config load error: %w", err)
			}
			client := api.NewClient(cfg)
			if err := ensureRole(cmd, client, "license status", "admin", "superadmin"); err != nil {
				return err
			}

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
				Valid         bool    `json:"valid"`
				Forced        bool    `json:"forced"`
				ForcedReason  string  `json:"forced_reason,omitempty"`
				Reason        string  `json:"reason,omitempty"`
				CustomerName  string  `json:"customer_name"`
				ExpiresAt     *string `json:"expires_at,omitempty"`
				ModulesActive int     `json:"modules_active"`
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
		Short: "Force an immediate validation against the cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config load error: %w", err)
			}
			client := api.NewClient(cfg)
			if err := ensureRole(cmd, client, "license check", "admin", "superadmin"); err != nil {
				return err
			}

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
			return printResponse(cmd, response, status, nil)
		},
	}
}

// ── license force / unforce (superadmin only) ─────────────────────

func newLicenseForceCmd() *cobra.Command {
	var customerID int
	var reason string

	cmd := &cobra.Command{
		Use:   "force",
		Short: "Force a customer license valid, bypassing Stripe (superadmin only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLicenseForce(cmd, customerID, true, reason)
		},
	}
	cmd.Flags().IntVar(&customerID, "customer-id", 0, "Customer ID to force")
	cmd.Flags().StringVar(&reason, "reason", "", "Reason for forcing (required)")
	_ = cmd.MarkFlagRequired("customer-id")
	_ = cmd.MarkFlagRequired("reason")
	return cmd
}

func newLicenseUnforceCmd() *cobra.Command {
	var customerID int
	var reason string

	cmd := &cobra.Command{
		Use:   "unforce",
		Short: "Stop forcing a customer license, trust Stripe again (superadmin only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLicenseForce(cmd, customerID, false, reason)
		},
	}
	cmd.Flags().IntVar(&customerID, "customer-id", 0, "Customer ID to unforce")
	cmd.Flags().StringVar(&reason, "reason", "", "Reason for unforcing (required)")
	_ = cmd.MarkFlagRequired("customer-id")
	_ = cmd.MarkFlagRequired("reason")
	return cmd
}

func runLicenseForce(cmd *cobra.Command, customerID int, force bool, reason string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config load error: %w", err)
	}
	client := api.NewClient(cfg)
	if err := ensureRole(cmd, client, "license management", "superadmin"); err != nil {
		return err
	}

	body, _ := json.Marshal(map[string]any{"force": force, "reason": reason})
	path := fmt.Sprintf("/api/v1/external/admin/customers/%d/license-force", customerID)
	response, status, err := client.Do("PATCH", path, body)
	if err != nil || status != 200 {
		return printResponse(cmd, response, status, err)
	}

	var result struct {
		Success bool   `json:"success"`
		Forced  bool   `json:"forced"`
		Reason  string `json:"reason"`
	}
	if err := json.Unmarshal(response, &result); err == nil {
		action := "forced"
		if !force {
			action = "unforced"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "License %s for customer %d (forced=%v).\n", action, customerID, result.Forced)
	}
	return nil
}

// ── permissions ───────────────────────────────────────────────────

// ensureRole verifies the connected CLI user has one of the allowed roles.
func ensureRole(cmd *cobra.Command, client *api.Client, scope string, allowed ...string) error {
	body, status, err := client.Do("GET", "/api/v1/external/profile", nil)
	if err != nil || status != 200 {
		return fmt.Errorf("cannot verify permissions for %s (http %d)", scope, status)
	}

	var parsed struct {
		Data struct {
			Role string `json:"role"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Errorf("cannot parse profile response: %w", err)
	}

	role := strings.ToLower(strings.TrimSpace(parsed.Data.Role))
	for _, candidate := range allowed {
		if role == candidate {
			return nil
		}
	}
	return fmt.Errorf("permission denied: %s requires %s (your role: %s)", scope, strings.Join(allowed, " or "), role)
}

// ── helpers ───────────────────────────────────────────────────────

func printResponse(cmd *cobra.Command, response []byte, status int, err error) error {
	if status != 0 {
		fmt.Fprintf(cmd.ErrOrStderr(), "http %d", status)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), ": %v", err)
		}
		fmt.Fprintln(cmd.ErrOrStderr())
		if len(response) > 0 {
			fmt.Fprintf(cmd.ErrOrStderr(), "%s\n", string(response))
		}
		return fmt.Errorf("http %d", status)
	}
	if len(response) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(response))
	}
	return nil
}
