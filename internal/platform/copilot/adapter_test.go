package copilot

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
)

func TestTargetPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("COPILOT_HOME", home)
	if got := AgentTargetPath("agent"); got != filepath.Join(home, "agents", "agent.md") {
		t.Fatalf("agent path = %s", got)
	}
	if got := SkillTargetPath("skill"); got != filepath.Join(home, "skills", "skill") {
		t.Fatalf("skill path = %s", got)
	}
}

// TestPlanUsesCustomCopilotHome verifies that when COPILOT_HOME is set, both
// the TargetPath and AllowedRoot of planned operations use the custom path.
func TestPlanUsesCustomCopilotHome(t *testing.T) {
	customHome := t.TempDir()
	t.Setenv("COPILOT_HOME", customHome)

	packageRoot := t.TempDir()
	src := filepath.Join(packageRoot, "agent.md")
	if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &manifest.Manifest{
		Agents: []manifest.Agent{
			{
				Name:    "test-agent",
				Source:  "agent.md",
				Targets: map[string]manifest.AgentTarget{string(platform.PlatformCopilotCLI): {}},
			},
		},
	}

	adapter := NewAdapter()
	ops, err := adapter.Plan(platform.PlanInput{
		PackageRoot:    packageRoot,
		Manifest:       m,
		ProjectRoot:    packageRoot,
		ConflictAction: "skip",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 1 {
		t.Fatalf("expected 1 op, got %d", len(ops))
	}
	op := ops[0]
	if !strings.HasPrefix(op.TargetPath, customHome) {
		t.Errorf("TargetPath %q should start with COPILOT_HOME %q", op.TargetPath, customHome)
	}
	if op.AllowedRoot != customHome {
		t.Errorf("AllowedRoot = %q, want %q", op.AllowedRoot, customHome)
	}
}
