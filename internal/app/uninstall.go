package app

import (
	"fmt"
	"os"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/backup"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/installer"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/security"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
)

// UninstallOptions controls the behavior of Uninstall.
type UninstallOptions struct {
	Package     string
	ProjectRoot string
	DryRun      bool
}

// UninstallResult reports what was done (or would be done on dry-run).
type UninstallResult struct {
	Package        string
	FilesRemoved   []string
	FilesRestored  []string
	MarkersRemoved []string
	BackupPath     string
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

	// Remove managed blocks from AGENTS.md files
	for _, block := range rec.AgentBlocks {
		if err := security.EnsureWithinRoot(opts.ProjectRoot, block.Path); err != nil {
			return nil, ExitSecurityViolation, fmt.Errorf("block path %s escapes project root: %w", block.Path, err)
		}
		if err := installer.RemoveManagedBlock(block.Path, block.AgentName, opts.ProjectRoot); err != nil {
			return nil, ExitError, fmt.Errorf("remove managed block %s: %w", block.Path, err)
		}
		result.MarkersRemoved = append(result.MarkersRemoved, block.Path)
	}

	// Restore backup if available
	if rec.BackupPath != "" {
		if err := backup.Restore(rec.BackupPath, opts.ProjectRoot); err != nil {
			return nil, ExitError, fmt.Errorf("restore backup: %w", err)
		}
		result.FilesRestored = rec.Files
	}

	// Delete created files and directories
	for _, f := range rec.CreatedFiles {
		if err := security.EnsureWithinRoot(opts.ProjectRoot, f); err != nil {
			return nil, ExitSecurityViolation, fmt.Errorf("file %s escapes project root: %w", f, err)
		}
		info, statErr := os.Stat(f)
		if os.IsNotExist(statErr) {
			continue
		}
		if statErr != nil {
			return nil, ExitError, fmt.Errorf("stat %s: %w", f, statErr)
		}
		if info.IsDir() {
			if err := os.RemoveAll(f); err != nil {
				return nil, ExitError, fmt.Errorf("remove %s: %w", f, err)
			}
		} else {
			if err := os.Remove(f); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return nil, ExitError, fmt.Errorf("remove %s: %w", f, err)
			}
		}
		result.FilesRemoved = append(result.FilesRemoved, f)
	}

	// Remove installation record
	st.Installations = append(st.Installations[:recIdx], st.Installations[recIdx+1:]...)
	if err := state.Write(st); err != nil {
		return nil, ExitError, err
	}

	return result, ExitSuccess, nil
}
