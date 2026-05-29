package app

import (
	"fmt"
	"os"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/backup"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/installer"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/security"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
)

// SkippedItem describes a file or managed block that was not removed during
// Uninstall because it resides outside the project root and Force was false.
type SkippedItem struct {
	Path      string // absolute path to the file or AGENTS.md containing the block
	AgentName string // non-empty when the item is a managed agent block
}

// UninstallOptions controls the behavior of Uninstall.
type UninstallOptions struct {
	Package     string
	ProjectRoot string
	DryRun      bool
	// Force allows removing files outside the project root without prompting.
	// When false, such files are collected in UninstallResult.FilesSkipped instead.
	Force bool
}

// UninstallResult reports what was done (or would be done on dry-run).
type UninstallResult struct {
	Package        string
	FilesRemoved   []string
	FilesRestored  []string
	MarkersRemoved []string
	// FilesSkipped contains items outside ProjectRoot that were not removed.
	// Non-empty only when Force is false and such items exist.
	FilesSkipped []SkippedItem
	BackupPath   string
}

// Uninstall removes an installed package from the project.
func Uninstall(opts UninstallOptions) (*UninstallResult, int, error) {
	st, err := state.Read()
	if err != nil {
		return nil, ExitError, err
	}

	recIdx := -1
	for i, inst := range st.Installations {
		if inst.Package == opts.Package {
			recIdx = i
			break
		}
	}
	if recIdx == -1 {
		return nil, ExitError, fmt.Errorf("package %q is not installed", opts.Package)
	}
	rec := st.Installations[recIdx]

	result := &UninstallResult{
		Package:    rec.Package,
		BackupPath: rec.BackupPath,
	}

	if opts.DryRun {
		result.FilesRemoved = rec.CreatedFiles
		for _, block := range rec.AgentBlocks {
			result.MarkersRemoved = append(result.MarkersRemoved, block.Path)
		}
		if rec.BackupPath != "" {
			result.FilesRestored = rec.Files
		}
		return result, ExitSuccess, nil
	}

	// Remove managed blocks from AGENTS.md files.
	for _, block := range rec.AgentBlocks {
		withinRoot := security.EnsureWithinRoot(opts.ProjectRoot, block.Path) == nil
		if !withinRoot {
			if !opts.Force {
				result.FilesSkipped = append(result.FilesSkipped, SkippedItem{
					Path:      block.Path,
					AgentName: block.AgentName,
				})
				continue
			}
			// Force: pass empty allowedRoot to bypass the root restriction.
			if err := installer.RemoveManagedBlock(block.Path, block.AgentName, ""); err != nil {
				return nil, ExitError, fmt.Errorf("remove managed block %s: %w", block.Path, err)
			}
		} else {
			if err := installer.RemoveManagedBlock(block.Path, block.AgentName, opts.ProjectRoot); err != nil {
				return nil, ExitError, fmt.Errorf("remove managed block %s: %w", block.Path, err)
			}
		}
		result.MarkersRemoved = append(result.MarkersRemoved, block.Path)
	}

	// Restore backup if available.
	if rec.BackupPath != "" {
		if err := backup.Restore(rec.BackupPath, opts.ProjectRoot); err != nil {
			return nil, ExitError, fmt.Errorf("restore backup: %w", err)
		}
		result.FilesRestored = rec.Files
	}

	// Delete created files and directories.
	for _, f := range rec.CreatedFiles {
		withinRoot := security.EnsureWithinRoot(opts.ProjectRoot, f) == nil
		if !withinRoot {
			if !opts.Force {
				result.FilesSkipped = append(result.FilesSkipped, SkippedItem{Path: f})
				continue
			}
		}
		if removed, rmErr := removeFileOrDir(f); rmErr != nil {
			return nil, ExitError, rmErr
		} else if removed {
			result.FilesRemoved = append(result.FilesRemoved, f)
		}
	}

	// Remove installation record.
	st.Installations = append(st.Installations[:recIdx], st.Installations[recIdx+1:]...)
	if err := state.Write(st); err != nil {
		return nil, ExitError, err
	}

	return result, ExitSuccess, nil
}

// RemoveSkippedItems forcefully removes files and managed blocks collected in
// UninstallResult.FilesSkipped. It bypasses the project-root security check.
// Callers must obtain explicit user consent before calling this function.
func RemoveSkippedItems(items []SkippedItem) (filesRemoved, markersRemoved []string, err error) {
	for _, item := range items {
		if item.AgentName != "" {
			// Pass empty allowedRoot to bypass the project-root restriction.
			if rmErr := installer.RemoveManagedBlock(item.Path, item.AgentName, ""); rmErr != nil {
				return filesRemoved, markersRemoved, fmt.Errorf("remove block %s: %w", item.Path, rmErr)
			}
			markersRemoved = append(markersRemoved, item.Path)
		} else {
			if removed, rmErr := removeFileOrDir(item.Path); rmErr != nil {
				return filesRemoved, markersRemoved, rmErr
			} else if removed {
				filesRemoved = append(filesRemoved, item.Path)
			}
		}
	}
	return filesRemoved, markersRemoved, nil
}

// removeFileOrDir removes f if it exists.
// Returns (true, nil) on success, (false, nil) if the path does not exist, or (false, err) on failure.
func removeFileOrDir(f string) (bool, error) {
	info, statErr := os.Stat(f)
	if os.IsNotExist(statErr) {
		return false, nil
	}
	if statErr != nil {
		return false, fmt.Errorf("stat %s: %w", f, statErr)
	}
	if info.IsDir() {
		if err := os.RemoveAll(f); err != nil {
			return false, fmt.Errorf("remove %s: %w", f, err)
		}
	} else {
		if err := os.Remove(f); err != nil {
			if os.IsNotExist(err) {
				return false, nil
			}
			return false, fmt.Errorf("remove %s: %w", f, err)
		}
	}
	return true, nil
}
