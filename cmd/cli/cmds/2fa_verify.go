// Package cmds provides CLI commands to interact with the goIAM backend.
package cmds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// Verify2FACmd returns the `2fa-verify` Cobra command,
// which sends a TOTP code to the server to verify and complete 2FA setup or login.
//
// Requires:
//   - A valid JWT token (via --token).
//   - A TOTP code provided by the user from their authenticator app.
//
// Sends:
//   - POST /secure/auth/2fa/verify with {"code": "123456"} in JSON body.
//
// Prints:
//   - A success or error message from the server.
func Verify2FACmd(apiURL *string, token *string) *cobra.Command {
	var code string

	cmd := &cobra.Command{
		Use:   "2fa-verify",
		Short: "Verify 2FA code",
		Run: func(cmd *cobra.Command, args []string) {
			// Prepare payload
			data := map[string]string{"code": code}
			body, _ := json.Marshal(data)

			// Prepare request
			req, _ := http.NewRequest("POST", *apiURL+"/secure/auth/2fa/verify", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+*token)

			// Send request
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer res.Body.Close()

			// Read and print response
			output, _ := io.ReadAll(res.Body)
			fmt.Println(string(output))
		},
	}

	cmd.Flags().StringVar(&code, "code", "", "TOTP code from your authenticator app")
	cmd.MarkFlagRequired("code")
	return cmd
}
