package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
)

func newInitCmd() *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "init [source]",
		Short: "Run the interactive setup wizard",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			src := ""
			if len(args) > 0 {
				src = args[0]
			}
			projectRoot, _ := os.Getwd()
			_, code, err := app.Init(projectRoot, src, yes)
			if err != nil {
				return exitError(code, err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Initialization completed.")
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Skip prompts")
	return cmd
}
