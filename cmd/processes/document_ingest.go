package processes

import "github.com/spf13/cobra"

func newDocumentIngestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "document-ingest",
		Short: "Business verbs for document-ingest",
	}
	// List commands
	cmd.AddCommand(newCollectionListCmd("list-documents", "List document-ingest documents", "document-ingest", "documents"))
	cmd.AddCommand(newCollectionListCmd("list-uploads", "List document-ingest uploads", "document-ingest", "uploads"))
	cmd.AddCommand(newCollectionListCmd("list-emails", "List document-ingest emails", "document-ingest", "emails"))
	// Get commands
	cmd.AddCommand(newDocumentIngestGetCmd("document", "documents", "document_id"))
	cmd.AddCommand(newDocumentIngestGetCmd("email", "emails", "id"))
	cmd.AddCommand(newDocumentIngestGetCmd("upload", "uploads", "id"))
	// Action commands (reprocess requires user session auth, not API key —
	// added later when CLI supports login flow)
	// cmd.AddCommand(newDocumentIngestReprocessCmd())
	return cmd
}
