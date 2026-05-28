package source

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const DefaultSource = "https://github.com/jeanmolossi/vibe-and-kalika-code/tree/main/packages/kalika-ofc"

type ResolvedSource struct {
	Root     string
	Location string
	Cleanup  func() error
}

func Resolve(source, projectRoot string) (*ResolvedSource, error) {
	if IsGitURL(source) {
		return CloneGitSource(source, projectRoot)
	}
	return ResolveLocalSource(source)
}

func IsGitURL(source string) bool {
	return strings.HasPrefix(source, "http://") ||
		strings.HasPrefix(source, "https://") ||
		strings.HasPrefix(source, "git@") ||
		strings.HasPrefix(source, "ssh://") ||
		strings.HasPrefix(source, "file://")
}

func ResolveLocalSource(source string) (*ResolvedSource, error) {
	root, err := filepath.Abs(source)
	if err != nil {
		return nil, fmt.Errorf("resolve source: %w", err)
	}
	if _, err := os.Stat(root); err != nil {
		return nil, fmt.Errorf("stat source: %w", err)
	}
	return &ResolvedSource{Root: root, Location: root, Cleanup: func() error { return nil }}, nil
}
