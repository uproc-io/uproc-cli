package update

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// InstallMethod describes how the CLI was installed.
type InstallMethod int

const (
	// InstallUnknown cannot determine the installation method.
	InstallUnknown InstallMethod = iota
	// InstallHomebrew the binary lives under a Homebrew Cellar.
	InstallHomebrew
	// InstallScoop the binary lives under the Scoop shim directory.
	InstallScoop
	// InstallStandalone the binary was installed manually.
	InstallStandalone
)

func (m InstallMethod) String() string {
	switch m {
	case InstallHomebrew:
		return "homebrew"
	case InstallScoop:
		return "scoop"
	case InstallStandalone:
		return "standalone"
	default:
		return "unknown"
	}
}

// DetectInstall inspects the executable path to infer the install method.
func DetectInstall(exePath string) InstallMethod {
	p := filepath.ToSlash(exePath)
	lower := strings.ToLower(p)
	switch {
	case strings.Contains(lower, "/homebrew/"), strings.Contains(lower, "/cellar/"):
		return InstallHomebrew
	case strings.Contains(lower, "/scoop/"):
		return InstallScoop
	default:
		return InstallStandalone
	}
}

// replaceExecutable atomically swaps the running binary for the downloaded one.
// On Windows an in-use executable cannot be replaced; the new binary is left
// next to it as "<name>.new" with instructions.
func replaceExecutable(newBin, exePath string) error {
	if exePath == "" {
		return errors.New("cannot resolve current executable path")
	}

	if runtime.GOOS == "windows" {
		newPath := exePath + ".new"
		if err := os.Rename(newBin, newPath); err != nil {
			return fmt.Errorf("could not stage new binary: %w", err)
		}
		fmt.Fprintf(os.Stderr, "New binary staged at %s.\n", newPath)
		fmt.Fprintln(os.Stderr, "Close any running uproc instance and replace the executable manually (or re-run `uproc self-update`).")
		return nil
	}

	if err := os.Chmod(newBin, 0o755); err != nil {
		return err
	}
	dir := filepath.Dir(exePath)
	staged := filepath.Join(dir, ".uproc-new")
	if err := os.Rename(newBin, staged); err != nil {
		return fmt.Errorf("could not stage new binary: %w", err)
	}
	if err := os.Rename(staged, exePath); err != nil {
		_ = os.Remove(staged)
		return fmt.Errorf("could not replace executable: %w", err)
	}
	return nil
}
