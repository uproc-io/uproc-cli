package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestPrintFriendlyProfileRendersKeyValues(t *testing.T) {
	body := []byte(`{
		"success": true,
		"message": "ok",
		"data": {
			"first_name": "Miquel",
			"last_name": "Colomer",
			"email": "mcolomer@killia.com",
			"role": "superadmin",
			"language": "es",
			"customer": "Killia",
			"customer_id": 1,
			"default_app": "dashboard",
			"modules": ["approval-management", "data-process"],
			"signature": "hi"
		}
	}`)

	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)

	if err := printFriendlyProfile(cmd, body); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{
		"Connected profile",
		"Miquel Colomer",
		"mcolomer@killia.com",
		"superadmin",
		"Killia (#1)",
		"Modules (2)",
		"approval-management, data-process",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in output:\n%s", want, out)
		}
	}
	if strings.Contains(out, "{") {
		t.Fatalf("friendly output should not contain raw JSON: %s", out)
	}
}

func TestPrintFriendlyProfileReturnsServerError(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetOut(&bytes.Buffer{})
	body := []byte(`{"success": false, "message": "user.not_found"}`)

	err := printFriendlyProfile(cmd, body)
	if err == nil || !strings.Contains(err.Error(), "user.not_found") {
		t.Fatalf("expected user.not_found error, got %v", err)
	}
}
