package security_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/security"
)

func TestRejectEscapingSymlinks(t *testing.T) {
	dir := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.txt")
	if err := os.WriteFile(outside, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "link.txt")); err != nil {
		t.Fatal(err)
	}
	if err := security.RejectEscapingSymlinks(dir); err == nil {
		t.Fatal("expected symlink escape to be rejected")
	}
}

// TestEnsureResolvedWithinRootRejectsSymlinkEscape verifies that a pre-existing symlink
// inside the allowed root that points outside is detected and rejected before any write.
func TestEnsureResolvedWithinRootRejectsSymlinkEscape(t *testing.T) {
	allowedRoot := t.TempDir()
	outside := t.TempDir()

	// Create a symlink inside allowedRoot that points to the outside dir
	if err := os.Symlink(outside, filepath.Join(allowedRoot, "subdir")); err != nil {
		t.Fatal(err)
	}

	// A write to allowedRoot/subdir/secret.txt would resolve to outside/secret.txt
	target := filepath.Join(allowedRoot, "subdir", "secret.txt")
	if err := security.EnsureResolvedWithinRoot(allowedRoot, target); err == nil {
		t.Fatal("expected symlink escape to be rejected by EnsureResolvedWithinRoot")
	}
}

// TestEnsureResolvedWithinRootAllowsLegitimate verifies that normal (non-symlink) paths pass.
func TestEnsureResolvedWithinRootAllowsLegitimate(t *testing.T) {
	allowedRoot := t.TempDir()
	target := filepath.Join(allowedRoot, "agents", "foo.md")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := security.EnsureResolvedWithinRoot(allowedRoot, target); err != nil {
		t.Fatalf("expected legitimate path to pass: %v", err)
	}
}
