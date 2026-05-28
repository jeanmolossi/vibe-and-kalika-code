package planner

import (
	"fmt"
	"sort"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform/claude"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform/codex"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform/copilot"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/security"
)

func Build(packageRoot, projectRoot string, m *manifest.Manifest, selected []platform.Platform, conflictAction string) (Plan, error) {
	adapters := []platform.PlatformAdapter{}
	want := map[platform.Platform]struct{}{}
	for _, p := range selected {
		want[p] = struct{}{}
	}
	for _, adapter := range []platform.PlatformAdapter{copilot.NewAdapter(), claude.NewAdapter(), codex.NewAdapter()} {
		if _, ok := want[adapter.Platform()]; ok {
			adapters = append(adapters, adapter)
		}
	}
	var plan Plan
	for _, adapter := range adapters {
		ops, err := adapter.Plan(platform.PlanInput{PackageRoot: packageRoot, Manifest: m, ProjectRoot: projectRoot, ConflictAction: conflictAction})
		if err != nil {
			return Plan{}, err
		}
		for _, op := range ops {
			for _, root := range adapter.AllowedRoots(projectRoot) {
				if err := security.EnsureWithinRoot(root, op.TargetPath); err == nil {
					goto safe
				}
			}
			return Plan{}, fmt.Errorf("target path outside allowed roots: %s", op.TargetPath)
		safe:
			if scripts, err := security.FindScripts(op.SourcePath); err == nil {
				for _, s := range scripts {
					op.Warnings = append(op.Warnings, platform.Warning{Message: fmt.Sprintf("script detected: %s", s)})
				}
			} else if security.IsScriptPath(op.SourcePath) {
				op.Warnings = append(op.Warnings, platform.Warning{Message: fmt.Sprintf("script detected: %s", op.SourcePath)})
			}
			plan.Operations = append(plan.Operations, op)
		}
	}
	sort.Slice(plan.Operations, func(i, j int) bool { return plan.Operations[i].TargetPath < plan.Operations[j].TargetPath })
	for _, op := range plan.Operations {
		plan.Warnings = append(plan.Warnings, op.Warnings...)
	}
	return plan, nil
}
