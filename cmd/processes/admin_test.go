package processes

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewAdminCmdContainsExpectedResources(t *testing.T) {
	cmd := newAdminCmd()
	if cmd == nil {
		t.Fatal("expected admin command")
	}

	expected := map[string]bool{
		"users":       false,
		"customers":   false,
		"credentials": false,
		"modules":     false,
		"tickets":     false,
		"logs":        false,
		"ai-requests": false,
		"usage":       false,
		"changelog":   false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected admin subcommand %s", name)
		}
	}
}

func TestAdminUsageContainsListAndSummary(t *testing.T) {
	cmd := newAdminUsageCmd()
	if cmd == nil {
		t.Fatal("expected admin usage command")
	}

	expected := map[string]bool{
		"list":    false,
		"summary": false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected usage subcommand %s", name)
		}
	}
}

func TestAdminUsersContainsCRUDSubset(t *testing.T) {
	cmd := newAdminUsersCmd()
	if cmd == nil {
		t.Fatal("expected admin users command")
	}

	expected := map[string]bool{
		"list":   false,
		"get":    false,
		"create": false,
		"update": false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected users subcommand %s", name)
		}
	}
}

func TestAdminTicketsContainsCRUDSubset(t *testing.T) {
	cmd := newAdminTicketsCmd()
	if cmd == nil {
		t.Fatal("expected admin tickets command")
	}

	expected := map[string]bool{
		"list":   false,
		"get":    false,
		"create": false,
		"update": false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected tickets subcommand %s", name)
		}
	}
}

func TestAdminCreateUpdateCommandsAreHidden(t *testing.T) {
	tests := []struct {
		name string
		cmd  *cobra.Command
	}{
		{name: "users create", cmd: newAdminUsersCreateCmd()},
		{name: "users update", cmd: newAdminUsersUpdateCmd()},
		{name: "customers create", cmd: newAdminCustomersCreateCmd()},
		{name: "customers update", cmd: newAdminCustomersUpdateCmd()},
		{name: "credentials create", cmd: newAdminCredentialsCreateCmd()},
		{name: "credentials update", cmd: newAdminCredentialsUpdateCmd()},
	}

	for _, tc := range tests {
		if tc.cmd == nil {
			t.Fatalf("expected command for %s", tc.name)
		}
		if !tc.cmd.Hidden {
			t.Fatalf("expected command %s to be hidden", tc.name)
		}
	}
}
