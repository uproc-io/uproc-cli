package processes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func newDataManagementCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data-management",
		Short: "Business verbs for data-management",
	}
	cmd.AddCommand(newDMListTablesCmd())
	cmd.AddCommand(newDMCreateTableCmd())
	cmd.AddCommand(newDMListColumnsCmd())
	cmd.AddCommand(newDMEnqueueImportCmd())
	cmd.AddCommand(newDMGetUploadStatusCmd())
	cmd.AddCommand(newDMGetLatestUploadCmd())
	cmd.AddCommand(newDMAnalyzeNormalizationCmd())
	cmd.AddCommand(newDMExecuteNormalizationCmd())
	cmd.AddCommand(newDMListPlansCmd())
	cmd.AddCommand(newDMGetPlanCmd())
	cmd.AddCommand(newDMRegeneratePlanCmd())
	cmd.AddCommand(newDMGetRowMappingsCmd())
	cmd.AddCommand(newDMSmartImportPreviewCmd())
	cmd.AddCommand(newDMSmartImportExecuteCmd())
	cmd.AddCommand(newDMExportTransferCmd())
	cmd.AddCommand(newDMImportTransferCmd())
	return cmd
}

func newDMExportTransferCmd() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "export-transfer",
		Short: "Export all Data Management schema and rows to a portable ZIP",
		RunE: func(cmd *cobra.Command, args []string) error {
			if output == "" {
				return fmt.Errorf("--output is required")
			}
			client, err := mustClient()
			if err != nil {
				return err
			}
			response, status, reqErr := client.Do("GET", "/api/v1/external/data-management/transfer", nil)
			if reqErr != nil || status != 200 {
				return printResponse(cmd, response, status, reqErr)
			}
			var result struct {
				Success bool `json:"success"`
				Data    struct {
					ZipBase64 string          `json:"zip_base64"`
					Report    json.RawMessage `json:"report"`
				} `json:"data"`
			}
			if err := json.Unmarshal(response, &result); err != nil || !result.Success || result.Data.ZipBase64 == "" {
				return printResponse(cmd, response, status, nil)
			}
			archive, err := base64.StdEncoding.DecodeString(result.Data.ZipBase64)
			if err != nil {
				return fmt.Errorf("cannot decode transfer ZIP: %w", err)
			}
			if err := os.WriteFile(output, archive, 0600); err != nil {
				return fmt.Errorf("cannot write %s: %w", output, err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Data Management transfer saved to %s (%d bytes)\n", output, len(archive))

			if len(result.Data.Report) == 0 {
				return nil
			}
			reportPath := output + ".report.json"
			if err := os.WriteFile(reportPath, result.Data.Report, 0600); err != nil {
				return fmt.Errorf("cannot write %s: %w", reportPath, err)
			}
			var report struct {
				Relations struct {
					References         int `json:"references"`
					Resolved           int `json:"resolved"`
					Unresolved         int `json:"unresolved"`
					IgnoredNonPositive int `json:"ignored_non_positive"`
				} `json:"relations"`
			}
			_ = json.Unmarshal(result.Data.Report, &report)
			fmt.Fprintf(
				cmd.OutOrStdout(),
				"Relations: %d resolved, %d unresolved, %d ignored (non-positive)\n",
				report.Relations.Resolved,
				report.Relations.Unresolved,
				report.Relations.IgnoredNonPositive,
			)
			fmt.Fprintf(cmd.OutOrStdout(), "Export report saved to %s\n", reportPath)
			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "", "Portable ZIP output path")
	return cmd
}

func newDMImportTransferCmd() *cobra.Command {
	var file string
	var confirm bool

	cmd := &cobra.Command{
		Use:   "import-transfer",
		Short: "Preview or replace all Data Management schema and rows from a portable ZIP",
		Long:  "Without --confirm, shows a replacement preview. --confirm permanently replaces only the authenticated customer's Data Management tables.\n\nUploads via multipart/form-data (no base64 overhead); runs async, CLI polls until completion.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file is required")
			}

			client, err := mustClient()
			if err != nil {
				return err
			}

			// Read the ZIP file
			archive, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("cannot read %s: %w", file, err)
			}

			// Build multipart body
			pr, pw := io.Pipe()
			mw := multipart.NewWriter(pw)

			// io.Pipe is synchronous: every Write blocks until the reader
			// consumes it, so the multipart encoding must run concurrently
			// with the HTTP request that streams the pipe as its body.
			go func() {
				var writeErr error
				defer func() {
					if writeErr != nil {
						pw.CloseWithError(writeErr)
						return
					}
					pw.Close()
				}()

				part, err := mw.CreateFormFile("file", filepath.Base(file))
				if err != nil {
					writeErr = fmt.Errorf("cannot create form file: %w", err)
					return
				}
				if _, err := part.Write(archive); err != nil {
					writeErr = fmt.Errorf("cannot write file: %w", err)
					return
				}
				if err := mw.WriteField("confirm", fmt.Sprintf("%v", confirm)); err != nil {
					writeErr = fmt.Errorf("cannot write confirm: %w", err)
					return
				}
				writeErr = mw.Close()
			}()

			respBody, status, err := client.DoMultipart("POST", "/api/v1/external/data-management/transfer/import", mw.FormDataContentType(), pr)
			if err != nil || status != 200 {
				return fmt.Errorf("multipart upload failed (http %d): %s", status, string(respBody))
			}

			var enqueue struct {
				Success  bool   `json:"success"`
				UploadID string `json:"upload_id"`
				Status   string `json:"status"`
			}
			if err := json.Unmarshal(respBody, &enqueue); err != nil || !enqueue.Success {
				return fmt.Errorf("upload not queued: %s", string(respBody))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Import queued (upload %s). Polling status...\n", enqueue.UploadID)

			var final []byte
			for {
				statusResp, st, err := client.Do("GET", "/api/v1/external/data-management/transfer/import/status/"+enqueue.UploadID, nil)
				if err != nil || st != 200 {
					fmt.Fprintf(cmd.ErrOrStderr(), "Error polling: %v (status %d)\n", err, st)
					return printResponse(cmd, statusResp, st, err)
				}
				var poll struct {
					Success bool            `json:"success"`
					Status  string          `json:"status"`
					Error   string          `json:"error,omitempty"`
					Data    json.RawMessage `json:"data,omitempty"`
				}
				if err := json.Unmarshal(statusResp, &poll); err != nil {
					return fmt.Errorf("cannot decode status: %w", err)
				}
				if poll.Status == "done" {
					final = statusResp
					break
				}
				if poll.Status == "error" {
					fmt.Fprintf(cmd.ErrOrStderr(), "Import failed: %s\n", poll.Error)
					return printResponse(cmd, statusResp, 400, fmt.Errorf("import error: %s", poll.Error))
				}
				// Still processing – poll every 1s
				time.Sleep(1 * time.Second)
			}
			return printResponse(cmd, final, 200, nil)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "Portable Data Management ZIP")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "Replace all Data Management tables for the authenticated customer")
	return cmd
}

func newDMListTablesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-tables",
		Short: "List data-management tables",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModuleAction(cmd, "data-management", "list_tables", map[string]any{})
		},
	}
}

func newDMCreateTableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-table <item>",
		Short: "Create a data-management table",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			payload, err := parseJSONObjectArg("item", args[0])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "data-management", "create_table", payload)
		},
	}
}

func newDMListColumnsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-columns <entity_id>",
		Short: "List columns for a data-management table",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "list_columns", "entity_id", args[0])
		},
	}
}

func newDMEnqueueImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enqueue-import <collection_name> <file_name> <content>",
		Short: "Queue a file import into a data-management table",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModuleAction(cmd, "data-management", "enqueue_import", map[string]any{
				"collection_name": args[0],
				"file_name":       args[1],
				"content":         args[2],
			})
		},
	}
}

func newDMGetUploadStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-upload-status <upload_id>",
		Short: "Check the status of a queued data-management import",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "get_upload_status", "upload_id", args[0])
		},
	}
}

func newDMGetLatestUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-latest-upload",
		Short: "Get the latest data-management import status",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModuleAction(cmd, "data-management", "get_latest_upload", map[string]any{})
		},
	}
}

func newDMAnalyzeNormalizationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "analyze-normalization <entity_id> [entity_id...]",
		Short: "Analyze entities and generate AI normalization plan",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityIDs := make([]int, 0, len(args))
			for _, raw := range args {
				id, err := parsePositiveIntArg("entity_id", raw)
				if err != nil {
					return err
				}
				entityIDs = append(entityIDs, id)
			}
			return runModuleAction(cmd, "data-management", "analyze_normalization", map[string]any{
				"entity_ids": entityIDs,
			})
		},
	}
}

func newDMExecuteNormalizationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "execute-normalization <plan_id>",
		Short: "Execute an accepted normalization plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "execute_normalization", "plan_id", args[0])
		},
	}
}

func newDMListPlansCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-plans",
		Short: "List normalization plans",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModuleAction(cmd, "data-management", "list_plans", map[string]any{})
		},
	}
}

func newDMGetPlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-plan <plan_id>",
		Short: "Get a normalization plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "get_plan", "plan_id", args[0])
		},
	}
}

func newDMRegeneratePlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "regenerate-plan <plan_id>",
		Short: "Regenerate a normalization plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "regenerate_plan", "plan_id", args[0])
		},
	}
}

func newDMGetRowMappingsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-row-mappings <plan_id>",
		Short: "Get row-level change mappings for a normalization plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "get_row_mappings", "plan_id", args[0])
		},
	}
}

func newDMSmartImportPreviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "smart-import-preview <raw_entity_id>",
		Short: "Preview which normalized entity a raw upload matches",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "smart_import_preview", "raw_entity_id", args[0])
		},
	}
}

func newDMSmartImportExecuteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "smart-import-execute <raw_entity_id> <target_entity_id>",
		Short: "Execute smart import from raw upload to normalized entity",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			rawID, err := parsePositiveIntArg("raw_entity_id", args[0])
			if err != nil {
				return err
			}
			targetID, err := parsePositiveIntArg("target_entity_id", args[1])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "data-management", "smart_import_execute", map[string]any{
				"raw_entity_id":    rawID,
				"target_entity_id": targetID,
			})
		},
	}
}
