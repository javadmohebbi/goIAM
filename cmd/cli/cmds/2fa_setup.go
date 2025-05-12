// Package cmds provides CLI commands to interact with the goIAM backend.
package cmds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/spf13/cobra"
)

// Setup2FACmd returns the `2fa-setup` Cobra command,
// which registers a TOTP secret with the server and returns a QR code URL.
//
// This command:
//   - Sends a POST request to /secure/auth/2fa/setup with the JWT token.
//   - Receives a TOTP secret and otpauth:// QR URL.
//   - Optionally displays the QR code in the terminal using `qrencode`, if installed.
//
// Requirements:
//   - A valid --token must be provided.
func Setup2FACmd(apiURL *string, token *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "2fa-setup",
		Short: "Setup 2FA and get QR code",
		Run: func(cmd *cobra.Command, args []string) {
			url := *apiURL + "/s/auth/2fa/setup"

			req, _ := http.NewRequest("POST", url, nil)
			req.Header.Set("Authorization", "Bearer "+*token)

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("Request failed:", err)
				return
			}
			defer res.Body.Close()

			var out map[string]string
			if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
				fmt.Println("Invalid JSON response")
				return
			}

			fmt.Println("Secret:", out["secret"])
			fmt.Println("QR URL:", out["otpauth_url"])

			// Try displaying QR code using `qrencode` if available
			qrCmd := exec.Command("qrencode", "-t", "UTF8", out["otpauth_url"])
			if output, err := qrCmd.Output(); err == nil {
				fmt.Println(string(output))
			} else {
				fmt.Println("qrencode not installed (QR code not shown)")
			}
		},
	}

	cmd.MarkFlagRequired("token")
	return cmd
}
