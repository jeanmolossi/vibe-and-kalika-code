package installer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/installer"
)

func TestMergeManagedBlockInsertUpdatePreserve(t *testing.T) {
	existing := "# Existing\n"
	merged := installer.MergeManagedBlock(existing, "agent-one", "hello")
	if !strings.Contains(merged, "BEGIN VKC AGENT: agent-one") || !strings.Contains(merged, "# Existing") {
		t.Fatalf("unexpected merged content: %s", merged)
	}
	updated := installer.MergeManagedBlock(merged, "agent-one", "updated")
	if strings.Count(updated, "BEGIN VKC AGENT: agent-one") != 1 || !strings.Contains(updated, "updated") || strings.Contains(updated, "hello") {
		t.Fatalf("unexpected updated content: %s", updated)
	}
}

// TestMergeAgentFileUsesAgentName verifies that the managed block key is derived
// from the agent name (from the manifest), not from the source filename.
func TestMergeAgentFileUsesAgentName(t *testing.T) {
	dir := t.TempDir()
	// source file intentionally has a different name than the agent
	srcPath := filepath.Join(dir, "differently-named.md")
	if err := os.WriteFile(srcPath, []byte("agent content"), 0o644); err != nil {
		t.Fatal(err)
	}
	targetPath := filepath.Join(dir, "AGENTS.md")

	// First install: block should use agent name, not source filename
	if err := installer.MergeAgentFile(targetPath, "my-custom-agent", srcPath, dir); err != nil {
		t.Fatal(err)
	}
	content, _ := os.ReadFile(targetPath)
	if !strings.Contains(string(content), "BEGIN VKC AGENT: my-custom-agent") {
		t.Fatalf("block should use agent name 'my-custom-agent', got:\n%s", content)
	}
	if strings.Contains(string(content), "differently-named") {
		t.Fatalf("block must not contain source filename 'differently-named', got:\n%s", content)
	}

	// Second install (update): block should be replaced, not duplicated
	if err := installer.MergeAgentFile(targetPath, "my-custom-agent", srcPath, dir); err != nil {
		t.Fatal(err)
	}
	content, _ = os.ReadFile(targetPath)
	if count := strings.Count(string(content), "BEGIN VKC AGENT: my-custom-agent"); count != 1 {
		t.Fatalf("block should appear exactly once after update, got %d occurrences:\n%s", count, content)
	}
}
