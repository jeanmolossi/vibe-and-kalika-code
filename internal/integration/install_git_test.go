package integration_test

import (
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
	runCmd(t, repo, "git", "init")
	runCmd(t, repo, "git", "config", "user.email", "test@example.com")
	runCmd(t, repo, "git", "config", "user.name", "Test")
	runCmd(t, repo, "git", "add", ".")
	runCmd(t, repo, "git", "commit", "-m", "init")
	projectRoot := t.TempDir()
	copilotHome := filepath.Join(t.TempDir(), ".copilot")
	t.Setenv("COPILOT_HOME", copilotHome)
	_, code, err := app.Install(app.InstallOptions{Source: "file://" + repo, ProjectRoot: projectRoot, Targets: []platform.Platform{platform.PlatformCopilotCLI}, Yes: true, ConflictAction: "backup-and-overwrite"})
	if err != nil || code != app.ExitSuccess {
		t.Fatalf("Install() err=%v code=%d", err, code)
	}
	if _, err := os.Stat(filepath.Join(copilotHome, "agents", "test-agent.md")); err != nil {
		t.Fatalf("agent missing: %v", err)
	}
}

func runCmd(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, out)
	}
}
