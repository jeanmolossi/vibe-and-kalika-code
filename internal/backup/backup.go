package backup

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/security"
)

type Result struct {
	Dir   string
	Files []string
}

func Create(projectRoot string) (*Result, error) {
	dir := filepath.Join(projectRoot, ".ai-setup", "backups", time.Now().UTC().Format("20060102-150405"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Result{Dir: dir}, nil
}

func CopyInto(result *Result, projectRoot, source string) error {
	if result == nil {
		return nil
	}
	rel, err := filepath.Rel(projectRoot, source)
	if err != nil || strings.HasPrefix(rel, "..") {
		rel = filepath.Base(source)
	}
	target := filepath.Join(result.Dir, rel)
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			subRel, err := filepath.Rel(source, path)
			if err != nil {
				return err
			}
			dst := filepath.Join(target, subRel)
			if d.IsDir() {
				return os.MkdirAll(dst, 0o755)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
				return err
			}
			if err := security.EnsureResolvedWithinRoot(result.Dir, dst); err != nil {
				return fmt.Errorf("backup destination escapes backup root: %w", err)
			}
			if err := os.WriteFile(dst, data, 0o644); err != nil {
				return err
			}
			result.Files = append(result.Files, dst)
			return nil
		})
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	if err := security.EnsureResolvedWithinRoot(result.Dir, target); err != nil {
		return fmt.Errorf("backup destination escapes backup root: %w", err)
	}
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		return err
	}
	result.Files = append(result.Files, target)
	return nil
}

func WriteReport(result *Result) error {
	if result == nil {
		return nil
	}
	var b strings.Builder
	b.WriteString("Backup report\n")
	for _, file := range result.Files {
		b.WriteString("- " + file + "\n")
	}
	return os.WriteFile(filepath.Join(result.Dir, "backup-report.txt"), []byte(b.String()), 0o644)
}

func Summary(result *Result) string {
	if result == nil {
		return ""
	}
	return fmt.Sprintf("%s (%d files)", result.Dir, len(result.Files))
}
