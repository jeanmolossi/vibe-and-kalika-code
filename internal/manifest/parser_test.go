package manifest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"
)

func TestParseFileValidManifest(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "packages", "valid-basic")
	m, err := manifest.ParseFile(root)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}
	if m.Name != "valid-basic" {
		t.Fatalf("unexpected name %q", m.Name)
	}
}

func TestParseFileMissingFields(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte("name: x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := manifest.ParseFile(dir)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}
	if m.Name != "x" || m.Version != "" {
		t.Fatalf("unexpected manifest: %+v", m)
	}
}
