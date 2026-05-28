package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/backup"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/security"
)

type Result struct {
	Applied []platform.PlannedOperation
	Backup  *backup.Result
}

func Apply(projectRoot string, m *manifest.Manifest, operations []platform.PlannedOperation) (*Result, error) {
	var needBackup bool
	for _, op := range operations {
		if op.Type == platform.OperationModify && op.Conflict != nil && op.Conflict.Action == "backup-and-overwrite" {
			needBackup = true
			break
		}
	}
	var backupResult *backup.Result
	var err error
	if needBackup {
		backupResult, err = backup.Create(projectRoot)
		if err != nil {
			return nil, err
		}
	}
	result := &Result{Backup: backupResult}
	for _, op := range operations {
		if err := applyOperation(projectRoot, op); err != nil {
			return nil, err
		}
		if op.Type == platform.OperationSkip {
			result.Applied = append(result.Applied, op)
			continue
		}
		if op.Conflict != nil && op.Conflict.Action == "backup-and-overwrite" {
			if _, err := os.Stat(op.TargetPath); err == nil {
				if err := backup.CopyInto(backupResult, projectRoot, op.TargetPath); err != nil {
					return nil, err
				}
			}
		}
		if strings.HasSuffix(op.TargetPath, "AGENTS.md") {
			agentName := op.AgentName
			if agentName == "" {
				agentName = strings.TrimSuffix(filepath.Base(op.SourcePath), filepath.Ext(op.SourcePath))
			}
			if err := MergeAgentFile(op.TargetPath, agentName, op.SourcePath, op.AllowedRoot); err != nil {
				return nil, err
			}
		} else {
			info, err := os.Stat(op.SourcePath)
			if err != nil {
				return nil, err
			}
			if info.IsDir() {
				if err := os.RemoveAll(op.TargetPath); err != nil && !os.IsNotExist(err) {
					return nil, err
				}
				if err := CopyDir(op.SourcePath, op.TargetPath, op.AllowedRoot); err != nil {
					return nil, err
				}
			} else if err := CopyFile(op.SourcePath, op.TargetPath, op.AllowedRoot); err != nil {
				return nil, err
			}
		}
		result.Applied = append(result.Applied, op)
	}
	if backupResult != nil {
		if err := backup.WriteReport(backupResult); err != nil {
			return nil, err
		}
	}
	sort.Slice(result.Applied, func(i, j int) bool { return result.Applied[i].TargetPath < result.Applied[j].TargetPath })
	_ = m
	_ = time.Now()
	return result, nil
}

func applyOperation(projectRoot string, op platform.PlannedOperation) error {
	if op.Type == platform.OperationSkip {
		return nil
	}
	allowedRoot := op.AllowedRoot
	if allowedRoot == "" {
		allowedRoot = projectRoot
	}
	if err := security.EnsureWithinRoot(allowedRoot, op.TargetPath); err != nil {
		return fmt.Errorf("refusing to write outside allowed roots: %s", op.TargetPath)
	}
	if err := os.MkdirAll(filepath.Dir(op.TargetPath), 0o755); err != nil {
		return err
	}
	// After MkdirAll, verify the resolved path (after symlink expansion) is still within the allowed root.
	return security.EnsureResolvedWithinRoot(allowedRoot, op.TargetPath)
}
