package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"bizzmod-cli/internal/config"
)

func TestProfileAddSavesAndActivatesInteractiveProfile(t *testing.T) {
	config.SetConfigPath(filepath.Join(t.TempDir(), "config.yml"))
	t.Cleanup(func() { config.SetConfigPath("") })

	command := newProfileAddCmd()
	command.SetArgs([]string{"hospicare"})
	command.SetIn(strings.NewReader("https://api.example.test\nuser-api-key\n"))
	var output bytes.Buffer
	command.SetOut(&output)

	if err := command.Execute(); err != nil {
		t.Fatalf("profile add failed: %v", err)
	}

	active, err := config.GetActiveProfileName()
	if err != nil {
		t.Fatalf("get active profile: %v", err)
	}
	if active != "hospicare" {
		t.Fatalf("active profile = %q, want hospicare", active)
	}
	profile, err := config.Load()
	if err != nil {
		t.Fatalf("load profile: %v", err)
	}
	if profile.APIURL != "https://api.example.test" || profile.UserAPIKey != "user-api-key" {
		t.Fatalf("unexpected saved profile: %#v", profile)
	}
}
