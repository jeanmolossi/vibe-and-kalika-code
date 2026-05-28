package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/planner"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/source"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/ui"
)

// Init runs the interactive setup wizard for vkc init.
// src pre-fills the source step when the user passes a positional argument.
// yes skips all interactive prompts and uses sensible defaults.
func Init(projectRoot, src string, yes bool) (*InstallResult, int, error) {
	// Step 1: welcome banner.
	fmt.Println(welcomeBanner())

	// Step 2: detect project root.
	if projectRoot == "" {
		var err error
		projectRoot, err = os.Getwd()
		if err != nil {
			return nil, ExitError, fmt.Errorf("detect project root: %w", err)
		}
	}

	// Step 3: detect supported platforms.
	detections := Detect(projectRoot)

	// Step 4: ask which platforms to target.
	selectedPlatforms, err := ui.AskPlatforms(detections, yes)
	if err != nil {
		return nil, ExitError, fmt.Errorf("platform selection: %w", err)
	}
	if len(selectedPlatforms) == 0 {
		return nil, ExitUserCancelled, fmt.Errorf("no platforms selected — installation cancelled")
	}

	// Step 5: ask for the package source.
	sourceInput, err := ui.AskSource(yes, src)
	if err != nil {
		return nil, ExitUserCancelled, err
	}
	if sourceInput.Source == "" {
		return nil, ExitUserCancelled, fmt.Errorf("no source provided — installation cancelled")
	}

	// Step 6: resolve (fetch) source.
	resolved, err := source.Resolve(sourceInput.Source, projectRoot)
	if err != nil {
		return nil, ExitSourceFetchError, fmt.Errorf("fetch source: %w", err)
	}
	defer resolved.Cleanup() //nolint:errcheck

	// Step 7: parse manifest.
	m, err := manifest.ParseFile(resolved.Root)
	if err != nil {
		return nil, ExitValidationError, fmt.Errorf("parse manifest: %w", err)
	}

	// Step 8: validate manifest.
	if issues, verr := manifest.Validate(resolved.Root, projectRoot, m); verr != nil {
		return nil, ExitValidationError, fmt.Errorf("%w: %s", verr, strings.Join(issues, "; "))
	}

	// Step 9: show package summary.
	ui.ShowPackageSummary(m)

	// Step 10: build dry-run plan (empty conflict action so conflicts surface).
	plan, err := planner.Build(resolved.Root, projectRoot, m, selectedPlatforms, "")
	if err != nil {
		return nil, ExitSecurityViolation, fmt.Errorf("build plan: %w", err)
	}

	// Step 11: show dry-run plan.
	ui.ShowDryRun(plan.Operations)

	// Step 12: resolve conflicts — ask per-conflict, derive a global action.
	conflictAction := "backup-and-overwrite"
	if planner.HasConflicts(plan) {
		actionCounts := make(map[string]int)
		for _, op := range plan.Operations {
			if op.Conflict == nil {
				continue
			}
			chosen, cerr := ui.AskConflictAction(op.TargetPath, yes)
			if cerr != nil {
				return nil, ExitError, fmt.Errorf("conflict resolution: %w", cerr)
			}
			actionCounts[chosen]++
		}
		// Use the single chosen action when every conflict was resolved the same way;
		// otherwise fall back to the safe "backup-and-overwrite" default.
		if len(actionCounts) == 1 {
			for a := range actionCounts {
				conflictAction = a
			}
		}
	}

	// Step 13: final confirmation before touching the filesystem.
	confirmed, err := ui.Confirm("Apply installation?", yes)
	if err != nil {
		return nil, ExitError, fmt.Errorf("confirmation: %w", err)
	}
	if !confirmed {
		return nil, ExitUserCancelled, fmt.Errorf("installation cancelled by user")
	}

	// Steps 14-17: install.
	// Pass resolved.Root (local path) so Install re-resolves cheaply without re-fetching.
	result, code, err := Install(InstallOptions{
		Source:         resolved.Root,
		ProjectRoot:    projectRoot,
		Targets:        selectedPlatforms,
		Yes:            yes,
		ConflictAction: conflictAction,
	})
	if err != nil {
		return nil, code, err
	}

	// Step 18: show final summary.
	ui.ShowFinalSummary(buildInstallSummary(result))

	return result, ExitSuccess, nil
}

// buildInstallSummary converts an InstallResult into the UI-safe InstallSummary.
func buildInstallSummary(result *InstallResult) ui.InstallSummary {
	s := ui.InstallSummary{
		ReportPath: result.ReportPath,
		BackupPath: result.BackupPath,
	}
	for _, op := range result.Plan.Operations {
		switch op.Type {
		case platform.OperationCreate:
			s.FilesCreated = append(s.FilesCreated, op.TargetPath)
		case platform.OperationModify:
			s.FilesModified = append(s.FilesModified, op.TargetPath)
		case platform.OperationSkip:
			s.FilesSkipped = append(s.FilesSkipped, op.TargetPath)
		}
	}
	return s
}

// welcomeBanner returns a short ASCII welcome message.
func welcomeBanner() string {
	return `
╔══════════════════════════════════════╗
║   vkc — AI Agent Package Manager    ║
║           Setup Wizard               ║
╚══════════════════════════════════════╝`
}
