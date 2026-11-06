package processes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAdminAuditContainsListShowCapture(t *testing.T) {
	cmd := newAdminAuditCmd()
	if cmd == nil {
		t.Fatal("expected admin audit command")
	}

	expected := map[string]bool{"list": false, "show": false, "capture": false}
	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Fatalf("expected audit subcommand %s", name)
		}
	}
}

func TestAuditNoResultMatchesAuditAfter(t *testing.T) {
	detail := map[string]any{
		"after": map[string]any{"ok": true, "matched": false},
	}
	if !auditNoResult(detail) {
		t.Fatal("expected no-result when after.matched=false")
	}
}

func TestAuditNoResultMatchesRunReason(t *testing.T) {
	detail := map[string]any{
		"after": map[string]any{"ok": true, "matched": true},
		"run": map[string]any{
			"result": `{"status":"completed","reason":"no_result","message":{}}`,
		},
	}
	if !auditNoResult(detail) {
		t.Fatal("expected no-result from run reason")
	}
}

func TestAuditNoResultFalseForMatchedResult(t *testing.T) {
	detail := map[string]any{
		"after": map[string]any{"ok": true, "matched": true},
		"run": map[string]any{
			"free_reason": "",
			"result":      `{"status":"completed","reason":"valid"}`,
		},
	}
	if auditNoResult(detail) {
		t.Fatal("did not expect no-result")
	}
}

func TestWriteAuditOccurrenceRedactsPII(t *testing.T) {
	outDir := t.TempDir()
	detail := map[string]any{
		"id":          float64(42),
		"created_at":  "2026-11-06T10:00:00+00:00",
		"trace_id":    "abc123def456",
		"entity_id":   "get-company-by-name",
		"module_slug": "data-process",
		"after":       map[string]any{"ok": true, "matched": false},
		"run": map[string]any{
			"execution_mode": "native",
			"free_reason":    "no_result",
			"status":         "completed",
			"params":         `{"name":"Mehmet-company","country":"Germany","email":"nessim.djerboua@everreal.co","ip":"20.52.126.5"}`,
			"result":         `{"status":"completed","reason":"no_result","message":{}}`,
		},
	}

	path, existed, err := writeAuditOccurrence(outDir, detail, true, false)
	if err != nil {
		t.Fatalf("write: %v", err)
	}
	if existed {
		t.Fatal("did not expect existing file")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	text := string(content)
	if strings.Contains(text, "nessim.djerboua@everreal.co") {
		t.Fatal("email was not redacted")
	}
	if strings.Contains(text, "20.52.126.5") {
		t.Fatal("IP was not redacted")
	}
	if !strings.Contains(text, "get-company-by-name") || !strings.Contains(text, "Mehmet-company") {
		t.Fatal("expected tool and params in occurrence")
	}
	if !strings.Contains(text, "## Params") || !strings.Contains(text, "Audit `after`") {
		t.Fatal("expected occurrence sections")
	}

	// Second write must dedupe by file name.
	_, existedAgain, err := writeAuditOccurrence(outDir, detail, true, false)
	if err != nil {
		t.Fatalf("second write: %v", err)
	}
	if !existedAgain {
		t.Fatal("expected existing occurrence to be skipped")
	}
}

func TestWriteAuditOccurrenceDryRunDoesNotWrite(t *testing.T) {
	outDir := t.TempDir()
	detail := map[string]any{
		"id":         float64(7),
		"created_at": "2026-11-06T10:00:00+00:00",
		"trace_id":   "trace7",
		"entity_id":  "get-company-by-name",
		"after":      map[string]any{"matched": false},
	}
	path, existed, err := writeAuditOccurrence(outDir, detail, true, true)
	if err != nil {
		t.Fatalf("dry-run: %v", err)
	}
	if existed {
		t.Fatal("unexpected existing file")
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Fatal("dry-run should not create the file")
	}
	if filepath.Dir(path) == outDir {
		t.Fatal("expected tool/occurrences folder in path")
	}
}
