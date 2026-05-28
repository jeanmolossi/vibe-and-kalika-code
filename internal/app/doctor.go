package app

import (
	"os/exec"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
)

// PlatformStatus holds doctor info for a single platform.
type PlatformStatus struct {
	Name     string
	Platform platform.Platform
	Detected bool
	BasePath string
}

// DoctorResult holds the full environment health report.
type DoctorResult struct {
	Platforms    []PlatformStatus
	GitAvailable bool
	StateValid   bool
}

func Doctor(projectRoot string) DoctorResult {
	_, gitErr := exec.LookPath("git")
	_, stateErr := state.Read(projectRoot)
	detections := Detect(projectRoot)

	platforms := make([]PlatformStatus, 0, len(detections))
	for _, d := range detections {
		platforms = append(platforms, PlatformStatus{
			Name:     platformName(d.Platform),
			Platform: d.Platform,
			Detected: d.Detected,
			BasePath: d.BasePath,
		})
	}
	return DoctorResult{
		Platforms:    platforms,
		GitAvailable: gitErr == nil,
		StateValid:   stateErr == nil,
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
