package report

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/backup"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
)

type InstallReportInput struct {
	ProjectRoot string
	Manifest    *manifest.Manifest
	Source      string
	Platforms   []platform.Platform
	Operations  []platform.PlannedOperation
	Backup      *backup.Result
}

func WriteInstallReport(input InstallReportInput) (string, error) {
	dir := filepath.Join(input.ProjectRoot, ".ai-setup", "reports")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, time.Now().UTC().Format("20060102-150405")+"-install-report.md")
	content := RenderMarkdown(input)
	return path, os.WriteFile(path, []byte(content), 0o644)
}

func RenderMarkdown(input InstallReportInput) string {
	var b strings.Builder
	b.WriteString("# VKC Install Report\n\n")
	b.WriteString("- Package: " + input.Manifest.Name + "\n")
	b.WriteString("- Version: " + input.Manifest.Version + "\n")
	b.WriteString("- Source: " + input.Source + "\n")
	b.WriteString("- Platforms: " + joinPlatforms(input.Platforms) + "\n")
	if input.Backup != nil {
		b.WriteString("- Backup path: " + input.Backup.Dir + "\n")
	}

	// Partition operations by type
	sections := []struct {
		heading string
		opType  platform.OperationType
	}{
		{"## Created Files", platform.OperationCreate},
		{"## Modified Files", platform.OperationModify},
		{"## Skipped Files", platform.OperationSkip},
	}

	for _, s := range sections {
		var lines []string
		for _, op := range input.Operations {
			if op.Type == s.opType {
				lines = append(lines, "- ["+string(op.Type)+"] "+op.TargetPath)
			}
		}
		if len(lines) > 0 {
			b.WriteString("\n" + s.heading + "\n")
			for _, l := range lines {
				b.WriteString(l + "\n")
			}
		}
	}

	// Conflicts section
	var conflicts []platform.PlannedOperation
	for _, op := range input.Operations {
		if op.Conflict != nil {
			conflicts = append(conflicts, op)
		}
	}
	if len(conflicts) > 0 {
		b.WriteString("\n## Conflicts\n")
		for _, op := range conflicts {
			b.WriteString("- [" + string(op.Type) + "] " + op.TargetPath + " (action: " + op.Conflict.Action + ")\n")
		}
	}

	// Warnings section
	var warnings []struct {
		path    string
		message string
	}
	for _, op := range input.Operations {
		for _, w := range op.Warnings {
			warnings = append(warnings, struct {
				path    string
				message string
			}{op.TargetPath, w.Message})
		}
	}
	if len(warnings) > 0 {
		b.WriteString("\n## Warnings\n")
		for _, w := range warnings {
			b.WriteString("- " + w.path + ": " + w.message + "\n")
		}
	}

	return b.String()
}

func joinPlatforms(platforms []platform.Platform) string {
	parts := make([]string, 0, len(platforms))
	for _, p := range platforms {
		parts = append(parts, string(p))
	}
	return strings.Join(parts, ", ")
}
