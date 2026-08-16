package update

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/mod/semver"
)

const (
	repoOwner = "uproc-io"
	repoName  = "uproc.cli"
	apiBase   = "https://api.github.com"
)

// ErrUnknownVersion is returned when the running binary has no released
// version (local builds use "dev").
var ErrUnknownVersion = errors.New("unknown version (dev build)")

type Asset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
	Size int64  `json:"size"`
}

type Release struct {
	TagName    string  `json:"tag_name"`
	Prerelease bool    `json:"prerelease"`
	Assets     []Asset `json:"assets"`
}

// Version returns the release tag normalized for semver comparison.
func (r Release) Version() string {
	return NormalizeVersion(r.TagName)
}

// NormalizeVersion prefixes a leading "v" when missing and validates it is
// semver-compatible. An empty string is returned for invalid input.
func NormalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	if v == "" || strings.EqualFold(v, "dev") {
		return ""
	}
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	if !semver.IsValid(v) {
		return ""
	}
	return v
}

// Updater fetches releases from the uproc.cli GitHub repository.
type Updater struct {
	HTTPClient *http.Client
}

func (u *Updater) client() *http.Client {
	if u.HTTPClient != nil {
		return u.HTTPClient
	}
	return http.DefaultClient
}

// Latest returns the most recent stable release (no prereleases).
func (u *Updater) Latest() (Release, error) {
	var rel Release
	if err := u.getJSON(fmt.Sprintf("%s/repos/%s/%s/releases/latest", apiBase, repoOwner, repoName), &rel); err != nil {
		return Release{}, err
	}
	return rel, nil
}

// LatestPrerelease returns the most recent release allowing prereleases.
func (u *Updater) LatestPrerelease() (Release, error) {
	var rels []Release
	if err := u.getJSON(fmt.Sprintf("%s/repos/%s/%s/releases?per_page=20", apiBase, repoOwner, repoName), &rels); err != nil {
		return Release{}, err
	}
	for _, rel := range rels {
		if rel.Prerelease {
			return rel, nil
		}
	}
	return Release{}, errors.New("no pre-release found")
}

// ReleaseForVersion returns the release matching a specific version tag.
func (u *Updater) ReleaseForVersion(version string) (Release, error) {
	tag := strings.TrimSpace(version)
	if tag == "" {
		return Release{}, errors.New("version is required")
	}
	if !strings.HasPrefix(tag, "v") {
		tag = "v" + tag
	}
	var rel Release
	if err := u.getJSON(fmt.Sprintf("%s/repos/%s/%s/releases/tags/%s", apiBase, repoOwner, repoName, tag), &rel); err != nil {
		return Release{}, err
	}
	return rel, nil
}

// Asset returns the release asset for the current platform.
func (u *Updater) Asset(rel Release, goos, goarch string) (Asset, error) {
	assetVersion := strings.TrimPrefix(rel.TagName, "v")
	suffix := ".tar.gz"
	if goos == "windows" {
		suffix = ".zip"
	}
	want := fmt.Sprintf("uproc.cli_%s_%s_%s%s", assetVersion, goos, goarch, suffix)
	for _, a := range rel.Assets {
		if a.Name == want {
			return a, nil
		}
	}
	return Asset{}, fmt.Errorf("no asset for %s/%s (expected %s) in release %s", goos, goarch, want, rel.TagName)
}

// Checksums fetches and parses the release checksums.txt into
// filename -> sha256 hex.
func (u *Updater) Checksums(rel Release) (map[string]string, error) {
	for _, a := range rel.Assets {
		if a.Name != "checksums.txt" {
			continue
		}
		res, err := u.client().Get(a.URL)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode >= 400 {
			return nil, fmt.Errorf("checksums download failed: http %d", res.StatusCode)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return parseChecksums(body), nil
	}
	return nil, errors.New("checksums.txt not found in release")
}

func parseChecksums(body []byte) map[string]string {
	sums := map[string]string{}
	for _, line := range strings.Split(string(body), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		sums[fields[len(fields)-1]] = fields[0]
	}
	return sums
}

// Apply downloads, verifies and swaps the running binary for the given asset.
// exePath is the path of the current executable.
func (u *Updater) Apply(rel Release, asset Asset, expectedSHA, exePath string) error {
	tmpDir, err := os.MkdirTemp("", "uproc-update-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, asset.Name)
	if err := u.download(asset.URL, archivePath); err != nil {
		return err
	}

	sum, err := fileSHA256(archivePath)
	if err != nil {
		return err
	}
	if expectedSHA != "" && !strings.EqualFold(sum, strings.TrimSpace(expectedSHA)) {
		return fmt.Errorf("checksum mismatch for %s: got %s, expected %s", asset.Name, sum, expectedSHA)
	}

	binPath, err := extractBinary(archivePath, asset.Name)
	if err != nil {
		return err
	}

	return replaceExecutable(binPath, exePath)
}

func (u *Updater) getJSON(url string, out any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "uproc-cli/self-update")

	res, err := u.client().Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return fmt.Errorf("github api returned http %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(out)
}

func (u *Updater) download(url, dest string) error {
	res, err := u.client().Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return fmt.Errorf("download failed: http %d", res.StatusCode)
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, res.Body)
	return err
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// extractBinary extracts the "uproc" binary from a tar.gz or zip archive.
func extractBinary(archivePath, assetName string) (string, error) {
	binName := "uproc"
	if runtime.GOOS == "windows" {
		binName = "uproc.exe"
	}
	if strings.HasSuffix(assetName, ".zip") {
		return extractBinaryZip(archivePath, binName)
	}
	return extractBinaryTarGz(archivePath, binName)
}

func extractBinaryTarGz(archivePath, binName string) (string, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if hdr.Typeflag == tar.TypeReg && filepath.Base(hdr.Name) == binName {
			tmp, err := os.CreateTemp("", "uproc-bin-*")
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(tmp, tr); err != nil {
				tmp.Close()
				return "", err
			}
			if err := tmp.Close(); err != nil {
				return "", err
			}
			if err := os.Chmod(tmp.Name(), 0o755); err != nil {
				return "", err
			}
			return tmp.Name(), nil
		}
	}
	return "", fmt.Errorf("binary %s not found in archive", binName)
}

func extractBinaryZip(archivePath, binName string) (string, error) {
	zr, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer zr.Close()
	for _, f := range zr.File {
		if filepath.Base(f.Name) != binName {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()
		tmp, err := os.CreateTemp("", "uproc-bin-*")
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(tmp, rc); err != nil {
			tmp.Close()
			return "", err
		}
		if err := tmp.Close(); err != nil {
			return "", err
		}
		if err := os.Chmod(tmp.Name(), 0o755); err != nil {
			return "", err
		}
		return tmp.Name(), nil
	}
	return "", fmt.Errorf("binary %s not found in archive", binName)
}
