package cli

import (
	"github.com/spf13/cobra"
)

// newBaseRootCmd creates the cobra root command with all subcommands registered.
func newBaseRootCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "vkc", Short: "Vibe and Kalika Code installer"}
	cmd.AddCommand(newInitCmd(), newDetectCmd(), newInstallCmd(), newValidateCmd(), newDoctorCmd(), newUpdateCmd(), newUninstallCmd())
	return cmd
}

// NewRootCmd builds the root command with REPL mode.
// Update checks are performed asynchronously inside the REPL model so that
// the CLI responds immediately on startup.
func NewRootCmd() *cobra.Command {
	cmd := newBaseRootCmd()

	cmd.Run = func(cmd *cobra.Command, _ []string) {
		runREPL(newBaseRootCmd, cmd.OutOrStdout())
	}

	return cmd
}
