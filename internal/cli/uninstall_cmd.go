package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
)

func newUninstallCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:               "uninstall <package>",
		Short:             "Uninstall a previously installed package",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeInstalledPackages,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve project root: %w", err)
			}
			res, code, err := app.Uninstall(app.UninstallOptions{
				Package:     args[0],
				ProjectRoot: projectRoot,
				DryRun:      dryRun,
			})
			if err != nil {
				return exitError(code, err)
			}
			out := cmd.OutOrStdout()
			if dryRun {
				fmt.Fprintf(out, "Would uninstall %s (dry-run)\n", res.Package)
			} else {
				fmt.Fprintf(out, "Uninstalled %s\n", res.Package)
			}
			fmt.Fprintf(out, "  removed:  %d files\n", len(res.FilesRemoved))
			fmt.Fprintf(out, "  restored: %d files\n", len(res.FilesRestored))
			fmt.Fprintf(out, "  markers:  %d removed\n", len(res.MarkersRemoved))
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be uninstalled without applying")
	return cmd
}

// completeInstalledPackages provides shell completion with the list of
// currently installed package names.
func completeInstalledPackages(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	store, err := state.Read()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	names := make([]string, 0, len(store.Installations))
	for _, inst := range store.Installations {
		names = append(names, inst.Package)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
