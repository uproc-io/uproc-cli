package cmd

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"bizzmod-cli/internal/update"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

func newSelfUpdateCmd(currentVersion string) *cobra.Command {
	var checkOnly bool
	var targetVersion string
	var allowPre bool

	cmd := &cobra.Command{
		Use:   "self-update",
		Short: "Update the CLI to the latest released version",
		Long: `Updates the CLI to the latest released version from GitHub.

Checks the uproc.cli GitHub release for the current platform, verifies the
artifact checksum and replaces the running binary atomically. When the CLI was
installed via Homebrew or Scoop, prints the package-manager update command
instead of touching the binary.`,
		RunE: func(c *cobra.Command, args []string) error {
			return runSelfUpdate(c, currentVersion, targetVersion, allowPre, checkOnly)
		},
	}

	cmd.Flags().BoolVar(&checkOnly, "check", false, "only check whether an update is available")
	cmd.Flags().StringVar(&targetVersion, "version", "", "update to a specific version (e.g. v0.1.5)")
	cmd.Flags().BoolVar(&allowPre, "pre", false, "allow pre-release versions")

	return cmd
}

func runSelfUpdate(c *cobra.Command, currentVersion, targetVersion string, allowPre, checkOnly bool) error {
	updater := &update.Updater{HTTPClient: &http.Client{Timeout: 30 * time.Second}}

	var rel update.Release
	var err error
	switch {
	case targetVersion != "":
		rel, err = updater.ReleaseForVersion(targetVersion)
	case allowPre:
		rel, err = updater.LatestPrerelease()
	default:
		rel, err = updater.Latest()
	}
	if err != nil {
		return err
	}

	current := update.NormalizeVersion(currentVersion)
	if current == "" {
		return update.ErrUnknownVersion
	}

	target := rel.Version()
	if target == "" {
		return fmt.Errorf("release %q has no valid semver tag", rel.TagName)
	}

	if semver.Compare(current, target) >= 0 {
		if checkOnly {
			fmt.Fprintf(c.OutOrStdout(), "uproc %s is up to date (latest: %s)\n", current, target)
			return nil
		}
		fmt.Fprintf(c.OutOrStdout(), "uproc %s is already up to date\n", current)
		return nil
	}

	if checkOnly {
		fmt.Fprintf(c.OutOrStdout(), "A new version %s is available (current: %s). Run `uproc self-update` to update.\n", target, current)
		return nil
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not resolve current executable: %w", err)
	}

	switch update.DetectInstall(exePath) {
	case update.InstallHomebrew:
		fmt.Fprintln(c.OutOrStdout(), "uproc is installed via Homebrew. Update it with:")
		fmt.Fprintln(c.OutOrStdout(), "  brew update && brew upgrade uproc")
		return nil
	case update.InstallScoop:
		fmt.Fprintln(c.OutOrStdout(), "uproc is installed via Scoop. Update it with:")
		fmt.Fprintln(c.OutOrStdout(), "  scoop update uproc")
		return nil
	}

	asset, err := updater.Asset(rel, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}
	sums, err := updater.Checksums(rel)
	if err != nil {
		return err
	}
	if err := updater.Apply(rel, asset, sums[asset.Name], exePath); err != nil {
		return err
	}

	fmt.Fprintf(c.OutOrStdout(), "Updated uproc from %s to %s.\n", current, target)
	return nil
}
