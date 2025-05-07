// Package cmds provides CLI commands to interact with the goIAM backend.
package cmds

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
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
//   - --organization-id
//   - --organization-name
//   - --organization-slug
//
// Example:
//
//	goiam register -u alice -e alice@example.com --organization-id 1
func RegisterCmd(apiURL *string) *cobra.Command {
	var username, password, email, phone, first, middle, last, address string
	var orgID string
	var orgName, orgSlug string

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

			var data map[string]any
			if orgID != "" {
				orgUint, err := strconv.ParseUint(orgID, 10, 32)
				if err != nil {
					fmt.Println("Invalid organization ID:", err)
					return
				}
				data = map[string]any{
					"organization_id": uint(orgUint),
				}
			} else {
				data = map[string]any{}
				if orgName != "" {
					data["organization_name"] = orgName
				}
				if orgSlug != "" {
					data["organization_slug"] = orgSlug
				}
			}

			data["username"] = username
			data["password"] = password
			data["email"] = email
			data["phone_number"] = phone
			data["first_name"] = first
			data["middle_name"] = middle
			data["last_name"] = last
			data["address"] = address

			res, _ := post(apiURL, "/auth/register", data, "")
			result, _ := io.ReadAll(res.Body)
			fmt.Println(string(result))
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
	cmd.Flags().StringVar(&orgID, "organization-id", "", "Organization ID (optional)")
	cmd.Flags().StringVar(&orgName, "organization-name", "", "Name of new organization (optional)")
	cmd.Flags().StringVar(&orgSlug, "organization-slug", "", "Slug of new organization (optional)")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("email")
	// cmd.MarkFlagRequired("organization-id")

	return cmd
}
