package cmds

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// LoginCmd returns the login command, which handles authentication and 2FA.
func LoginCmd(apiURL *string) *cobra.Command {
	var username, password string
	var useStdin bool

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login and handle optional 2FA verification",
		Run: func(cmd *cobra.Command, args []string) {
			// Determine password source
			if password == "" {
				fi, _ := os.Stdin.Stat()
				if (fi.Mode() & os.ModeCharDevice) == 0 {
					// stdin is piped
					passBytes, err := io.ReadAll(os.Stdin)
					if err != nil {
						fmt.Println("Failed to read piped password:", err)
						return
					}
					password = strings.TrimSpace(string(passBytes))
				} else {
					// stdin is terminal (secure prompt)
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

			payload := map[string]string{
				"username": username,
				"password": password,
			}
			body, _ := json.Marshal(payload)

			res, err := http.Post(*apiURL+"/auth/login", "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer res.Body.Close()

			output, _ := io.ReadAll(res.Body)

			// Case: 2FA required
			if res.StatusCode == http.StatusAccepted {
				fmt.Println("2FA required.")

				unverifiedToken := extractToken(output)
				if unverifiedToken == "" {
					fmt.Println("Could not extract token from login response.")
					return
				}

				reader := bufio.NewReader(os.Stdin)
				for attempt := 1; attempt <= 3; attempt++ {
					fmt.Printf("Enter TOTP or backup code (attempt %d/3): ", attempt)
					codeInput, _ := reader.ReadString('\n')
					codeInput = strings.TrimSpace(codeInput)

					verifyBody := map[string]string{"code": codeInput}
					vbody, _ := json.Marshal(verifyBody)

					req, _ := http.NewRequest("POST", *apiURL+"/secure/auth/2fa/verify", bytes.NewBuffer(vbody))
					req.Header.Set("Authorization", "Bearer "+unverifiedToken)
					req.Header.Set("Content-Type", "application/json")

					verifyRes, err := http.DefaultClient.Do(req)
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
							fmt.Println("Token:", token)
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
				// Normal login
				token := extractToken(output)
				if token != "" {
					fmt.Println("Token:", token)
				} else {
					fmt.Println(string(output))
				}
			}
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username (required)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password (optional; will be read from stdin if not provided)")
	cmd.Flags().BoolVar(&useStdin, "stdin", false, "Deprecated: use --password='' instead")
	cmd.MarkFlagRequired("username")

	return cmd
}

// extractToken extracts the "token" value from a JSON response
func extractToken(raw []byte) string {
	var result map[string]string
	if err := json.Unmarshal(raw, &result); err == nil {
		if token, ok := result["token"]; ok {
			return token
		}
	}
	return ""
}
