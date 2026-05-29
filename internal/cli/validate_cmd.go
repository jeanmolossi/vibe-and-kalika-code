package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/ui"
)

func newValidateCmd() *cobra.Command {
	return &cobra.Command{Use: "validate [source]", Short: "Validate a package", Args: cobra.MaximumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		src := ""
		if len(args) > 0 {
			src = args[0]
		}
		if src == "" {
			sourceInput, err := ui.AskSource(false, "")
			if err != nil {
				return exitError(app.ExitUserCancelled, err)
			}
			src = sourceInput.Source
		}
		projectRoot, _ := os.Getwd()
		res, code, err := app.ValidateSource(projectRoot, src)
		if err != nil {
			return exitError(code, err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Valid package %s %s\n", res.Manifest.Name, res.Manifest.Version)
		return nil
	}}
}
