package processes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"text/tabwriter"

	"bizzmod-cli/internal/config"
	"github.com/spf13/cobra"
)

var lookPathFn = exec.LookPath
var commandOutputFn = func(name string, args ...string) (string, error) {
	output, err := exec.Command(name, args...).CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

type updateCheckRow struct {
	ID       string
	Expected string
	Actual   string
	Status   string
}

func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update verification commands (dry-run only)",
	}

	cmd.AddCommand(newUpdateCheckCmd())
	return cmd
}

func newUpdateCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check <user_api_key>",
		Short: "Validate user server update readiness (dry-run only)",
		Long:  "Runs read-only checks against install plan and local server state. No deployment changes are executed.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			cfg.UserAPIKey = args[0]
			if err := config.Validate(cfg); err != nil {
				return err
			}

			client, err := mustClientFromConfig(cfg)
			if err != nil {
				return err
			}

			path := "/api/v1/external/install?dry_run=true"
			body, status, reqErr := client.Do("GET", path, nil)
			if err := printResponse(cmd, body, status, reqErr); err != nil {
				return err
			}

			requiredServices, requiredEnvVars := parseInstallRequirements(body)
			rows := runLocalUpdateChecks(cfg, requiredServices, requiredEnvVars)
			printUpdateCheckRows(cmd, rows)

			failed := 0
			for _, row := range rows {
				if row.Status == "fail" {
					failed++
				}
			}

			if failed > 0 {
				return fmt.Errorf("update check failed: %d check(s) did not pass", failed)
			}
			return nil
		},
	}

	return cmd
}

func parseInstallRequirements(payload []byte) ([]string, []string) {
	requiredServices := []string{"docker", "dokploy", "bizzmod-back", "bizzmod-front", "whisper-asr"}
	requiredEnv := []string{"WHISPER_SERVICE_URL", "WHISPER_SERVICE_TIMEOUT"}

	var parsed map[string]any
	if err := json.Unmarshal(payload, &parsed); err != nil {
		return requiredServices, requiredEnv
	}

	data, _ := parsed["data"].(map[string]any)
	install, _ := data["install"].(map[string]any)

	if services, ok := install["required_services"].([]any); ok && len(services) > 0 {
		requiredServices = make([]string, 0, len(services))
		for _, service := range services {
			value := strings.TrimSpace(fmt.Sprintf("%v", service))
			if value != "" {
				requiredServices = append(requiredServices, value)
			}
		}
	}

	if envItems, ok := install["required_env"].([]any); ok && len(envItems) > 0 {
		requiredEnv = []string{}
		for _, item := range envItems {
			entry, _ := item.(map[string]any)
			name := strings.TrimSpace(fmt.Sprintf("%v", entry["name"]))
			if name != "" {
				requiredEnv = append(requiredEnv, name)
			}
		}
	}

	sort.Strings(requiredServices)
	sort.Strings(requiredEnv)
	return requiredServices, requiredEnv
}

func runLocalUpdateChecks(cfg config.Config, requiredServices []string, requiredEnv []string) []updateCheckRow {
	rows := []updateCheckRow{}

	if _, err := lookPathFn("docker"); err == nil {
		rows = append(rows, updateCheckRow{ID: "docker", Expected: "binary installed", Actual: "available", Status: "pass"})
	} else {
		rows = append(rows, updateCheckRow{ID: "docker", Expected: "binary installed", Actual: err.Error(), Status: "fail"})
	}

	if _, err := lookPathFn("dokploy"); err == nil {
		rows = append(rows, updateCheckRow{ID: "dokploy", Expected: "binary installed", Actual: "available", Status: "pass"})
	} else {
		rows = append(rows, updateCheckRow{ID: "dokploy", Expected: "binary installed", Actual: err.Error(), Status: "fail"})
	}

	dockerPsOutput, dockerErr := commandOutputFn("docker", "ps", "--format", "{{.Names}}")
	containerNames := []string{}
	if dockerErr == nil && dockerPsOutput != "" {
		containerNames = strings.Split(dockerPsOutput, "\n")
	}

	for _, service := range requiredServices {
		trimmedService := strings.TrimSpace(service)
		if trimmedService == "" || trimmedService == "docker" || trimmedService == "dokploy" {
			continue
		}

		if dockerErr != nil {
			rows = append(rows, updateCheckRow{
				ID:       fmt.Sprintf("service:%s", trimmedService),
				Expected: "running container",
				Actual:   dockerErr.Error(),
				Status:   "fail",
			})
			continue
		}

		found := false
		for _, name := range containerNames {
			normalizedName := strings.TrimSpace(strings.ToLower(name))
			if normalizedName == "" {
				continue
			}
			if strings.Contains(normalizedName, strings.ToLower(trimmedService)) {
				found = true
				break
			}
		}

		if found {
			rows = append(rows, updateCheckRow{
				ID:       fmt.Sprintf("service:%s", trimmedService),
				Expected: "running container",
				Actual:   "running",
				Status:   "pass",
			})
		} else {
			rows = append(rows, updateCheckRow{
				ID:       fmt.Sprintf("service:%s", trimmedService),
				Expected: "running container",
				Actual:   "not running",
				Status:   "fail",
			})
		}
	}

	for _, variableName := range requiredEnv {
		name := strings.TrimSpace(variableName)
		if name == "" {
			continue
		}
		value := strings.TrimSpace(os.Getenv(name))
		if value == "" {
			rows = append(rows, updateCheckRow{
				ID:       fmt.Sprintf("env:%s", name),
				Expected: "set",
				Actual:   "missing",
				Status:   "fail",
			})
		} else {
			rows = append(rows, updateCheckRow{
				ID:       fmt.Sprintf("env:%s", name),
				Expected: "set",
				Actual:   "present",
				Status:   "pass",
			})
		}
	}

	rows = append(rows, checkHTTPAvailability("health:backend", cfg.APIURL))

	whisperURL := strings.TrimSpace(os.Getenv("WHISPER_SERVICE_URL"))
	if whisperURL != "" {
		rows = append(rows, checkHTTPAvailability("health:whisper", whisperURL))
	} else {
		rows = append(rows, updateCheckRow{
			ID:       "health:whisper",
			Expected: "HTTP 2xx/3xx",
			Actual:   "missing WHISPER_SERVICE_URL",
			Status:   "fail",
		})
	}

	return rows
}

func checkHTTPAvailability(id string, baseURL string) updateCheckRow {
	trimmedURL := strings.TrimSpace(baseURL)
	if trimmedURL == "" {
		return updateCheckRow{ID: id, Expected: "HTTP 2xx/3xx", Actual: "missing URL", Status: "fail"}
	}

	healthURL := strings.TrimRight(trimmedURL, "/")
	if id == "health:backend" {
		healthURL += "/health"
	} else if id == "health:whisper" {
		healthURL += "/docs"
	}

	resp, err := http.Get(healthURL)
	if err != nil {
		return updateCheckRow{ID: id, Expected: "HTTP 2xx/3xx", Actual: err.Error(), Status: "fail"}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return updateCheckRow{ID: id, Expected: "HTTP 2xx/3xx", Actual: fmt.Sprintf("HTTP %d", resp.StatusCode), Status: "pass"}
	}

	return updateCheckRow{ID: id, Expected: "HTTP 2xx/3xx", Actual: fmt.Sprintf("HTTP %d", resp.StatusCode), Status: "fail"}
}

func printUpdateCheckRows(cmd *cobra.Command, rows []updateCheckRow) {
	if len(rows) == 0 {
		return
	}

	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), "Local update checks (dry-run):")

	tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "CHECK\tEXPECTED\tACTUAL\tSTATUS")
	for _, row := range rows {
		_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", row.ID, row.Expected, row.Actual, row.Status)
	}
	_ = tw.Flush()
}
