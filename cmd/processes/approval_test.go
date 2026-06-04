package processes

import "testing"

func TestNewApprovalCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newApprovalCmd()
	if cmd == nil {
		t.Fatal("expected approval command")
	}

	expected := map[string]bool{
		"approve":  false,
		"reject":   false,
		"reassign": false,
		"cancel":   false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected approval subcommand %s", name)
		}
	}
}
