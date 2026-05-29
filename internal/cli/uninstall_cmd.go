package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/ui"
)

func newUninstallCmd() *cobra.Command {
	var dryRun, force bool
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
				Force:       force,
			})
			if err != nil {
				return exitError(code, err)
			}

			out := cmd.OutOrStdout()

			// Interactively ask about items outside the project root that were skipped.
			if !dryRun && !force && len(res.FilesSkipped) > 0 {
				fmt.Fprintf(out, "\n⚠  %d item(s) fora do diretório do projeto foram ignorados:\n", len(res.FilesSkipped))
				for _, item := range res.FilesSkipped {
					if item.AgentName != "" {
						fmt.Fprintf(out, "   - %s (bloco de agente: %s)\n", item.Path, item.AgentName)
					} else {
						fmt.Fprintf(out, "   - %s\n", item.Path)
					}
				}
				confirmed, cerr := ui.Confirm("Deseja remover esses arquivos também?", false)
				if cerr != nil {
					return fmt.Errorf("confirmação: %w", cerr)
				}
				if confirmed {
					extraFiles, extraMarkers, rerr := app.RemoveSkippedItems(res.FilesSkipped)
					if rerr != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "aviso: remoção parcial: %v\n", rerr)
					}
					res.FilesRemoved = append(res.FilesRemoved, extraFiles...)
					res.MarkersRemoved = append(res.MarkersRemoved, extraMarkers...)
				}
			}

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
	cmd.Flags().BoolVar(&force, "force", false, "Remove files outside the project root without prompting")
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
