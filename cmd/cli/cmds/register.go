// Package cmds provides CLI commands to interact with the goIAM backend.
package cmds

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// RegisterCmd creates the Cobra CLI command for registering a new user.
//
// This command collects user input via flags and prompts, then sends a POST request to
// the goIAM backend to create a new user. It supports secure password entry and piping.
//
// Required flags:
//   - --username or -u
//   - --email or -e
//
// Optional flags:
//   - --password or -p (if omitted, will prompt securely)
//   - --phone
//   - --first
//   - --middle
//   - --last
//   - --address
//   - --organization-name
//
// Example:
//
//	goiam register -u alice -e alice@example.com --organization-name "Acme Corp"
func RegisterCmd(apiURL *string, token *string) *cobra.Command {
	var username, password, email, phone, first, middle, last, address string
	var orgName string

	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new user",
		Run: func(cmd *cobra.Command, args []string) {

			// If password is not supplied via --password, read from stdin or terminal securely.
			if password == "" {
				fi, _ := os.Stdin.Stat()

				// If stdin is piped
				if (fi.Mode() & os.ModeCharDevice) == 0 {
					reader := bufio.NewReader(os.Stdin)
					passBytes, err := reader.ReadBytes('\n')
					if err != nil && err != io.EOF {
						fmt.Println("Failed to read piped password:", err)
						return
					}
					password = strings.TrimSpace(string(passBytes))
				} else {
					// If terminal: prompt securely and confirm password
					for {
						fmt.Print("Enter password: ")
						passBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
						fmt.Println()
						if err != nil {
							fmt.Println("Failed to read password:", err)
							return
						}
						if string(passBytes) == "" {
							fmt.Println("Empty password is not allowed")
							fmt.Println()
							continue
						}

						fmt.Print("Enter password confirmation: ")
						confPassBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
						fmt.Println()
						if err != nil {
							fmt.Println("Failed to read password confirmation:", err)
							return
						}
						if string(confPassBytes) != string(passBytes) {
							fmt.Println("Password and confirmation are not matched")
							fmt.Println()
							continue
						}

						password = string(passBytes)
						break
					}
				}
			}

			data := map[string]any{
				"username":     username,
				"password":     password,
				"email":        email,
				"phone_number": phone,
				"first_name":   first,
				"middle_name":  middle,
				"last_name":    last,
				"address":      address,
			}

			if orgName != "" {
				data["organization_name"] = orgName
			}

			res, _ := request(http.MethodPost, apiURL, "/auth/register", data, "")
			result, _ := io.ReadAll(res.Body)
			fmt.Println("Result:", string(result))
		},
	}

	// Define CLI flags
	cmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password (optional; will be read securely or from stdin)")
	cmd.Flags().StringVarP(&email, "email", "e", "", "Email")
	cmd.Flags().StringVar(&phone, "phone", "", "Phone")
	cmd.Flags().StringVar(&first, "first", "", "First name")
	cmd.Flags().StringVar(&middle, "middle", "", "Middle name")
	cmd.Flags().StringVar(&last, "last", "", "Last name")
	cmd.Flags().StringVar(&address, "address", "", "Address")
	cmd.Flags().StringVar(&orgName, "organization-name", "", "Name of new organization (optional)")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("email")

	return cmd
}
