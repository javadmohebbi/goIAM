package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Disable2FACmd(apiURL *string, token *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "2fa-disable",
		Short: "Disable 2FA",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🚫 2fa-disable command not implemented yet")
		},
	}
	return cmd
}
