package claude

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/security"
)

type Adapter struct{}

func NewAdapter() *Adapter                     { return &Adapter{} }
func (a *Adapter) Platform() platform.Platform { return platform.PlatformClaudeCode }

func BasePath(projectRoot string) string { return filepath.Join(projectRoot, ".claude") }
func AgentTargetPath(projectRoot, name string) string {
	return filepath.Join(BasePath(projectRoot), "agents", name+".md")
}
func SkillTargetPath(projectRoot, name string) string {
	return filepath.Join(BasePath(projectRoot), "skills", name)
}

func (a *Adapter) Detect(projectRoot string) platform.DetectionResult {
	base := BasePath(projectRoot)
	_, err := os.Stat(base)
	return platform.DetectionResult{Detected: err == nil, Platform: a.Platform(), BasePath: base, AgentsPath: filepath.Join(base, "agents"), SkillsPath: filepath.Join(base, "skills")}
}

func (a *Adapter) AllowedRoots(projectRoot string) []string {
	base := BasePath(projectRoot)
	return []string{filepath.Join(base, "agents"), filepath.Join(base, "skills")}
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
		target := AgentTargetPath(input.ProjectRoot, agent.Name)
		op := platform.PlannedOperation{Platform: a.Platform(), SourcePath: src, TargetPath: target, AllowedRoot: BasePath(input.ProjectRoot), Description: fmt.Sprintf("Install agent %s", agent.Name)}
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
		op := platform.PlannedOperation{Platform: a.Platform(), SourcePath: src, TargetPath: target, AllowedRoot: BasePath(input.ProjectRoot), Description: fmt.Sprintf("Install skill %s", skill.Name)}
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
