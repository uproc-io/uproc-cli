package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/mod/semver"
)

const (
	checkInterval     = 24 * time.Hour
	updateCheckEnv    = "UPROC_NO_UPDATE_CHECK"
	notifyHTTPTimeout = 5 * time.Second
)

type checkCache struct {
	CheckedAt time.Time `json:"checked_at"`
	Latest    string    `json:"latest"`
}

func cachePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "uproc", "update-check.json"), nil
}

func readCache() (checkCache, error) {
	path, err := cachePath()
	if err != nil {
		return checkCache{}, err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return checkCache{}, err
	}
	var c checkCache
	if err := json.Unmarshal(b, &c); err != nil {
		return checkCache{}, err
	}
	return c, nil
}

func writeCache(c checkCache) error {
	path, err := cachePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o600)
}

func updateCheckDisabled() bool {
	switch os.Getenv(updateCheckEnv) {
	case "1", "true", "yes", "TRUE", "YES":
		return true
	}
	return false
}

// NotifyIfUpdateAvailable prints a non-blocking notice to stderr when a newer
// release exists. It uses a cached result for up to 24h and fails silently on
// any error (offline, API limits, etc.).
func NotifyIfUpdateAvailable(currentVersion string) {
	if updateCheckDisabled() {
		return
	}
	current := NormalizeVersion(currentVersion)
	if current == "" {
		return
	}

	cached, err := readCache()
	latest := ""
	fresh := err == nil && time.Since(cached.CheckedAt) < checkInterval
	if fresh {
		latest = cached.Latest
	} else {
		u := &Updater{HTTPClient: &http.Client{Timeout: notifyHTTPTimeout}}
		rel, err := u.Latest()
		if err == nil {
			latest = rel.Version()
			_ = writeCache(checkCache{CheckedAt: time.Now(), Latest: latest})
		}
	}

	if latest == "" || semver.Compare(current, latest) >= 0 {
		return
	}
	fmt.Fprintf(os.Stderr, "A new version %s is available. Run `uproc self-update` to update.\n", latest)
}
