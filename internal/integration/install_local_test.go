package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
)

func TestInstallLocalPackage(t *testing.T) {
	projectRoot := t.TempDir()
	copilotHome := filepath.Join(t.TempDir(), ".copilot")
	t.Setenv("COPILOT_HOME", copilotHome)
	pkg := filepath.Join("..", "..", "testdata", "packages", "valid-basic")
	res, code, err := app.Install(app.InstallOptions{Source: pkg, ProjectRoot: projectRoot, Targets: []platform.Platform{platform.PlatformCopilotCLI}, Yes: true, ConflictAction: "backup-and-overwrite"})
	if err != nil || code != app.ExitSuccess {
		t.Fatalf("Install() err=%v code=%d", err, code)
	}
	if _, err := os.Stat(filepath.Join(copilotHome, "agents", "test-agent.md")); err != nil {
		t.Fatalf("agent missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(copilotHome, "skills", "test-skill", "SKILL.md")); err != nil {
		t.Fatalf("skill missing: %v", err)
	}
	if _, err := os.Stat(res.ReportPath); err != nil {
		t.Fatalf("report missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projectRoot, ".ai-setup", "installed.yaml")); err != nil {
		t.Fatalf("state missing: %v", err)
	}
}
