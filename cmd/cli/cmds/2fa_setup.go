package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Setup2FACmd(apiURL *string, token *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "2fa-setup",
		Short: "Setup 2FA and get QR code",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🔐 2fa-setup command not implemented yet")
		},
	}
	return cmd
}
