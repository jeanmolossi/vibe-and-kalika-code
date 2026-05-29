package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
)

// newBaseRootCmd creates the cobra root command with all subcommands registered.
func newBaseRootCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "vkc", Short: "Vibe and Kalika Code installer"}
	cmd.AddCommand(newInitCmd(), newDetectCmd(), newInstallCmd(), newValidateCmd(), newDoctorCmd(), newUpdateCmd(), newUninstallCmd())
	return cmd
}

// NewRootCmd builds the root command with REPL mode and startup version checks.
func NewRootCmd() *cobra.Command {
	cmd := newBaseRootCmd()

	cmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		// Skip update checks during shell completion to avoid slow network
		// calls that break autocomplete response times.
		if cmd.Name() == cobra.ShellCompRequestCmd || cmd.Name() == cobra.ShellCompNoDescRequestCmd {
			return nil
		}

		cwd, err := os.Getwd()
		if err != nil {
			cwd = ""
		}

		check := app.CheckUpdates(cwd)

		w := cmd.ErrOrStderr()
		if check.CLIUpdateAvailable {
			fmt.Fprintf(w, "[!] vkc update available: %s → %s (run /update --self or vkc update --self)\n",
				check.CLICurrentVersion, check.CLILatestVersion)
		}

		for _, pkg := range check.PackageUpdates {
			fmt.Fprintf(w, "[!] Package update available: %s %s → %s (run /update or vkc update)\n",
				pkg.Name, pkg.CurrentVersion, pkg.LatestVersion)
		}

		return nil
	}

	cmd.Run = func(cmd *cobra.Command, _ []string) {
		runREPL(newBaseRootCmd, cmd.OutOrStdout())
	}

	return cmd
}
