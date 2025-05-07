// Package cmds provides CLI commands to interact with the goIAM backend.
package cmds

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// LoginCmd returns the `login` Cobra command which logs in a user,
// handles optional 2FA flow, and prints the JWT token on success.
//
// Parameters:
//   - apiURL: Pointer to the base URL of the goIAM API.
//
// Returns:
//   - *cobra.Command: The Cobra command for user login.
func LoginCmd(apiURL *string) *cobra.Command {
	var username, password string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login and handle optional 2FA verification",
		Run: func(cmd *cobra.Command, args []string) {
			// If password is not provided with --password, read from stdin or securely via terminal
			if password == "" {
				fi, _ := os.Stdin.Stat()
				if (fi.Mode() & os.ModeCharDevice) == 0 {
					// Read from piped stdin (e.g., echo "pass" | ...)
					reader := bufio.NewReader(os.Stdin)
					passBytes, err := reader.ReadBytes('\n')
					if err != nil && err != io.EOF {
						fmt.Println("Failed to read piped password:", err)
						return
					}
					password = strings.TrimSpace(string(passBytes))
				} else {
					// Read from terminal with masking
					fmt.Print("Enter password: ")
					passBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
					fmt.Println()
					if err != nil {
						fmt.Println("Failed to read password:", err)
						return
					}
					password = strings.TrimSpace(string(passBytes))
				}
			}

			// Construct login payload
			payload := map[string]any{
				"username": username,
				"password": password,
			}
			// body, _ := json.Marshal(payload)

			// Perform login request
			// res, err := http.Post(*apiURL+"/auth/login", "application/json", bytes.NewBuffer(body))
			res, err := post(
				apiURL,
				"/auth/login",
				payload,
				"",
			)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer res.Body.Close()

			output, _ := io.ReadAll(res.Body)

			// If 2FA is required, the API responds with 202 and a temporary token
			if res.StatusCode == http.StatusAccepted {
				fmt.Println("2FA required.")
				unverifiedToken := extractToken(output)
				if unverifiedToken == "" {
					fmt.Println("Could not extract token from login response.")
					return
				}

				// Always read 2FA input from /dev/tty, not stdin, to support piped workflows
				tty, err := os.Open("/dev/tty")
				if err != nil {
					fmt.Println("Cannot open /dev/tty for 2FA input:", err)
					return
				}
				defer tty.Close()

				reader := bufio.NewReader(tty)
				for attempt := 1; attempt <= 3; attempt++ {
					fmt.Printf("Enter TOTP or backup code (attempt %d/3): ", attempt)
					codeInput, _ := reader.ReadString('\n')
					codeInput = strings.TrimSpace(codeInput)

					verifyBody := map[string]any{"code": codeInput}
					// vbody, _ := json.Marshal(verifyBody)

					// Send 2FA verification request
					// req, _ := http.NewRequest("POST", *apiURL+"/secure/auth/2fa/verify", bytes.NewBuffer(vbody))
					// req.Header.Set("Authorization", "Bearer "+unverifiedToken)
					// req.Header.Set("Content-Type", "application/json")

					// verifyRes, err := http.DefaultClient.Do(req)
					verifyRes, err := post(
						apiURL,
						"/secure/auth/2fa/verify",
						verifyBody,
						"",
						map[string]string{
							"Authorization": "Bearer " + unverifiedToken,
						},
					)
					if err != nil {
						fmt.Println("Verification request error:", err)
						return
					}
					defer verifyRes.Body.Close()

					vout, _ := io.ReadAll(verifyRes.Body)

					if verifyRes.StatusCode == http.StatusOK {
						fmt.Println("2FA verified.")
						token := extractToken(vout)
						if token != "" {
							fmt.Printf("\n{\"token\": \"%s\"}\n", token)
						} else {
							fmt.Println(string(vout))
						}
						return
					}

					fmt.Println("Invalid code.")
					if attempt == 3 {
						fmt.Println("Too many failed attempts.")
						return
					}
				}
			} else {
				// Non-2FA login success
				token := extractToken(output)
				if token != "" {
					fmt.Println("Token:", token)
				} else {
					fmt.Println(string(output))
				}
			}
		},
	}

	// Command-line flags
	cmd.Flags().StringVarP(&username, "username", "u", "", "Username (required)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password (optional; will be read securely or from stdin)")
	cmd.MarkFlagRequired("username")

	return cmd
}

// extractToken parses the response body to extract a "token" field from JSON.
//
// Parameters:
//   - raw: Raw byte slice of the response body.
//
// Returns:
//   - string: The extracted token value if present, otherwise empty.
func extractToken(raw []byte) string {
	var result map[string]string
	if err := json.Unmarshal(raw, &result); err == nil {
		if token, ok := result["token"]; ok {
			return token
		}
	}
	return ""
}
