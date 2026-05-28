package cli

import (
	"fmt"
	"os"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	return &cobra.Command{Use: "validate <source>", Short: "Validate a package", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot, _ := os.Getwd()
		res, code, err := app.ValidateSource(projectRoot, args[0])
		if err != nil {
			return exitError(code, err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Valid package %s %s\n", res.Manifest.Name, res.Manifest.Version)
		return nil
	}}
}
