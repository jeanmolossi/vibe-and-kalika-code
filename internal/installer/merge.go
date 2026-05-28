package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/security"
)

func ManagedBlock(agentName, content string) string {
	trimmed := strings.TrimSpace(content)
	return fmt.Sprintf("<!-- BEGIN VKC AGENT: %s -->\n## Agent: %s\n\n%s\n<!-- END VKC AGENT: %s -->\n", agentName, agentName, trimmed, agentName)
}

func MergeManagedBlock(existing, agentName, content string) string {
	block := ManagedBlock(agentName, content)
	pattern := regexp.MustCompile(`(?s)<!-- BEGIN VKC AGENT: ` + regexp.QuoteMeta(agentName) + ` -->.*?<!-- END VKC AGENT: ` + regexp.QuoteMeta(agentName) + ` -->\n?`)
	if pattern.MatchString(existing) {
		return strings.TrimRight(pattern.ReplaceAllString(existing, block), "\n") + "\n"
	}
	if strings.TrimSpace(existing) == "" {
		return block
	}
	return strings.TrimRight(existing, "\n") + "\n\n" + block
}

func MergeAgentFile(targetPath, agentName, sourcePath, allowedRoot string) error {
	contentBytes, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	existing, _ := os.ReadFile(targetPath)
	merged := MergeManagedBlock(string(existing), agentName, string(contentBytes))
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}
	if allowedRoot != "" {
		if err := security.EnsureResolvedWithinRoot(allowedRoot, targetPath); err != nil {
			return fmt.Errorf("merge destination escapes allowed root: %w", err)
		}
	}
	return os.WriteFile(targetPath, []byte(merged), 0o644)
}

// RemoveManagedBlock removes the managed block for agentName from the file at targetPath.
// If the file does not exist or the block is not present, it returns nil (idempotent).
// If allowedRoot is non-empty, the path is verified to be within that root before any I/O.
func RemoveManagedBlock(targetPath, agentName, allowedRoot string) error {
	if allowedRoot != "" {
		if err := security.EnsureResolvedWithinRoot(allowedRoot, targetPath); err != nil {
			return fmt.Errorf("block path escapes allowed root: %w", err)
		}
	}
	data, err := os.ReadFile(targetPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	pattern := regexp.MustCompile(`(?s)<!-- BEGIN VKC AGENT: ` + regexp.QuoteMeta(agentName) + ` -->.*?<!-- END VKC AGENT: ` + regexp.QuoteMeta(agentName) + ` -->\n?`)
	if !pattern.MatchString(string(data)) {
		return nil
	}
	removed := pattern.ReplaceAllString(string(data), "")
	// Normalize: collapse more than one consecutive blank line into one
	blankLines := regexp.MustCompile(`\n{3,}`)
	normalized := blankLines.ReplaceAllString(removed, "\n\n")
	// Ensure final newline
	normalized = strings.TrimRight(normalized, "\n") + "\n"
	if strings.TrimSpace(normalized) == "" {
		return os.Remove(targetPath)
	}
	return os.WriteFile(targetPath, []byte(normalized), 0o644)
}
