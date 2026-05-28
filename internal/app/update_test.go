package app

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/version"
)

const version100 = "1.0.0"

const devVersion = "dev"

// --- helpers -----------------------------------------------------------------

type mockHTTPClient struct {
	responses []*http.Response
	idx       int
}

func (m *mockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	if m.idx >= len(m.responses) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	}
	resp := m.responses[m.idx]
	m.idx++
	return resp, nil
}

func jsonBody(v any) io.ReadCloser {
	b, _ := json.Marshal(v)
	return io.NopCloser(strings.NewReader(string(b)))
}

func stringBody(s string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(s))
}

// --- isOutdated --------------------------------------------------------------

func TestIsOutdated(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		installed string
		latest    string
		want      bool
	}{
		{"dev vs release", devVersion, version100, true},
		{"same version with v prefix", "v1.0.0", version100, false},
		{"same version no prefix", version100, version100, false},
		{"older installed", version100, "1.1.0", true},
		{"both dev", devVersion, devVersion, false},
		{"v prefix both sides", "v1.2.3", "v1.2.3", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := isOutdated(tc.installed, tc.latest)
			if got != tc.want {
				t.Errorf("isOutdated(%q, %q) = %v; want %v", tc.installed, tc.latest, got, tc.want)
			}
		})
	}
}

// --- readYesNo ---------------------------------------------------------------

func TestReadYesNo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  bool
	}{
		{"y\n", true},
		{"Y\n", true},
		{"n\n", false},
		{"N\n", false},
		{"yes\n", false}, // only "y" accepted
		{"\n", false},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			got, err := readYesNo(strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("readYesNo(%q) = %v; want %v", tc.input, got, tc.want)
			}
		})
	}
}

// --- Update (empty state) ----------------------------------------------------

func TestUpdate_EmptyState(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	res, code, err := Update(UpdateOptions{
		ProjectRoot: dir,
		Yes:         true,
		DryRun:      false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d; want %d", code, ExitSuccess)
	}
	if len(res.Updated) != 0 || len(res.Skipped) != 0 || len(res.UpToDate) != 0 {
		t.Errorf("expected empty result, got %+v", res)
	}
}

// --- SelfUpdate: already latest ----------------------------------------------

func TestSelfUpdate_AlreadyLatest(t *testing.T) {
	t.Parallel()

	// Build a tag that will match the current version.Version default ("dev").
	// Using the package-level default avoids hardcoding "vdev" which would
	// break if tests are built with a real version via ldflags.
	currentTag := "v" + version.Version
	rel := githubRelease{TagName: currentTag}
	client := &mockHTTPClient{
		responses: []*http.Response{
			{StatusCode: http.StatusOK, Body: jsonBody(rel)},
		},
	}

	res, code, err := selfUpdateWithClient(SelfUpdateOptions{}, client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d; want %d", code, ExitSuccess)
	}
	if !res.AlreadyLatest {
		t.Errorf("AlreadyLatest = false; want true")
	}
}

// --- SelfUpdate: dry-run shows update plan -----------------------------------

func TestSelfUpdate_DryRun(t *testing.T) {
	t.Parallel()

	rel := githubRelease{TagName: "v2.0.0"}
	client := &mockHTTPClient{
		responses: []*http.Response{
			{StatusCode: http.StatusOK, Body: jsonBody(rel)},
		},
	}

	res, code, err := selfUpdateWithClient(SelfUpdateOptions{DryRun: true}, client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != ExitSuccess {
		t.Errorf("exit code = %d; want %d", code, ExitSuccess)
	}
	if res.AlreadyLatest {
		t.Errorf("AlreadyLatest should be false for dry-run with newer version")
	}
	if !res.DryRun {
		t.Errorf("DryRun should be true")
	}
	if res.LatestVersion != "v2.0.0" {
		t.Errorf("LatestVersion = %q; want %q", res.LatestVersion, "v2.0.0")
	}
}

// --- fetchExpectedChecksum ---------------------------------------------------

func TestFetchExpectedChecksum(t *testing.T) {
	t.Parallel()

	checksumContent := `abc123  vkc_1.0.0_linux_amd64.tar.gz
def456  vkc_1.0.0_darwin_arm64.tar.gz
`
	client := &mockHTTPClient{
		responses: []*http.Response{
			{StatusCode: http.StatusOK, Body: stringBody(checksumContent)},
		},
	}

	hash, err := fetchExpectedChecksum(client, "http://example.com/checksums.txt", "vkc_1.0.0_linux_amd64.tar.gz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash != "abc123" {
		t.Errorf("hash = %q; want %q", hash, "abc123")
	}
}

func TestFetchExpectedChecksum_NotFound(t *testing.T) {
	t.Parallel()

	client := &mockHTTPClient{
		responses: []*http.Response{
			{StatusCode: http.StatusOK, Body: stringBody("abc123  other_file.tar.gz\n")},
		},
	}

	_, err := fetchExpectedChecksum(client, "http://example.com/checksums.txt", "missing.tar.gz")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestFetchExpectedChecksum_HTTPError(t *testing.T) {
	t.Parallel()

	client := &mockHTTPClient{
		responses: []*http.Response{
			{StatusCode: http.StatusForbidden, Body: stringBody("rate limited\n")},
		},
	}

	_, err := fetchExpectedChecksum(client, "http://example.com/checksums.txt", "any.tar.gz")
	if err == nil {
		t.Fatal("expected error for HTTP 403, got nil")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("error should mention HTTP status, got: %v", err)
	}
}

func TestDownloadFile_HTTPError(t *testing.T) {
	t.Parallel()

	client := &mockHTTPClient{
		responses: []*http.Response{
			{StatusCode: http.StatusNotFound, Body: stringBody("not found\n")},
		},
	}

	dest := t.TempDir() + "/asset.tar.gz"
	err := downloadFile(client, "http://example.com/asset.tar.gz", dest)
	if err == nil {
		t.Fatal("expected error for HTTP 404, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("error should mention HTTP status, got: %v", err)
	}
}

// --- findAssetURL ------------------------------------------------------------

func TestFindAssetURL(t *testing.T) {
	t.Parallel()

	assets := []githubAsset{
		{Name: "vkc_1.0.0_linux_amd64.tar.gz", BrowserDownloadURL: "http://example.com/linux"},
		{Name: "vkc_1.0.0_checksums.txt", BrowserDownloadURL: "http://example.com/checksums"},
	}

	if url := findAssetURL(assets, "vkc_1.0.0_linux_amd64.tar.gz"); url != "http://example.com/linux" {
		t.Errorf("findAssetURL = %q; want linux url", url)
	}
	if url := findAssetURL(assets, "vkc_1.0.0_checksums.txt"); url != "http://example.com/checksums" {
		t.Errorf("findAssetURL = %q; want checksums url", url)
	}
	if url := findAssetURL(assets, "missing"); url != "" {
		t.Errorf("findAssetURL = %q; want empty string", url)
	}
}
