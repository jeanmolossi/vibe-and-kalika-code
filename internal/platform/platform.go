package platform

import "github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"

type Platform string

const (
	PlatformCopilotCLI Platform = "copilot-cli"
	PlatformClaudeCode Platform = "claude-code"
	PlatformCodexCLI   Platform = "codex-cli"
)

type DetectionResult struct {
	Detected   bool
	Platform   Platform
	BasePath   string
	AgentsPath string
	SkillsPath string
	Notes      []string
}

type PlanInput struct {
	PackageRoot    string
	Manifest       *manifest.Manifest
	ProjectRoot    string
	ConflictAction string
}

type ValidateInput struct {
	PackageRoot string
	Manifest    *manifest.Manifest
	ProjectRoot string
}

type OperationType string

const (
	OperationCreate OperationType = "create"
	OperationModify OperationType = "modify"
	OperationSkip   OperationType = "skip"
)

type Conflict struct {
	ExistingPath string
	Action       string
}

type Warning struct {
	Message string
}

type PlannedOperation struct {
	Type        OperationType
	Platform    Platform
	SourcePath  string
	TargetPath  string
	AllowedRoot string // the root this operation is allowed to write to
	AgentName   string // for Codex managed-block: the manifest agent name
	Description string
	Conflict    *Conflict
	Warnings    []Warning
}

type PlatformAdapter interface {
	Platform() Platform
	Detect(projectRoot string) DetectionResult
	Plan(input PlanInput) ([]PlannedOperation, error)
	Validate(input ValidateInput) error
	AllowedRoots(projectRoot string) []string
}
