package integration_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
)

// TestInstallCopilotAndClaude verifies that a package can be installed into
// both Copilot CLI and Claude Code simultaneously.
func TestInstallCopilotAndClaude(t *testing.T) {
	t.Setenv("VKC_STATE_DIR", t.TempDir())
	projectRoot := t.TempDir()
	copilotHome := filepath.Join(t.TempDir(), ".copilot")
	t.Setenv("COPILOT_HOME", copilotHome)

	pkg := filepath.Join("..", "..", "testdata", "packages", "valid-multi")
	targets := []platform.Platform{platform.PlatformCopilotCLI, platform.PlatformClaudeCode}

	_, code, err := app.Install(app.InstallOptions{
		Source:         pkg,
		ProjectRoot:    projectRoot,
		Targets:        targets,
		Yes:            true,
		ConflictAction: conflictActionBackupOverwrite,
	})
	if err != nil || code != app.ExitSuccess {
		t.Fatalf("Install() err=%v code=%d", err, code)
	}

	// Copilot CLI assertions
	agentCopilot := filepath.Join(copilotHome, "agents", "kalika-reviewer.md")
	if _, err := os.Stat(agentCopilot); err != nil {
		t.Errorf("copilot agent missing at %s: %v", agentCopilot, err)
	}
	skillCopilot := filepath.Join(copilotHome, "skills", "code-review")
	if _, err := os.Stat(skillCopilot); err != nil {
		t.Errorf("copilot skill dir missing at %s: %v", skillCopilot, err)
	}

	// Claude Code assertions
	agentClaude := filepath.Join(projectRoot, ".claude", "agents", "kalika-reviewer.md")
	if _, err := os.Stat(agentClaude); err != nil {
		t.Errorf("claude agent missing at %s: %v", agentClaude, err)
	}
	skillClaude := filepath.Join(projectRoot, ".claude", "skills", "code-review")
	if _, err := os.Stat(skillClaude); err != nil {
		t.Errorf("claude skill dir missing at %s: %v", skillClaude, err)
	}
}

// TestInstallCodexAgentIntoExistingAgentsMD verifies that installing a Codex
// agent appends a managed block to an existing AGENTS.md without destroying
// pre-existing content.
func TestInstallCodexAgentIntoExistingAgentsMD(t *testing.T) {
	t.Setenv("VKC_STATE_DIR", t.TempDir())
	projectRoot := t.TempDir()

	// Create pre-existing AGENTS.md
	agentsMD := filepath.Join(projectRoot, "AGENTS.md")
	existingContent := "# Project Agents\n\nThis file documents the project agents.\n"
	if err := os.WriteFile(agentsMD, []byte(existingContent), 0o644); err != nil {
		t.Fatalf("failed to create AGENTS.md: %v", err)
	}

	pkg := filepath.Join("..", "..", "testdata", "packages", "valid-multi")
	_, code, err := app.Install(app.InstallOptions{
		Source:         pkg,
		ProjectRoot:    projectRoot,
		Targets:        []platform.Platform{platform.PlatformCodexCLI},
		Yes:            true,
		ConflictAction: conflictActionBackupOverwrite,
	})
	if err != nil || code != app.ExitSuccess {
		t.Fatalf("Install() err=%v code=%d", err, code)
	}

	data, err := os.ReadFile(agentsMD)
	if err != nil {
		t.Fatalf("failed to read AGENTS.md: %v", err)
	}
	content := string(data)

	// Pre-existing content must be preserved
	if !strings.Contains(content, existingContent) {
		t.Errorf("pre-existing AGENTS.md content was deleted; got:\n%s", content)
	}

	// Managed block must be present
	marker := "<!-- BEGIN VKC AGENT: kalika-reviewer -->"
	if !strings.Contains(content, marker) {
		t.Errorf("managed block marker %q not found in AGENTS.md; got:\n%s", marker, content)
	}
}
