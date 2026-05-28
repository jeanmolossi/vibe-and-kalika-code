package app

import "fmt"

func Init(projectRoot, src string, yes bool) (*InstallResult, int, error) {
	if src == "" {
		return nil, ExitUserCancelled, fmt.Errorf("source is required for non-interactive init")
	}
	return Install(InstallOptions{Source: src, ProjectRoot: projectRoot, Yes: yes, ConflictAction: "backup-and-overwrite"})
}
