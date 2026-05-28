package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
)

func TestConflictSkipAndBackupOverwrite(t *testing.T) {
	projectRoot := t.TempDir()
	copilotHome := filepath.Join(t.TempDir(), ".copilot")
	t.Setenv("COPILOT_HOME", copilotHome)
	if err := os.MkdirAll(filepath.Join(copilotHome, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(copilotHome, "agents", "test-agent.md")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	pkg := filepath.Join("..", "..", "testdata", "packages", "with-conflicts")
	_, code, err := app.Install(app.InstallOptions{Source: pkg, ProjectRoot: projectRoot, Targets: []platform.Platform{platform.PlatformCopilotCLI}, Yes: true, ConflictAction: "skip"})
	if err != nil || code != app.ExitSuccess {
		t.Fatalf("skip install err=%v code=%d", err, code)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "old" {
		t.Fatalf("expected skip to preserve file, got %q", string(data))
	}
	res, code, err := app.Install(app.InstallOptions{Source: pkg, ProjectRoot: projectRoot, Targets: []platform.Platform{platform.PlatformCopilotCLI}, Yes: true, ConflictAction: conflictActionBackupOverwrite})
	if err != nil || code != app.ExitSuccess {
		t.Fatalf("overwrite install err=%v code=%d", err, code)
	}
	if res.BackupPath == "" {
		t.Fatal("expected backup path")
	}
	newData, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(newData) == "old" {
		t.Fatal("expected overwrite to replace file")
	}
}
