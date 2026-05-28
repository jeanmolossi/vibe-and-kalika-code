package planner

import "github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"

type Plan struct {
	Operations []platform.PlannedOperation
	Warnings   []platform.Warning
}
