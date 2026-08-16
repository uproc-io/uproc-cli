package update

import (
	"strings"
	"testing"
)

func TestNormalizeVersion(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"0.1.4", "v0.1.4"},
		{"v0.1.4", "v0.1.4"},
		{"v1.2.3-rc.1", "v1.2.3-rc.1"},
		{"dev", ""},
		{"", ""},
		{"latest", ""},
	}
	for _, c := range cases {
		if got := NormalizeVersion(c.in); got != c.want {
			t.Errorf("NormalizeVersion(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestParseChecksums(t *testing.T) {
	body := "  abc123  uproc.cli_0.1.4_darwin_arm64.tar.gz\n" +
		"def456  uproc.cli_0.1.4_darwin_amd64.tar.gz\n"
	sums := parseChecksums([]byte(body))
	if got := sums["uproc.cli_0.1.4_darwin_arm64.tar.gz"]; got != "abc123" {
		t.Errorf("unexpected hash for darwin_arm64: %q", got)
	}
	if got := sums["uproc.cli_0.1.4_darwin_amd64.tar.gz"]; got != "def456" {
		t.Errorf("unexpected hash for darwin_amd64: %q", got)
	}
	if len(sums) != 2 {
		t.Errorf("expected 2 entries, got %d", len(sums))
	}
}

func TestAssetSelection(t *testing.T) {
	rel := Release{
		TagName: "v0.1.4",
		Assets: []Asset{
			{Name: "checksums.txt", URL: "https://example.com/checksums.txt"},
			{Name: "uproc.cli_0.1.4_darwin_arm64.tar.gz", URL: "https://example.com/darwin-arm64"},
			{Name: "uproc.cli_0.1.4_darwin_amd64.tar.gz", URL: "https://example.com/darwin-amd64"},
			{Name: "uproc.cli_0.1.4_windows_amd64.zip", URL: "https://example.com/win-amd64"},
		},
	}
	u := &Updater{}

	cases := []struct {
		goos, goarch, want string
		wantErr            bool
	}{
		{"darwin", "arm64", "uproc.cli_0.1.4_darwin_arm64.tar.gz", false},
		{"darwin", "amd64", "uproc.cli_0.1.4_darwin_amd64.tar.gz", false},
		{"windows", "amd64", "uproc.cli_0.1.4_windows_amd64.zip", false},
		{"linux", "arm64", "", true},
	}
	for _, c := range cases {
		got, err := u.Asset(rel, c.goos, c.goarch)
		if c.wantErr {
			if err == nil {
				t.Errorf("Asset(%s/%s): expected error, got %+v", c.goos, c.goarch, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("Asset(%s/%s): unexpected error: %v", c.goos, c.goarch, err)
			continue
		}
		if got.Name != c.want {
			t.Errorf("Asset(%s/%s) = %q, want %q", c.goos, c.goarch, got.Name, c.want)
		}
	}
}

func TestDetectInstall(t *testing.T) {
	cases := []struct {
		path string
		want InstallMethod
	}{
		{"/opt/homebrew/Cellar/uproc/0.1.4/bin/uproc", InstallHomebrew},
		{"/usr/local/Cellar/uproc/0.1.4/bin/uproc", InstallHomebrew},
		{"/home/linuxbrew/.linuxbrew/Cellar/uproc/bin/uproc", InstallHomebrew},
		{"/Users/user/scoop/shims/uproc.exe", InstallScoop},
		{"/usr/local/bin/uproc", InstallStandalone},
		{"/tmp/uproc", InstallStandalone},
	}
	for _, c := range cases {
		if got := DetectInstall(c.path); got != c.want {
			t.Errorf("DetectInstall(%q) = %v, want %v", c.path, got, c.want)
		}
	}
}

func TestChecksumsNotFound(t *testing.T) {
	rel := Release{TagName: "v0.1.4", Assets: []Asset{{Name: "uproc.cli_0.1.4_darwin_arm64.tar.gz"}}}
	u := &Updater{}
	_, err := u.Checksums(rel)
	if err == nil || !strings.Contains(err.Error(), "checksums.txt not found") {
		t.Errorf("expected checksums.txt not found error, got %v", err)
	}
}
