package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
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

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Run environment health checks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			projectRoot, _ := os.Getwd()
			res := app.Doctor(projectRoot)
			out := cmd.OutOrStdout()

			fmt.Fprintln(out, doctorInfo.Render("\nEnvironment health:\n"))
			fmt.Fprintf(out, "  %s git available\n", doctorCheck(res.GitAvailable))
			if !res.GitAvailable {
				fmt.Fprintln(out, doctorInfo.Render("      install git to use Git-based package sources"))
			}

			fmt.Fprintf(out, "  %s installation state readable\n", doctorCheck(res.StateValid))
			if !res.StateValid {
				fmt.Fprintln(out, doctorInfo.Render("      .ai-setup/installed.yaml is missing or corrupt — run vkc install first"))
			}

			fmt.Fprintln(out, doctorInfo.Render("\nPlatform status:"))
			for _, p := range res.Platforms {
				status := doctorCheck(p.Detected)
				fmt.Fprintf(out, "  %s %s\n", status, p.Name)
				if p.BasePath != "" {
					fmt.Fprintf(out, "      home: %s\n", p.BasePath)
				}
			}
			fmt.Fprintln(out)
			return nil
		},
	}
}
