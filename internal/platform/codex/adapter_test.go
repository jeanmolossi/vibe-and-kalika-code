package codex

import (
	"path/filepath"
	"testing"
)

func TestTargetPaths(t *testing.T) {
	root := t.TempDir()
	if got := AgentTargetPath(root); got != filepath.Join(root, "AGENTS.md") {
		t.Fatalf("agent path = %s", got)
	}
	if got := SkillTargetPath(root, "skill"); got != filepath.Join(root, ".agents", "skills", "skill") {
		t.Fatalf("skill path = %s", got)
	}
}
