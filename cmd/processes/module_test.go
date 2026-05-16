package processes

import "testing"

func TestNewModuleCmdContainsSettingsReadOnlyCommands(t *testing.T) {
	cmd := newModuleCmd()
	if cmd == nil {
		t.Fatal("expected module command")
	}

	expected := map[string]bool{
		"list":               false,
		"get":                false,
		"overview":           false,
		"collections":        false,
		"collection":         false,
		"data":               false,
		"settings-tabs":      false,
		"settings-tab":       false,
		"upload":             false,
		"webhook":            false,
		"submit-public-form": false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected module subcommand %s", name)
		}
	}
}
