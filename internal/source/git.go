package source

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
)

func CloneGitSource(source, projectRoot string) (*ResolvedSource, error) {
	cloneBase := filepath.Join(projectRoot, ".ai-setup", "sources")
	if err := os.MkdirAll(cloneBase, 0o755); err != nil {
		return nil, fmt.Errorf("create clone base: %w", err)
	}
	dir := filepath.Join(cloneBase, fmt.Sprintf("clone-%d", time.Now().UnixNano()))
	cloneURL := source
	if strings.HasPrefix(source, "file://") {
		cloneURL = strings.TrimPrefix(source, "file://")
	}
	if _, err := git.PlainClone(dir, false, &git.CloneOptions{URL: cloneURL}); err != nil {
		return nil, fmt.Errorf("clone git source: %w", err)
	}
	return &ResolvedSource{Root: dir, Location: source, Cleanup: func() error { return os.RemoveAll(dir) }}, nil
}
