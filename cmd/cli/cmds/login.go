package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

func LoginCmd(apiURL *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login and get JWT",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🔐 login command not implemented yet")
		},
	}
	return cmd
}
