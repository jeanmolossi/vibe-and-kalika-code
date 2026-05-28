package claude

import (
	"path/filepath"
	"testing"
)

func TestTargetPaths(t *testing.T) {
	root := t.TempDir()
	if got := AgentTargetPath(root, "agent"); got != filepath.Join(root, ".claude", "agents", "agent.md") {
		t.Fatalf("agent path = %s", got)
	}
	if got := SkillTargetPath(root, "skill"); got != filepath.Join(root, ".claude", "skills", "skill") {
		t.Fatalf("skill path = %s", got)
	}
}
