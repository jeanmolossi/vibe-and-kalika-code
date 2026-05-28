package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
)

var (
	detectOK     = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	detectMissed = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	detectLabel  = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)

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

func newDetectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "detect",
		Short: "Detect supported AI coding platforms",
		RunE: func(cmd *cobra.Command, _ []string) error {
			projectRoot, _ := os.Getwd()
			results := app.Detect(projectRoot)
			out := cmd.OutOrStdout()

			fmt.Fprintln(out, detectLabel.Render("\nDetected platforms:\n"))
			for _, res := range results {
				name := platformDisplayName(res.Platform)
				if res.Detected {
					fmt.Fprintf(out, "  %s %s\n", detectOK.Render("[✓]"), name)
				} else {
					fmt.Fprintf(out, "  %s %s\n", detectMissed.Render("[ ]"), name)
				}
				if res.BasePath != "" {
					fmt.Fprintf(out, "      home:   %s\n", res.BasePath)
				}
				if res.AgentsPath != "" {
					fmt.Fprintf(out, "      agents: %s\n", res.AgentsPath)
				}
				if res.SkillsPath != "" {
					fmt.Fprintf(out, "      skills: %s\n", res.SkillsPath)
				}
				for _, note := range res.Notes {
					fmt.Fprintf(out, "      note:   %s\n", note)
				}
				fmt.Fprintln(out)
			}
			return nil
		},
	}
}
