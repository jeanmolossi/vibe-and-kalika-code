package manifest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"
)

func TestValidateValidManifest(t *testing.T) {
	projectRoot := t.TempDir()
	copilotHome := filepath.Join(projectRoot, ".copilot-home")
	t.Setenv("COPILOT_HOME", copilotHome)
	root := filepath.Join("..", "..", "testdata", "packages", "valid-basic")
	m, err := manifest.ParseFile(root)
	if err != nil {
		t.Fatal(err)
	}
	issues, err := manifest.Validate(root, projectRoot, m)
	if err != nil || len(issues) != 0 {
		t.Fatalf("Validate() issues=%v err=%v", issues, err)
	}
}

func TestValidatePathTraversalFails(t *testing.T) {
	projectRoot := t.TempDir()
	t.Setenv("COPILOT_HOME", filepath.Join(projectRoot, ".copilot-home"))
	root := filepath.Join("..", "..", "testdata", "packages", "invalid-path-traversal")
	m, err := manifest.ParseFile(root)
	if err != nil {
		t.Fatal(err)
	}
	issues, err := manifest.Validate(root, projectRoot, m)
	if err == nil || !strings.Contains(strings.Join(issues, " "), "path traversal") {
		t.Fatalf("expected path traversal error, got issues=%v err=%v", issues, err)
	}
}

func TestValidateUnknownTargetFails(t *testing.T) {
	dir := t.TempDir()
	manifestYAML := `name: x
version: v
description: d
targets: [alien-cli]
`
	if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte(manifestYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := manifest.ParseFile(dir)
	if err != nil {
		t.Fatal(err)
	}
	issues, err := manifest.Validate(dir, t.TempDir(), m)
	if err == nil || !strings.Contains(strings.Join(issues, " "), "unsupported target") {
		t.Fatalf("expected unsupported target, got issues=%v err=%v", issues, err)
	}
}

func TestValidateMissingSkillMarkdownFails(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "skills", "bad-skill"), 0o755); err != nil {
		t.Fatal(err)
	}
	manifestYAML := `name: x
version: v
description: d
targets: [copilot-cli]
skills:
  - name: bad-skill
    source: skills/bad-skill
    targets:
      copilot-cli:
        scope: user
`
	if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte(manifestYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := manifest.ParseFile(dir)
	if err != nil {
		t.Fatal(err)
	}
	issues, err := manifest.Validate(dir, t.TempDir(), m)
	if err == nil || !strings.Contains(strings.Join(issues, " "), "SKILL.md was not found") {
		t.Fatalf("expected missing SKILL.md, got issues=%v err=%v", issues, err)
	}
}
