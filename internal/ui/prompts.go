package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/source"
)

const sourceTypeGit = "git"

// InstallSummary holds data for ShowFinalSummary, avoiding an import of the app package.
type InstallSummary struct {
	FilesCreated  []string
	FilesModified []string
	FilesSkipped  []string
	ReportPath    string
	BackupPath    string
}

// SourceInput holds the result of source selection.
type SourceInput struct {
	SourceType string // "local" or "git"
	Source     string
}

// Confirm asks for a yes/no confirmation.
func Confirm(message string, assumeYes bool) (bool, error) {
	if assumeYes {
		return true, nil
	}
	if os.Getenv("CI") != "" {
		return false, nil
	}
	var confirmed bool
	form := huh.NewForm(huh.NewGroup(huh.NewConfirm().Title(message).Value(&confirmed)))
	if err := form.Run(); err != nil {
		return false, err
	}
	return confirmed, nil
}

// AskSource asks the user to choose a local directory or git URL as the package source.
// If defaultSrc is non-empty the prompt is skipped and that value is returned directly.
func AskSource(assumeYes bool, defaultSrc string) (SourceInput, error) {
	nonInteractive := assumeYes || os.Getenv("CI") != ""

	// If a source was already provided (positional arg), use it directly.
	if defaultSrc != "" {
		return SourceInput{
			SourceType: sourceKind(defaultSrc),
			Source:     defaultSrc,
		}, nil
	}

	if nonInteractive {
		return SourceInput{
			SourceType: sourceTypeGit,
			Source:     source.DefaultSource,
		}, nil
	}

	// Step 1: pick type.
	var sourceType string
	form1 := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Package source type").
			Options(
				huh.NewOption("Local directory", "local"),
				huh.NewOption("Git repository URL", sourceTypeGit),
			).
			Value(&sourceType),
	))
	if err := form1.Run(); err != nil {
		return SourceInput{}, err
	}

	// Step 2: enter path/URL.
	placeholder := "./path/to/package"
	if sourceType == sourceTypeGit {
		placeholder = "https://github.com/user/repo"
	}
	var sourcePath string
	pathInput := huh.NewInput().
		Title("Source path or URL").
		Placeholder(placeholder).
		Value(&sourcePath)

	if sourceType != sourceTypeGit {
		pathInput = pathInput.SuggestionsFunc(func() []string {
			return pathSuggestions(sourcePath)
		}, &sourcePath)
	}

	form2 := huh.NewForm(huh.NewGroup(pathInput))
	if err := form2.Run(); err != nil {
		return SourceInput{}, err
	}

	sourcePath = strings.TrimSpace(sourcePath)
	if sourcePath == "" {
		return SourceInput{}, fmt.Errorf("source path or URL is required")
	}

	return SourceInput{SourceType: sourceType, Source: sourcePath}, nil
}

// AskPlatforms shows a multi-select checkbox for platform selection.
// When running non-interactively it returns all detected platforms (or all platforms if none detected).
func AskPlatforms(detections []platform.DetectionResult, assumeYes bool) ([]platform.Platform, error) {
	nonInteractive := assumeYes || os.Getenv("CI") != ""

	if nonInteractive {
		var selected []platform.Platform
		for _, d := range detections {
			if d.Detected {
				selected = append(selected, d.Platform)
			}
		}
		if len(selected) == 0 {
			for _, d := range detections {
				selected = append(selected, d.Platform)
			}
		}
		return selected, nil
	}

	options := make([]huh.Option[platform.Platform], 0, len(detections))
	for _, d := range detections {
		options = append(options, huh.NewOption(platformLabel(d), d.Platform))
	}

	var selected []platform.Platform
	form := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[platform.Platform]().
			Title("Select installation targets").
			Options(options...).
			Value(&selected),
	))
	if err := form.Run(); err != nil {
		return nil, err
	}

	return selected, nil
}

// AskConflictAction asks the user what to do when a file conflict is found.
func AskConflictAction(targetPath string, assumeYes bool) (string, error) {
	if assumeYes || os.Getenv("CI") != "" {
		return "backup-and-overwrite", nil
	}

	var action string
	form := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title(fmt.Sprintf("Conflict: %s", targetPath)).
			Description("File already exists. Choose an action:").
			Options(
				huh.NewOption("Backup existing file and overwrite", "backup-and-overwrite"),
				huh.NewOption("Overwrite without backup", "overwrite"),
				huh.NewOption("Skip this file", "skip"),
			).
			Value(&action),
	))
	if err := form.Run(); err != nil {
		return "", err
	}
	return action, nil
}

// ShowDryRun prints the planned operations to stdout.
func ShowDryRun(ops []platform.PlannedOperation) {
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	create := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	modify := lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
	skip := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	warn := lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	conflict := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	fmt.Println(header.Render("\n📋 Dry-run plan:"))

	for _, op := range ops {
		switch op.Type {
		case platform.OperationCreate:
			fmt.Printf("  %s  %s\n", create.Render("[create]"), op.TargetPath)
		case platform.OperationModify:
			fmt.Printf("  %s  %s\n", modify.Render("[modify]"), op.TargetPath)
		case platform.OperationSkip:
			fmt.Printf("  %s    %s\n", skip.Render("[skip]"), op.TargetPath)
		default:
			fmt.Printf("  [?]    %s\n", op.TargetPath)
		}
		if op.Conflict != nil {
			fmt.Printf("           %s existing: %s\n", conflict.Render("conflict"), op.Conflict.ExistingPath)
		}
		for _, w := range op.Warnings {
			fmt.Printf("           %s %s\n", warn.Render("warn:"), w.Message)
		}
	}
	fmt.Println()
}

// ShowPackageSummary prints the package name, version, description, agent and skill counts.
func ShowPackageSummary(m *manifest.Manifest) {
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	label := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	fmt.Println(header.Render("\n📦 Package summary:"))
	fmt.Printf("  %s %s\n", label.Render("Name:"), m.Name)
	fmt.Printf("  %s %s\n", label.Render("Version:"), m.Version)
	if m.Description != "" {
		fmt.Printf("  %s %s\n", label.Render("Description:"), m.Description)
	}
	if m.Author != "" {
		fmt.Printf("  %s %s\n", label.Render("Author:"), m.Author)
	}
	fmt.Printf("  %s %d\n", label.Render("Agents:"), len(m.Agents))
	fmt.Printf("  %s %d\n", label.Render("Skills:"), len(m.Skills))
	fmt.Println()
}

// ShowFinalSummary prints the installation result to stdout.
func ShowFinalSummary(s InstallSummary) {
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("46"))
	label := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	fmt.Println(header.Render("\n✅ Installation complete!"))
	fmt.Printf("  %s %d\n", label.Render("Files created:"), len(s.FilesCreated))
	fmt.Printf("  %s %d\n", label.Render("Files modified:"), len(s.FilesModified))
	fmt.Printf("  %s %d\n", label.Render("Files skipped:"), len(s.FilesSkipped))
	if s.ReportPath != "" {
		fmt.Printf("  %s %s\n", label.Render("Report:"), s.ReportPath)
	}
	if s.BackupPath != "" {
		fmt.Printf("  %s %s\n", label.Render("Backup:"), s.BackupPath)
	}
	fmt.Println()
}

// --- helpers ---

// sourceKind returns "git" for git URLs, "local" otherwise.
func sourceKind(src string) string {
	for _, prefix := range []string{"http://", "https://", "git@", "ssh://", "file://"} {
		if strings.HasPrefix(src, prefix) {
			return sourceTypeGit
		}
	}
	return "local"
}

// platformLabel builds the human-readable checkbox label for a detection result.
func platformLabel(d platform.DetectionResult) string {
	name := platformDisplayName(d.Platform)
	if d.Detected {
		path := d.BasePath
		if path == "" {
			path = d.AgentsPath
		}
		return fmt.Sprintf("%s  (detected at %s)", name, path)
	}
	return fmt.Sprintf("%s  (not detected, can still install project files)", name)
}

func platformDisplayName(p platform.Platform) string {
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

// pathSuggestions returns directory entries that match the given prefix for
// path autocompletion in the local source input.
//
// It preserves the exact prefix format typed by the user so that textinput's
// HasPrefix filter matches correctly (e.g. "./f" yields "./foo", not "foo").
func pathSuggestions(prefix string) []string {
	if prefix == "" {
		return nil
	}

	var dir, partial string
	i := strings.LastIndex(prefix, "/")
	if i < 0 {
		dir = "."
		partial = prefix
	} else {
		dir = prefix[:i+1]
		partial = prefix[i+1:]
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	suggestions := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if partial != "" && !strings.HasPrefix(name, partial) {
			continue
		}
		if dir == "." {
			suggestions = append(suggestions, name)
		} else {
			suggestions = append(suggestions, dir+name)
		}
	}
	return suggestions
}
