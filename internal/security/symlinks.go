package security

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func RejectEscapingSymlinks(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type()&os.ModeSymlink == 0 {
			return nil
		}
		resolved, err := filepath.EvalSymlinks(path)
		if err != nil {
			return fmt.Errorf("resolve symlink %s: %w", path, err)
		}
		if err := EnsureWithinRoot(root, resolved); err != nil {
			return fmt.Errorf("symlink escapes package root: %s", path)
		}
		return nil
	})
}
