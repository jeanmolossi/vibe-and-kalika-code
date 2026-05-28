package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
)

// newUpdateCmd returns the cobra command for `vkc update`.
func newUpdateCmd() *cobra.Command {
	var self, dryRun, yes bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update installed packages or the vkc CLI itself",
		Long: `Update checks all installed packages against their sources and re-installs
any that have a newer version available.

With --self, it fetches the latest GitHub release and replaces the running binary.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cmd.OutOrStdout()

			if self {
				res, code, err := app.SelfUpdate(app.SelfUpdateOptions{
					DryRun: dryRun,
					Yes:    yes,
				})
				if err != nil {
					return exitError(code, err)
				}
				if res.AlreadyLatest {
					fmt.Fprintln(out, "vkc is already up to date.")
					return nil
				}
				if res.DryRun {
					fmt.Fprintf(out, "Would update vkc from %s to %s\n", res.CurrentVersion, res.LatestVersion)
					return nil
				}
				fmt.Fprintf(out, "✓ Updated vkc from %s to %s\n", res.CurrentVersion, res.LatestVersion)
				return nil
			}

			projectRoot, _ := os.Getwd()

			res, code, err := app.Update(app.UpdateOptions{
				ProjectRoot: projectRoot,
				DryRun:      dryRun,
				Yes:         yes,
			})
			if err != nil {
				return exitError(code, err)
			}

			if len(res.Updated) == 0 && len(res.Skipped) == 0 {
				fmt.Fprintln(out, "All packages are up to date.")
				return nil
			}

			for _, name := range res.Updated {
				fmt.Fprintf(out, "✓ Updated %s\n", name)
			}
			for _, name := range res.Skipped {
				fmt.Fprintf(out, "~ Skipped %s\n", name)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&self, "self", false, "Update the vkc CLI itself")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show plan without applying")
	cmd.Flags().BoolVar(&yes, "yes", false, "Skip prompts")

	return cmd
}
