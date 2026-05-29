package source

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
)

// ParseGitHubTreeURL parses a GitHub tree URL into its components.
// Input:  https://github.com/owner/repo/tree/branch/sub/dir
// Output: cloneURL="https://github.com/owner/repo.git", branch="branch", subdir="sub/dir"
// Returns ok=false if the URL doesn't match the pattern.
func ParseGitHubTreeURL(rawURL string) (cloneURL, branch, subdir string, ok bool) {
	const prefix = "https://github.com/"
	if !strings.HasPrefix(rawURL, prefix) {
		return "", "", "", false
	}
	const splitIntoTwo = 2
	rest := strings.TrimPrefix(rawURL, prefix)
	parts := strings.SplitN(rest, "/tree/", splitIntoTwo)
	if len(parts) != splitIntoTwo {
		return "", "", "", false
	}
	repoPath := parts[0]  // "owner/repo"
	afterTree := parts[1] // "branch/sub/dir"

	// repoPath must have exactly one slash (owner/repo)
	repoParts := strings.SplitN(repoPath, "/", splitIntoTwo)
	if len(repoParts) != splitIntoTwo || repoParts[0] == "" || repoParts[1] == "" {
		return "", "", "", false
	}

	// Split branch from subdir
	branchAndSub := strings.SplitN(afterTree, "/", splitIntoTwo)
	if len(branchAndSub) == 0 || branchAndSub[0] == "" {
		return "", "", "", false
	}

	cloneURL = "https://github.com/" + repoPath + ".git"
	branch = branchAndSub[0]
	if len(branchAndSub) > 1 {
		subdir = branchAndSub[1]
	}
	return cloneURL, branch, subdir, true
}

func CloneGitSource(source, projectRoot string) (*ResolvedSource, error) {
	stateDir, err := state.StateDir()
	if err != nil {
		return nil, fmt.Errorf("resolve state directory: %w", err)
	}
	cloneBase := filepath.Join(stateDir, "sources")
	if err := os.MkdirAll(cloneBase, 0o755); err != nil {
		return nil, fmt.Errorf("create clone base: %w", err)
	}
	dir := filepath.Join(cloneBase, fmt.Sprintf("clone-%d", time.Now().UnixNano()))

	// Check if this is a GitHub tree URL.
	if ghCloneURL, ghBranch, ghSubdir, ok := ParseGitHubTreeURL(source); ok {
		_, err := git.PlainClone(dir, false, &git.CloneOptions{
			URL:           ghCloneURL,
			ReferenceName: plumbing.NewBranchReferenceName(ghBranch),
			SingleBranch:  true,
			Depth:         1,
		})
		if err != nil {
			return nil, fmt.Errorf("clone git source: %w", err)
		}
		root := dir
		if ghSubdir != "" {
			root = filepath.Join(dir, filepath.FromSlash(ghSubdir))
			// Guard against path traversal: the resolved root must stay within dir.
			dirClean := filepath.Clean(dir) + string(filepath.Separator)
			if !strings.HasPrefix(filepath.Clean(root)+string(filepath.Separator), dirClean) {
				return nil, fmt.Errorf("subdir %q escapes clone root — refusing unsafe path", ghSubdir)
			}
			if _, err := os.Stat(root); err != nil {
				return nil, fmt.Errorf("subdir %q not found after clone: %w", ghSubdir, err)
			}
		}
		return &ResolvedSource{Root: root, Location: source, Cleanup: func() error { return os.RemoveAll(dir) }}, nil
	}

	cloneURL := source
	if strings.HasPrefix(source, "file://") {
		cloneURL = strings.TrimPrefix(source, "file://")
	}
	if _, err := git.PlainClone(dir, false, &git.CloneOptions{URL: cloneURL}); err != nil {
		return nil, fmt.Errorf("clone git source: %w", err)
	}
	return &ResolvedSource{Root: dir, Location: source, Cleanup: func() error { return os.RemoveAll(dir) }}, nil
}
