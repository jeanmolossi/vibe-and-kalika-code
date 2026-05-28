package app

import (
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/planner"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/source"
)

type ValidateResult struct {
	Manifest *manifest.Manifest
	Plan     planner.Plan
}

func ValidateSource(projectRoot, src string) (*ValidateResult, int, error) {
	resolved, err := source.Resolve(src, projectRoot)
	if err != nil {
		return nil, ExitSourceFetchError, err
	}
	defer resolved.Cleanup()
	m, err := manifest.ParseFile(resolved.Root)
	if err != nil {
		return nil, ExitValidationError, err
	}
	if issues, err := manifest.Validate(resolved.Root, projectRoot, m); err != nil {
		return nil, ExitValidationError, errWithIssues(err, issues)
	}
	var targets []platform.Platform
	for _, t := range m.Targets {
		targets = append(targets, platform.Platform(t))
	}
	plan, err := planner.Build(resolved.Root, projectRoot, m, targets, "skip")
	if err != nil {
		return nil, ExitSecurityViolation, err
	}
	return &ValidateResult{Manifest: m, Plan: plan}, ExitSuccess, nil
}
