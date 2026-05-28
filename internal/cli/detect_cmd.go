package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
)

func newDetectCmd() *cobra.Command {
	return &cobra.Command{Use: "detect", Short: "Detect supported platforms", RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot, _ := os.Getwd()
		for _, res := range app.Detect(projectRoot) {
			fmt.Fprintf(cmd.OutOrStdout(), "%s detected=%t base=%s\n", res.Platform, res.Detected, res.BasePath)
		}
		return nil
	}}
}
