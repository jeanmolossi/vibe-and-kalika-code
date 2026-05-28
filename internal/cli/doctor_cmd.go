package cli

import (
	"fmt"
	"os"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{Use: "doctor", Short: "Run environment health checks", RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot, _ := os.Getwd()
		res := app.Doctor(projectRoot)
		fmt.Fprintf(cmd.OutOrStdout(), "git=%t state=%t\n", res.GitAvailable, res.StateValid)
		return nil
	}}
}
