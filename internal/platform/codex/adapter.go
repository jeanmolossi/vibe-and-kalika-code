package codex

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/security"
)

type Adapter struct{}

func NewAdapter() *Adapter                     { return &Adapter{} }
func (a *Adapter) Platform() platform.Platform { return platform.PlatformCodexCLI }

func AgentTargetPath(projectRoot string) string { return filepath.Join(projectRoot, "AGENTS.md") }
func SkillTargetPath(projectRoot, name string) string {
	return filepath.Join(projectRoot, ".agents", "skills", name)
}

func (a *Adapter) Detect(projectRoot string) platform.DetectionResult {
	skills := filepath.Join(projectRoot, ".agents", "skills")
	agents := AgentTargetPath(projectRoot)
	_, skillsErr := os.Stat(skills)
	_, agentsErr := os.Stat(agents)
	return platform.DetectionResult{
		Detected:   skillsErr == nil || agentsErr == nil,
		Platform:   a.Platform(),
		BasePath:   projectRoot,
		AgentsPath: agents,
		SkillsPath: skills,
		Notes:      []string{"Agents are merged into AGENTS.md managed sections."},
	}
}

func (a *Adapter) AllowedRoots(projectRoot string) []string {
	return []string{AgentTargetPath(projectRoot), filepath.Join(projectRoot, ".agents", "skills")}
}

func (a *Adapter) Validate(input platform.ValidateInput) error { _ = input; return nil }

func (a *Adapter) Plan(input platform.PlanInput) ([]platform.PlannedOperation, error) {
	var ops []platform.PlannedOperation
	for _, agent := range input.Manifest.Agents {
		if _, ok := agent.Targets[string(a.Platform())]; !ok {
			continue
		}
		src, err := security.ResolveWithinRoot(input.PackageRoot, agent.Source)
		if err != nil {
			return nil, err
		}
		target := AgentTargetPath(input.ProjectRoot)
		op := platform.PlannedOperation{
			Platform:    a.Platform(),
			SourcePath:  src,
			TargetPath:  target,
			AllowedRoot: input.ProjectRoot,
			AgentName:   agent.Name,
			Description: fmt.Sprintf("Merge agent %s into AGENTS.md", agent.Name),
		}
		if _, err := os.Stat(target); err == nil {
			op.Type = platform.OperationModify
			op.Conflict = &platform.Conflict{ExistingPath: target, Action: input.ConflictAction}
		} else {
			op.Type = platform.OperationCreate
		}
		ops = append(ops, op)
	}
	for _, skill := range input.Manifest.Skills {
		if _, ok := skill.Targets[string(a.Platform())]; !ok {
			continue
		}
		src, err := security.ResolveWithinRoot(input.PackageRoot, skill.Source)
		if err != nil {
			return nil, err
		}
		target := SkillTargetPath(input.ProjectRoot, skill.Name)
		op := platform.PlannedOperation{
			Platform:    a.Platform(),
			SourcePath:  src,
			TargetPath:  target,
			AllowedRoot: filepath.Join(input.ProjectRoot, ".agents"),
			Description: fmt.Sprintf("Install skill %s", skill.Name),
		}
		if _, err := os.Stat(target); err == nil {
			op.Type = platform.OperationModify
			op.Conflict = &platform.Conflict{ExistingPath: target, Action: input.ConflictAction}
		} else {
			op.Type = platform.OperationCreate
		}
		ops = append(ops, op)
	}
	return ops, nil
}

var _ platform.PlatformAdapter = (*Adapter)(nil)
