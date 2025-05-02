package cmds

import (
    "github.com/spf13/cobra"
)

func RegisterCommands(root *cobra.Command, apiURL *string, token *string) {
    root.AddCommand(RegisterCmd(apiURL))
    root.AddCommand(LoginCmd(apiURL))
    root.AddCommand(Setup2FACmd(apiURL, token))
    root.AddCommand(Verify2FACmd(apiURL, token))
    root.AddCommand(Disable2FACmd(apiURL, token))
    root.AddCommand(RegenBackupCodesCmd(apiURL, token))
}
