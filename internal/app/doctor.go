package app

import (
	"os"
	"os/exec"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
)

// PlatformStatus holds doctor info for a single platform.
type PlatformStatus struct {
	Name      string
	Platform  platform.Platform
	Detected  bool
	BasePath  string
	AgentsDir string
	SkillsDir string
	Writable  bool
}

// DoctorResult holds the full environment health report.
type DoctorResult struct {
	ProjectRoot  string
	Platforms    []PlatformStatus
	GitAvailable bool
	StateValid   bool
	EnvVars      map[string]string
}

// isWritable reports whether path is writable by attempting to create a temp file there.
func isWritable(path string) bool {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return false
	}
	tmp, err := os.CreateTemp(path, ".vkc-write-check-*")
	if err != nil {
		return false
	}
	_ = tmp.Close()
	_ = os.Remove(tmp.Name())
	return true
}

// Doctor runs environment health checks for all supported platforms.
func Doctor(projectRoot string) DoctorResult {
	_, gitErr := exec.LookPath("git")
	_, stateErr := state.Read(projectRoot)
	detections := Detect(projectRoot)

	platforms := make([]PlatformStatus, 0, len(detections))
	for _, d := range detections {
		platforms = append(platforms, PlatformStatus{
			Name:      platformName(d.Platform),
			Platform:  d.Platform,
			Detected:  d.Detected,
			BasePath:  d.BasePath,
			AgentsDir: d.AgentsPath,
			SkillsDir: d.SkillsPath,
			Writable:  isWritable(d.AgentsPath),
		})
	}
	return DoctorResult{
		ProjectRoot:  projectRoot,
		Platforms:    platforms,
		GitAvailable: gitErr == nil,
		StateValid:   stateErr == nil,
		EnvVars: map[string]string{
			"COPILOT_HOME": os.Getenv("COPILOT_HOME"),
			"CODEX_HOME":   os.Getenv("CODEX_HOME"),
		},
	}
}

func platformName(p platform.Platform) string {
	switch p {
	case platform.PlatformCopilotCLI:
		return "GitHub Copilot CLI"
	case platform.PlatformClaudeCode:
		return "Claude Code"
	case platform.PlatformCodexCLI:
		return "OpenAI Codex CLI"
	default:
		return string(p)
	}
}
