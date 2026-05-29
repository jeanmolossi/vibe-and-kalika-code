package integration_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
)

func TestInstallGitPackage(t *testing.T) {
	repo := t.TempDir()
	copyDir(t, filepath.Join("..", "..", "testdata", "packages", "valid-basic"), repo)
	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "test@example.com")
	runGit(t, repo, "config", "user.name", "Test")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-m", "init")
	t.Setenv("VKC_STATE_DIR", t.TempDir())
	projectRoot := t.TempDir()
	copilotHome := filepath.Join(t.TempDir(), ".copilot")
	t.Setenv("COPILOT_HOME", copilotHome)
	_, code, err := app.Install(app.InstallOptions{
		Source:         "file://" + repo,
		ProjectRoot:    projectRoot,
		Targets:        []platform.Platform{platform.PlatformCopilotCLI},
		Yes:            true,
		ConflictAction: conflictActionBackupOverwrite,
	})
	if err != nil || code != app.ExitSuccess {
		t.Fatalf("Install() err=%v code=%d", err, code)
	}
	if _, err := os.Stat(filepath.Join(copilotHome, "agents", "test-agent.md")); err != nil {
		t.Fatalf("agent missing: %v", err)
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.CommandContext(context.Background(), "git", args...) //nolint:gosec // git is a known safe command
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}
