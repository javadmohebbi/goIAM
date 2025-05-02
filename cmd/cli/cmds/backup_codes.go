package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

func RegenBackupCodesCmd(apiURL *string, token *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup-codes",
		Short: "Regenerate backup codes",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🧯 backup-codes command not implemented yet")
		},
	}
	return cmd
}
