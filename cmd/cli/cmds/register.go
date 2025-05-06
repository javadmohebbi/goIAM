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

func RegisterCmd(apiURL *string) *cobra.Command {
	var username, password, email, phone, first, middle, last, address string
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new user",
		Run: func(cmd *cobra.Command, args []string) {

			// Detect if password was provided directly
			if password == "" {
				fi, _ := os.Stdin.Stat()
				if (fi.Mode() & os.ModeCharDevice) == 0 {
					// stdin is being piped (but we want to preserve for 2FA)
					reader := bufio.NewReader(os.Stdin)
					passBytes, err := reader.ReadBytes('\n')
					if err != nil && err != io.EOF {
						fmt.Println("Failed to read piped password:", err)
						return
					}
					password = strings.TrimSpace(string(passBytes))
				} else {

					for {

						// terminal: securely read password
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

						// check conf pass
						// terminal: securely read password
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

			data := map[string]string{
				"username": username, "password": password,
				"email": email, "phone_number": phone,
				"first_name": first, "middle_name": middle,
				"last_name": last, "address": address,
			}
			post(apiURL, "/auth/register", data, "")
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	// cmd.Flags().StringVar(&password, "password", "", "Password")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password (optional; will be read securely or from stdin)")
	cmd.Flags().StringVarP(&email, "email", "e", "", "Email")
	cmd.Flags().StringVar(&phone, "phone", "", "Phone")
	cmd.Flags().StringVar(&first, "first", "", "First name")
	cmd.Flags().StringVar(&middle, "middle", "", "Middle name")
	cmd.Flags().StringVar(&last, "last", "", "Last name")
	cmd.Flags().StringVar(&address, "address", "", "Address")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("email")
	// cmd.MarkFlagRequired("password")
	return cmd
}

func post(apiURL *string, path string, data map[string]string, token string) {
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", *apiURL+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer res.Body.Close()
	result, _ := io.ReadAll(res.Body)
	fmt.Println(string(result))
}
