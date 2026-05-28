package source_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/source"
)

// TestIsGitURL verifies URL classification for all supported schemes.
func TestIsGitURL(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"https://github.com/org/repo.git", true},
		{"http://github.com/org/repo.git", true},
		{"git@github.com:org/repo.git", true},
		{"ssh://git@github.com/org/repo.git", true},
		{"file:///tmp/repo", true},
		{"./local/path", false},
		{"/absolute/path", false},
		{"local-dir", false},
	}
	for _, tc := range cases {
		got := source.IsGitURL(tc.input)
		if got != tc.want {
			t.Errorf("IsGitURL(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

// TestResolveLocalSource resolves a local directory and returns the absolute path.
func TestResolveLocalSource(t *testing.T) {
	dir := t.TempDir()
	res, err := source.ResolveLocalSource(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer res.Cleanup() //nolint:errcheck // test cleanup
	abs, _ := filepath.Abs(dir)
	if res.Root != abs {
		t.Errorf("Root = %q, want %q", res.Root, abs)
	}
}

// TestResolveLocalSourceMissing returns an error when the path does not exist.
func TestResolveLocalSourceMissing(t *testing.T) {
	_, err := source.ResolveLocalSource("/no/such/path/here")
	if err == nil {
		t.Fatal("expected error for missing path, got nil")
	}
}

// TestResolveLocalSourceCleanupIsNoop verifies the cleanup func does nothing for a local path.
func TestResolveLocalSourceCleanupIsNoop(t *testing.T) {
	dir := t.TempDir()
	res, err := source.ResolveLocalSource(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := res.Cleanup(); err != nil {
		t.Errorf("cleanup should be a no-op, got: %v", err)
	}
	// Directory must still exist after noop cleanup.
	if _, serr := os.Stat(dir); serr != nil {
		t.Errorf("directory was removed by cleanup: %v", serr)
	}
}

// TestResolveDispatchesLocalForNonGitURL confirms Resolve uses the local provider.
func TestResolveDispatchesLocalForNonGitURL(t *testing.T) {
	dir := t.TempDir()
	projectRoot := t.TempDir()
	res, err := source.Resolve(dir, projectRoot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer res.Cleanup() //nolint:errcheck // test cleanup
	abs, _ := filepath.Abs(dir)
	if res.Root != abs {
		t.Errorf("Root = %q, want %q", res.Root, abs)
	}
}
