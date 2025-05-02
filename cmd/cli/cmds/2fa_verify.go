package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Verify2FACmd(apiURL *string, token *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "2fa-verify",
		Short: "Verify 2FA code",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("✅ 2fa-verify command not implemented yet")
		},
	}
	return cmd
}
