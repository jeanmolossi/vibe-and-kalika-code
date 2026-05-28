package app

import (
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform/claude"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform/codex"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform/copilot"
)

func Detect(projectRoot string) []platform.DetectionResult {
	adapters := []platform.PlatformAdapter{copilot.NewAdapter(), claude.NewAdapter(), codex.NewAdapter()}
	results := make([]platform.DetectionResult, 0, len(adapters))
	for _, adapter := range adapters {
		results = append(results, adapter.Detect(projectRoot))
	}
	return results
}
