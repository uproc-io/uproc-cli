package processes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newDocumentIngestGetCmd(cmdSuffix, collection, idField string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-" + cmdSuffix + " <id>",
		Short: "Show full details for a " + cmdSuffix + " by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parsePositiveIntArg(idField, args[0])
			if err != nil {
				return err
			}
			client, err := mustClient()
			if err != nil {
				return err
			}

			// Use the external list endpoint with filter params.
			// The backend only allows a fixed set of filter fields
			// (see module_common_ops.py ALLOWED_FILTER_FIELDS); "id" isn't
			// included, so we filter client-side from the returned rows.
			path := fmt.Sprintf(
				"/api/v1/external/modules/document-ingest/data/%s?filter_field=id&filter_value=%d&per_page=1",
				collection,
				id,
			)

			body, status, reqErr := client.Do("GET", path, nil)
			if reqErr != nil {
				return printResponse(cmd, body, status, reqErr)
			}

			// Parse and filter client-side: return only the row whose ID matches
			return printFilteredResponse(cmd, body, status, cmdSuffix, idField, id)
		},
	}
	return cmd
}

// printFilteredResponse parses the API response and renders only the row
// whose "id" or "document_id" matches the requested ID.
func printFilteredResponse(cmd *cobra.Command, body []byte, status int, cmdSuffix, idField string, id int) error {
	var parsed any
	if err := json.Unmarshal(body, &parsed); err != nil {
		return printResponse(cmd, body, status, err)
	}

	rootMap, ok := parsed.(map[string]any)
	if !ok {
		return printResponse(cmd, body, status, nil)
	}

	var rows []any
	if dataMap, ok := rootMap["data"].(map[string]any); ok {
		if r, ok := dataMap["rows"]; ok {
			if rList, ok := r.([]any); ok {
				rows = rList
			}
		}
	}

	// Try top-level rows if no data wrapper
	if len(rows) == 0 {
		if r, ok := rootMap["rows"]; ok {
			if rList, ok := r.([]any); ok {
				rows = rList
			}
		}
	}

	// Find matching row by id or document_id
	for _, row := range rows {
		rowMap, ok := row.(map[string]any)
		if !ok {
			continue
		}
		// Check id field
		if idVal, ok := rowMap["id"]; ok {
			if fmt.Sprintf("%v", idVal) == fmt.Sprintf("%d", id) {
				// Return as a single object (triggers map rendering, not table)
				out, _ := json.Marshal(map[string]any{
					"data": rowMap,
				})
				return printResponse(cmd, out, status, nil)
			}
		}
		// Check document_id (documents collection)
		if idField == "document_id" {
			if didVal, ok := rowMap["document_id"]; ok {
				if fmt.Sprintf("%v", didVal) == fmt.Sprintf("%d", id) {
					out, _ := json.Marshal(map[string]any{
						"data": rowMap,
					})
					return printResponse(cmd, out, status, nil)
				}
			}
		}
	}

	fmt.Fprintln(cmd.OutOrStdout(), "No "+cmdSuffix+" found with ID "+fmt.Sprintf("%d", id))
	return nil
}

func newDocumentIngestReprocessCmd() *cobra.Command {
	var mode string

	cmd := &cobra.Command{
		Use:   "reprocess <document_id>",
		Short: "Reprocess a document-ingest document (re-run AI extraction)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parsePositiveIntArg("document_id", args[0])
			if err != nil {
				return err
			}
			if mode == "" {
				mode = "auto"
			}
			return runModuleAction(cmd, "document-ingest", "reprocess", map[string]any{
				"document_id": id,
				"mode":        mode,
			})
		},
	}

	cmd.Flags().StringVarP(&mode, "mode", "m", "auto", "Reprocess mode: auto, full, or fields-only")
	return cmd
}

func newDocumentIngestUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload <file_path_or_pattern...>",
		Short: "Upload a PDF document to document-ingest for extraction",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			uploadPaths, err := resolveUploadPaths(args)
			if err != nil {
				return err
			}

			total := len(uploadPaths)
			success := 0
			failed := make([]string, 0)

			path := "/api/v1/external/modules/document-ingest/collections/uploads/inputs/file"

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

				body, _, reqErr := client.Do("POST", path, b)
				if reqErr != nil {
					fmt.Fprintf(cmd.OutOrStdout(), "FAIL %s -> %v\n", filePath, reqErr)
					failed = append(failed, filePath)
					continue
				}

				// Check success in response
				var resp map[string]any
				if json.Unmarshal(body, &resp) == nil {
					if successVal, ok := resp["success"].(bool); ok && !successVal {
						msg := fmt.Sprintf("%v", resp["message"])
						fmt.Fprintf(cmd.OutOrStdout(), "FAIL %s -> %s\n", filePath, msg)
						failed = append(failed, filePath)
						continue
					}
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