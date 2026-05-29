package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/installer"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/planner"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/report"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/source"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
)

const (
	ExitSuccess            = 0
	ExitError              = 1
	ExitValidationError    = 2
	ExitSecurityViolation  = 3
	ExitUserCancelled      = 4
	ExitSourceFetchError   = 5
	ExitConflictUnresolved = 6
)

type InstallOptions struct {
	Source         string
	ProjectRoot    string
	Targets        []platform.Platform
	Yes            bool
	ConflictAction string
	DryRun         bool
}

type InstallResult struct {
	Manifest   *manifest.Manifest
	Plan       planner.Plan
	ReportPath string
	BackupPath string
}

func Install(opts InstallOptions) (*InstallResult, int, error) {
	resolved, err := source.Resolve(opts.Source, opts.ProjectRoot)
	if err != nil {
		return nil, ExitSourceFetchError, err
	}
	defer resolved.Cleanup() //nolint:errcheck // temp directory cleanup error is non-actionable
	m, err := manifest.ParseFile(resolved.Root)
	if err != nil {
		return nil, ExitValidationError, err
	}
	if issues, err := manifest.Validate(resolved.Root, opts.ProjectRoot, m); err != nil {
		return nil, ExitValidationError, fmt.Errorf("%w: %s", err, strings.Join(issues, "; "))
	}
	targets := opts.Targets
	if len(targets) == 0 {
		for _, t := range m.Targets {
			targets = append(targets, platform.Platform(t))
		}
	}
	plan, err := planner.Build(resolved.Root, opts.ProjectRoot, m, targets, opts.ConflictAction)
	if err != nil {
		return nil, ExitSecurityViolation, err
	}
	if planner.HasConflicts(plan) && opts.ConflictAction == "skip" {
		planner.ApplyConflictAction(&plan, "skip")
	}
	if planner.HasConflicts(plan) && opts.ConflictAction == "" {
		return &InstallResult{Manifest: m, Plan: plan}, ExitConflictUnresolved, fmt.Errorf("conflicts detected")
	}
	if opts.DryRun {
		return &InstallResult{Manifest: m, Plan: plan}, ExitSuccess, nil
	}
	instResult, err := installer.Apply(opts.ProjectRoot, m, plan.Operations)
	if err != nil {
		return nil, ExitError, err
	}
	reportPath, err := report.WriteInstallReport(report.InstallReportInput{
		Manifest:   m,
		Source:     opts.Source,
		Platforms:  targets,
		Operations: instResult.Applied,
		Backup:     instResult.Backup,
	})
	if err != nil {
		return nil, ExitError, err
	}
	st, err := state.Read()
	if err != nil {
		return nil, ExitError, err
	}
	installRecord := state.Installation{
		Package:     m.Name,
		Version:     m.Version,
		Source:      opts.Source,
		InstalledAt: time.Now().UTC().Format(time.RFC3339),
		ReportPath:  reportPath,
	}
	if instResult.Backup != nil {
		installRecord.BackupPath = instResult.Backup.Dir
	}
	for _, p := range targets {
		installRecord.Platforms = append(installRecord.Platforms, string(p))
	}
	for _, op := range instResult.Applied {
		installRecord.Files = append(installRecord.Files, op.TargetPath)
		if strings.HasSuffix(op.TargetPath, "AGENTS.md") {
			installRecord.ManagedMarkers = append(installRecord.ManagedMarkers, op.TargetPath)
			agentName := op.AgentName
			if agentName == "" {
				agentName = m.Name
			}
			installRecord.AgentBlocks = append(installRecord.AgentBlocks, state.AgentBlock{
				Path:      op.TargetPath,
				AgentName: agentName,
			})
		}
	}
	for _, op := range instResult.Applied {
		if op.Type == platform.OperationCreate && !strings.HasSuffix(op.TargetPath, "AGENTS.md") {
			installRecord.CreatedFiles = append(installRecord.CreatedFiles, op.TargetPath)
		}
	}
	st.Installations = append(st.Installations, installRecord)
	if err := state.Write(st); err != nil {
		return nil, ExitError, err
	}
	result := &InstallResult{Manifest: m, Plan: plan, ReportPath: reportPath}
	if instResult.Backup != nil {
		result.BackupPath = instResult.Backup.Dir
	}
	return result, ExitSuccess, nil
}
