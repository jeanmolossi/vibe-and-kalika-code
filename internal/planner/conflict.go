package planner

import "github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"

func HasConflicts(plan Plan) bool {
	for _, op := range plan.Operations {
		if op.Conflict != nil {
			return true
		}
	}
	return false
}

func ApplyConflictAction(plan *Plan, action string) {
	for i := range plan.Operations {
		if plan.Operations[i].Conflict != nil {
			plan.Operations[i].Conflict.Action = action
			if action == "skip" {
				plan.Operations[i].Type = platform.OperationSkip
			}
		}
	}
}
