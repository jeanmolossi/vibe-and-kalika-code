package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/ui"
)

func newInstallCmd() *cobra.Command {
	var yes, dryRun bool
	var targetsCSV, conflictAction string
	cmd := &cobra.Command{
		Use:   "install [source]",
		Short: "Install a package from a local dir or git URL",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, _ := os.Getwd()
			src := ""
			if len(args) > 0 {
				src = args[0]
			}
			if src == "" {
				sourceInput, err := ui.AskSource(yes, "")
				if err != nil {
					return exitError(app.ExitUserCancelled, err)
				}
				src = sourceInput.Source
			}
			out := cmd.OutOrStdout()
			res, code, err := app.Install(app.InstallOptions{
				Source:         src,
				ProjectRoot:    projectRoot,
				Targets:        parseTargets(targetsCSV),
				Yes:            yes,
				DryRun:         dryRun,
				ConflictAction: conflictAction,
			})
			if code == app.ExitConflictUnresolved && res != nil {
				var conflictingPaths []string
				for _, op := range res.Plan.Operations {
					if op.Conflict != nil {
						conflictingPaths = append(conflictingPaths, op.TargetPath)
					}
				}
				fmt.Fprintf(out, "\n  %d conflict(s) detected:\n", len(conflictingPaths))
				for _, p := range conflictingPaths {
					fmt.Fprintf(out, "   - %s\n", p)
				}
				var chosenAction string
				if yes {
					chosenAction = "backup-and-overwrite"
				} else {
					var askErr error
					chosenAction, askErr = ui.AskConflictAction(
						strings.Join(conflictingPaths, ", "),
						false,
					)
					if askErr != nil {
						return exitError(app.ExitUserCancelled, askErr)
					}
				}
				res, code, err = app.Install(app.InstallOptions{
					Source:         src,
					ProjectRoot:    projectRoot,
					Targets:        parseTargets(targetsCSV),
					Yes:            yes,
					DryRun:         dryRun,
					ConflictAction: chosenAction,
				})
			}
			if err != nil {
				return exitError(code, err)
			}
			fmt.Fprintf(out, "✓ Installed %s %s\n", res.Manifest.Name, res.Manifest.Version)
			printInstallSummary(out, res)
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Skip prompts")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show plan without applying")
	cmd.Flags().StringVar(&targetsCSV, "targets", "", "Comma-separated targets")
	cmd.Flags().StringVar(&conflictAction, "conflict-action", "", "skip|overwrite|backup-and-overwrite")
	return cmd
}

func printInstallSummary(out interface{ Write([]byte) (int, error) }, res *app.InstallResult) {
	created, modified, skipped := countOps(res)
	fmt.Fprintf(out, "  created:  %d files\n", created)
	fmt.Fprintf(out, "  modified: %d files\n", modified)
	fmt.Fprintf(out, "  skipped:  %d files\n", skipped)
	if res.ReportPath != "" {
		fmt.Fprintf(out, "  report:   %s\n", res.ReportPath)
	}
}

func countOps(res *app.InstallResult) (created, modified, skipped int) {
	for _, op := range res.Plan.Operations {
		switch op.Type {
		case platform.OperationCreate:
			created++
		case platform.OperationModify:
			modified++
		case platform.OperationSkip:
			skipped++
		}
	}
	return created, modified, skipped
}

func parseTargets(csv string) []platform.Platform {
	if csv == "" {
		return nil
	}
	parts := strings.Split(csv, ",")
	out := make([]platform.Platform, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, platform.Platform(trimmed))
		}
	}
	return out
}
