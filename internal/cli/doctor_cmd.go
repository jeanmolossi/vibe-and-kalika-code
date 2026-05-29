package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
)

const (
	envVarNotSet  = "(not set)"
	labelWritable = "writable"
	labelNotWrite = "not writable"
)

var (
	doctorOK   = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	doctorFail = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	doctorInfo = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
)

func doctorCheck(ok bool) string {
	if ok {
		return doctorOK.Render("[✓]")
	}
	return doctorFail.Render("[✗]")
}

func writableLabel(ok bool) string {
	if ok {
		return labelWritable
	}
	return labelNotWrite
}

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Run environment health checks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			projectRoot, _ := os.Getwd()
			res := app.Doctor(projectRoot)
			out := cmd.OutOrStdout()

			fmt.Fprintf(out, "\nProject root: %s\n", res.ProjectRoot)

			fmt.Fprintln(out, doctorInfo.Render("\nEnvironment variables:"))
			for _, key := range []string{"COPILOT_HOME", "CODEX_HOME"} {
				val := res.EnvVars[key]
				display := val
				if display == "" {
					display = envVarNotSet
				}
				fmt.Fprintf(out, "  %-14s %s\n", key, doctorInfo.Render(display))
			}

			fmt.Fprintln(out, doctorInfo.Render("\nEnvironment health:\n"))
			fmt.Fprintf(out, "  %s git available\n", doctorCheck(res.GitAvailable))
			if !res.GitAvailable {
				fmt.Fprintln(out, doctorInfo.Render("      install git to use Git-based package sources"))
			}

			fmt.Fprintf(out, "  %s installation state readable\n", doctorCheck(res.StateValid))
			if !res.StateValid {
				fmt.Fprintln(out, doctorInfo.Render("      ~/.ai-setup/installed.yaml is missing or corrupt — run vkc install first"))
			}

			fmt.Fprintln(out, doctorInfo.Render("\nPlatform status:"))
			for _, p := range res.Platforms {
				fmt.Fprintf(out, "  %s %s\n", doctorCheck(p.Detected), p.Name)
				if p.BasePath != "" {
					fmt.Fprintf(out, "      home:    %s\n", p.BasePath)
				}
				if p.AgentsDir != "" {
					fmt.Fprintf(out, "      agents:  %s\n", p.AgentsDir)
				}
				if p.SkillsDir != "" {
					fmt.Fprintf(out, "      skills:  %s\n", p.SkillsDir)
				}
				fmt.Fprintf(out, "      write:   %s %s\n", doctorCheck(p.Writable), writableLabel(p.Writable))
			}
			fmt.Fprintln(out)
			return nil
		},
	}
}
