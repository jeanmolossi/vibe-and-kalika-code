package planner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/planner"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
)

func TestExistingFileBecomesConflict(t *testing.T) {
	projectRoot := t.TempDir()
	copilotHome := filepath.Join(projectRoot, "copilot-home")
	t.Setenv("COPILOT_HOME", copilotHome)
	if err := os.MkdirAll(filepath.Join(copilotHome, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(copilotHome, "agents", "test-agent.md"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	pkgRoot := filepath.Join("..", "..", "testdata", "packages", "valid-basic")
	m, err := manifest.ParseFile(pkgRoot)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := planner.Build(pkgRoot, projectRoot, m, []platform.Platform{platform.PlatformCopilotCLI}, "skip")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Operations) == 0 || plan.Operations[0].Conflict == nil {
		t.Fatalf("expected conflict, got %+v", plan.Operations)
	}
}
