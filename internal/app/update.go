package app

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/source"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/version"
)

// UpdateOptions configures a package-update run.
type UpdateOptions struct {
	ProjectRoot string
	Yes         bool
	DryRun      bool
}

// UpdateResult reports which packages were updated, already up to date, or skipped.
type UpdateResult struct {
	Updated  []string // package names that were updated
	UpToDate []string // package names already at the latest version
	Skipped  []string // package names skipped (dry-run or user declined)
}

// SelfUpdateOptions configures a CLI self-update run.
type SelfUpdateOptions struct {
	DryRun bool
	Yes    bool
}

// SelfUpdateResult reports the outcome of a self-update attempt.
type SelfUpdateResult struct {
	CurrentVersion string
	LatestVersion  string
	AlreadyLatest  bool
	DryRun         bool
}

// githubRelease is a minimal representation of the GitHub releases API response.
type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// httpClientTimeout is the timeout applied to all HTTP calls for self-update.
const httpClientTimeout = 30 * time.Second

// httpClient is the shared HTTP client used for self-update downloads.
// Using an interface so tests can inject a mock.
type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

var defaultHTTPClient httpDoer = &http.Client{Timeout: httpClientTimeout}

// Update checks each installed package for a newer version and re-installs outdated ones.
func Update(opts UpdateOptions) (*UpdateResult, int, error) {
	st, err := state.Read(opts.ProjectRoot)
	if err != nil {
		return nil, ExitError, fmt.Errorf("read state: %w", err)
	}

	result := &UpdateResult{}

	for _, inst := range st.Installations {
		latest, err := fetchLatestVersion(inst.Source, opts.ProjectRoot)
		if err != nil {
			// Non-fatal: skip this package and report it.
			result.Skipped = append(result.Skipped, inst.Package)
			continue
		}

		if !isOutdated(inst.Version, latest) {
			result.UpToDate = append(result.UpToDate, inst.Package)
			continue
		}

		if opts.DryRun {
			result.Skipped = append(result.Skipped, inst.Package)
			continue
		}

		if !opts.Yes {
			confirmed, err := confirmUpdate(inst.Package, inst.Version, latest, os.Stdin)
			if err != nil || !confirmed {
				result.Skipped = append(result.Skipped, inst.Package)
				continue
			}
		}

		_, _, err = Install(InstallOptions{
			Source:         inst.Source,
			ProjectRoot:    opts.ProjectRoot,
			Yes:            true,
			ConflictAction: "overwrite",
		})
		if err != nil {
			result.Skipped = append(result.Skipped, inst.Package)
			continue
		}

		result.Updated = append(result.Updated, inst.Package)
	}

	return result, ExitSuccess, nil
}

// SelfUpdate fetches the latest GitHub release and replaces the running binary.
func SelfUpdate(opts SelfUpdateOptions) (*SelfUpdateResult, int, error) {
	return selfUpdateWithClient(opts, defaultHTTPClient)
}

// selfUpdateWithClient is the internal implementation that accepts an httpDoer for testing.
func selfUpdateWithClient(opts SelfUpdateOptions, client httpDoer) (*SelfUpdateResult, int, error) {
	current := version.Version

	rel, err := fetchLatestRelease(client)
	if err != nil {
		return nil, ExitError, fmt.Errorf("fetch latest release: %w", err)
	}

	latest := strings.TrimPrefix(rel.TagName, "v")
	currentNorm := strings.TrimPrefix(current, "v")

	if currentNorm == latest {
		return &SelfUpdateResult{
			CurrentVersion: current,
			LatestVersion:  rel.TagName,
			AlreadyLatest:  true,
		}, ExitSuccess, nil
	}

	if opts.DryRun {
		return &SelfUpdateResult{
			CurrentVersion: current,
			LatestVersion:  rel.TagName,
			DryRun:         true,
		}, ExitSuccess, nil
	}

	if !opts.Yes {
		confirmed, promptErr := confirmSelfUpdate(current, rel.TagName, os.Stdin)
		if promptErr != nil || !confirmed {
			return nil, ExitUserCancelled, fmt.Errorf("update canceled")
		}
	}

	if err := downloadAndReplace(client, rel, latest); err != nil {
		return nil, ExitError, err
	}

	return &SelfUpdateResult{
		CurrentVersion: current,
		LatestVersion:  rel.TagName,
	}, ExitSuccess, nil
}

// fetchLatestVersion resolves a package source and returns its manifest version.
func fetchLatestVersion(src, projectRoot string) (string, error) {
	resolved, err := source.Resolve(src, projectRoot)
	if err != nil {
		return "", err
	}
	defer resolved.Cleanup() //nolint:errcheck // temp directory cleanup error is non-actionable

	m, err := manifest.ParseFile(resolved.Root)
	if err != nil {
		return "", err
	}

	return m.Version, nil
}

// isOutdated returns true when the installed version differs from the latest.
func isOutdated(installed, latest string) bool {
	return strings.TrimPrefix(installed, "v") != strings.TrimPrefix(latest, "v")
}

// confirmUpdate prompts the user to confirm updating a single package.
func confirmUpdate(pkg, from, to string, r io.Reader) (bool, error) {
	fmt.Printf("Update %s from %s to %s? [y/N] ", pkg, from, to)
	return readYesNo(r)
}

// confirmSelfUpdate prompts the user to confirm a CLI self-update.
func confirmSelfUpdate(from, to string, r io.Reader) (bool, error) {
	fmt.Printf("Update vkc from %s to %s? [y/N] ", from, to)
	return readYesNo(r)
}

// readYesNo reads one line and returns true only for "y" or "Y".
func readYesNo(r io.Reader) (bool, error) {
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		return false, scanner.Err()
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y", nil
}

// fetchLatestRelease queries the GitHub releases API.
func fetchLatestRelease(client httpDoer) (*githubRelease, error) {
	const apiURL = "https://api.github.com/repos/jeanmolossi/vibe-and-kalika-code/releases/latest"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiURL, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "vkc/"+version.Version)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("decode release: %w", err)
	}

	return &rel, nil
}

// downloadAndReplace downloads the release asset, verifies its checksum, and
// atomically replaces the running binary.
func downloadAndReplace(client httpDoer, rel *githubRelease, latestVersion string) error {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Determine asset file name.
	var assetName string
	if goos == "windows" {
		assetName = fmt.Sprintf("vkc_%s_%s_%s.zip", latestVersion, goos, goarch)
	} else {
		assetName = fmt.Sprintf("vkc_%s_%s_%s.tar.gz", latestVersion, goos, goarch)
	}
	checksumName := fmt.Sprintf("vkc_%s_checksums.txt", latestVersion)

	assetURL := findAssetURL(rel.Assets, assetName)
	if assetURL == "" {
		return fmt.Errorf("no release asset found for %s/%s (looking for %s)", goos, goarch, assetName)
	}
	checksumURL := findAssetURL(rel.Assets, checksumName)
	if checksumURL == "" {
		return fmt.Errorf("checksums asset not found (%s)", checksumName)
	}

	// Download and parse checksums.
	expectedHash, err := fetchExpectedChecksum(client, checksumURL, assetName)
	if err != nil {
		return fmt.Errorf("fetch checksum: %w", err)
	}

	// Locate the running binary first so we can create the staging temp dir on
	// the same filesystem, avoiding EXDEV (cross-device link) errors from os.Rename.
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locate executable: %w", err)
	}
	execDir := filepath.Dir(execPath)

	// Download the asset into a temp dir on the same device as the binary.
	tmpDir, err := os.MkdirTemp(execDir, "vkc-update-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir) //nolint:errcheck // temp dir cleanup on exit, non-actionable

	assetPath := filepath.Join(tmpDir, assetName)
	if err := downloadFile(client, assetURL, assetPath); err != nil {
		return fmt.Errorf("download asset: %w", err)
	}

	// Verify checksum.
	if err := verifyChecksum(assetPath, expectedHash); err != nil {
		return fmt.Errorf("checksum mismatch: %w", err)
	}

	// Extract the binary.
	var binPath string
	if goos == "windows" {
		binPath, err = extractZip(assetPath, tmpDir)
	} else {
		binPath, err = extractTarGz(assetPath, tmpDir)
	}
	if err != nil {
		return fmt.Errorf("extract binary: %w", err)
	}

	if err := os.Chmod(binPath, 0o755); err != nil {
		return fmt.Errorf("chmod binary: %w", err)
	}

	// os.Rename is atomic on Unix when src and dst are on the same filesystem.
	if err := os.Rename(binPath, execPath); err != nil {
		return fmt.Errorf("replace binary: %w", err)
	}

	return nil
}

// findAssetURL returns the download URL for the named asset, or empty string.
func findAssetURL(assets []githubAsset, name string) string {
	for _, a := range assets {
		if a.Name == name {
			return a.BrowserDownloadURL
		}
	}
	return ""
}

// fetchExpectedChecksum downloads the checksums file and returns the SHA256
// hex string for the given file name.
func fetchExpectedChecksum(client httpDoer, url, filename string) (string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "vkc/"+version.Version)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d fetching checksums", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		// Format: "<hash>  <filename>"
		parts := strings.Fields(line)
		if len(parts) == 2 && parts[1] == filename {
			return parts[0], nil
		}
	}

	return "", fmt.Errorf("checksum for %s not found in checksums file", filename)
}

// downloadFile fetches url and writes the body to dest.
func downloadFile(client httpDoer, url, dest string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "vkc/"+version.Version)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d downloading asset", resp.StatusCode)
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	// Limit download to 100 MiB to guard against runaway responses.
	const maxAssetSize = 100 << 20
	_, err = io.Copy(f, io.LimitReader(resp.Body, maxAssetSize))
	return err
}

// verifyChecksum computes the SHA256 of path and compares it to expected.
func verifyChecksum(path, expected string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	got := hex.EncodeToString(h.Sum(nil))
	if got != expected {
		return fmt.Errorf("got %s, want %s", got, expected)
	}
	return nil
}

// extractTarGz extracts a .tar.gz archive and returns the path to the "vkc" binary.
func extractTarGz(src, destDir string) (string, error) {
	f, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gr.Close() //nolint:errcheck // deferred gzip reader close, non-actionable

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		base := filepath.Base(hdr.Name)
		if base != "vkc" {
			continue
		}

		destPath := filepath.Join(destDir, base)
		out, err := os.Create(destPath)
		if err != nil {
			return "", err
		}

		if _, err := io.Copy(out, tr); err != nil { //nolint:gosec // controlled extraction
			out.Close()
			return "", err
		}
		out.Close()

		return destPath, nil
	}

	return "", fmt.Errorf("vkc binary not found in archive")
}

// extractZip extracts a .zip archive and returns the path to the "vkc.exe" binary.
func extractZip(src, destDir string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close() //nolint:errcheck // deferred zip reader close, non-actionable

	for _, f := range r.File {
		base := filepath.Base(f.Name)
		if base != "vkc.exe" && base != "vkc" {
			continue
		}

		destPath := filepath.Join(destDir, base)
		rc, err := f.Open()
		if err != nil {
			return "", err
		}

		out, err := os.Create(destPath)
		if err != nil {
			rc.Close()
			return "", err
		}

		_, err = io.Copy(out, rc) //nolint:gosec // controlled extraction
		out.Close()
		rc.Close()

		if err != nil {
			return "", err
		}

		return destPath, nil
	}

	return "", fmt.Errorf("vkc binary not found in zip archive")
}
