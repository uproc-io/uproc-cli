package processes

import "testing"

func TestNewFormsCmdContainsSubmitPublic(t *testing.T) {
	cmd := newFormsCmd()
	if cmd == nil {
		t.Fatal("expected forms command")
	}

	expected := map[string]bool{
		"submit-public":             false,
		"publish":                   false,
		"archive":                   false,
		"archive-submission":        false,
		"restore":                   false,
		"mark-submission-processed": false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected forms subcommand %s", name)
		}
	}
}

func TestRunFormsActionRejectsInvalidIDs(t *testing.T) {
	cmd := newFormsCmd()
	if err := runFormsAction(cmd, "publish", "form_id", "0"); err == nil {
		t.Fatal("expected error for invalid form_id")
	}
	if err := runFormsAction(cmd, "mark_submission_processed", "submission_id", "abc"); err == nil {
		t.Fatal("expected error for invalid submission_id")
	}
}
