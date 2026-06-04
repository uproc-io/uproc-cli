package processes

import "testing"

func TestNewCandidateCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newCandidateCmd()
	if cmd == nil {
		t.Fatal("expected candidate command")
	}

	expected := map[string]bool{
		"create-profile":     false,
		"create-job-opening": false,
		"create-application": false,
		"move-stage":         false,
		"update-status":      false,
		"create-evaluation":  false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected candidate subcommand %s", name)
		}
	}
}
