package app

import (
	"os/exec"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
)

type DoctorResult struct {
	Platforms    []struct{ Detected bool }
	GitAvailable bool
	StateValid   bool
}

func Doctor(projectRoot string) DoctorResult {
	_, gitErr := exec.LookPath("git")
	_, stateErr := state.Read(projectRoot)
	detections := Detect(projectRoot)
	platforms := make([]struct{ Detected bool }, 0, len(detections))
	for _, detection := range detections {
		platforms = append(platforms, struct{ Detected bool }{Detected: detection.Detected})
	}
	return DoctorResult{Platforms: platforms, GitAvailable: gitErr == nil, StateValid: stateErr == nil}
}
