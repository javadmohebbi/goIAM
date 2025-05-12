// Package cmds contains all CLI subcommands for interacting with the goIAM API.
package cmds

import (
	"github.com/spf13/cobra"
)

// RegisterCommands adds all available subcommands to the root Cobra command.
//
// Parameters:
//   - root: pointer to the root Cobra command (`goiam`).
//   - apiURL: pointer to the base URL of the API.
//   - token: pointer to the JWT token for authenticated routes.
func RegisterCommands(root *cobra.Command, apiURL *string, token *string) {
	root.AddCommand(RegisterCmd(apiURL, token))         // User registration command
	root.AddCommand(LoginCmd(apiURL))                   // User login command
	root.AddCommand(Setup2FACmd(apiURL, token))         // 2FA setup command
	root.AddCommand(Verify2FACmd(apiURL, token))        // 2FA verification command
	root.AddCommand(Disable2FACmd(apiURL, token))       // Disable 2FA command
	root.AddCommand(RegenBackupCodesCmd(apiURL, token)) // Regenerate backup codes
	root.AddCommand(UpdateProfileCmd(apiURL, token))    // Update user profile
	root.AddCommand(UserAddCmd(apiURL, token))          // Add user
}
