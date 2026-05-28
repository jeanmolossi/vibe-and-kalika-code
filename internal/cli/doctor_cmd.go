package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{Use: "doctor", Short: "Run environment health checks", RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot, _ := os.Getwd()
		res := app.Doctor(projectRoot)
		fmt.Fprintf(cmd.OutOrStdout(), "git=%t state=%t\n", res.GitAvailable, res.StateValid)
		return nil
	}}
}
