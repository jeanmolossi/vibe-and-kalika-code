package security

import (
	"fmt"
	"path/filepath"
	"strings"
)

func ValidateRelativePath(path string) error {
	if path == "" {
		return fmt.Errorf("path is required")
	}
	if filepath.IsAbs(path) {
		return fmt.Errorf("absolute paths are not allowed: %s", path)
	}
	clean := filepath.Clean(path)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path traversal is not allowed: %s", path)
	}
	return nil
}

func ResolveWithinRoot(root, rel string) (string, error) {
	if err := ValidateRelativePath(rel); err != nil {
		return "", err
	}
	resolved := filepath.Join(root, filepath.Clean(rel))
	if err := EnsureWithinRoot(root, resolved); err != nil {
		return "", err
	}
	return resolved, nil
}

func EnsureWithinRoot(root, target string) error {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve root: %w", err)
	}
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("resolve target: %w", err)
	}
	rel, err := filepath.Rel(rootAbs, targetAbs)
	if err != nil {
		return fmt.Errorf("compare path: %w", err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path escapes allowed root: %s", target)
	}
	return nil
}

// EnsureResolvedWithinRoot resolves symlinks in the target path and verifies
// the resolved path is still within root. For files that don't exist yet, it
// resolves the parent directory. Use this before any write operation to prevent
// symlink escape attacks.
func EnsureResolvedWithinRoot(root, target string) error {
	var resolved string
	if r, err := filepath.EvalSymlinks(target); err == nil {
		resolved = r
	} else {
		// target doesn't exist yet; resolve the parent directory
		parent := filepath.Dir(target)
		resolvedParent, err := filepath.EvalSymlinks(parent)
		if err != nil {
			// parent doesn't exist either; fall back to raw path check
			return EnsureWithinRoot(root, target)
		}
		resolved = filepath.Join(resolvedParent, filepath.Base(target))
	}
	if err := EnsureWithinRoot(root, resolved); err != nil {
		return fmt.Errorf("resolved path %q escapes allowed root %q", resolved, root)
	}
	return nil
}
