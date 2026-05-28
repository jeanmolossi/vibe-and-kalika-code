package cli

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "vkc", Short: "Vibe and Kalika Code installer"}
	cmd.AddCommand(newInitCmd(), newDetectCmd(), newInstallCmd(), newValidateCmd(), newDoctorCmd(), newUpdateCmd())
	return cmd
}
