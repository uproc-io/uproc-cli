package processes

import "testing"

func TestNewSupportCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newSupportCmd()
	if cmd == nil {
		t.Fatal("expected support command")
	}

	expected := map[string]bool{
		"create-ticket": false,
		"assign-ticket": false,
		"reply-ticket":  false,
		"mark-resolved": false,
		"close-ticket":  false,
		"reopen-ticket": false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected support subcommand %s", name)
		}
	}
}
