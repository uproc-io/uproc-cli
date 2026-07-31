package processes

import "testing"

func TestNewInvoiceCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newSalesInvoicesCmd()
	if cmd == nil {
		t.Fatal("expected invoice command")
	}

	expected := map[string]bool{
		"list":       false,
		"issue":      false,
		"rectify":    false,
		"send":       false,
		"get-pdf":    false,
		"add":        false,
		"list-lines": false,
		"update":     false,
		"delete":     false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected invoice subcommand %s", name)
		}
	}
}
