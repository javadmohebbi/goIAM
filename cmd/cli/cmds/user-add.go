package cmds

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func UserAddCmd(apiURL *string, token *string) *cobra.Command {
	var username, email, fName, mName, lName, phone string

	// organization will be added automatically by the caller
	// and the token will be sent to the API
	// var orgSlug string

	cmd := &cobra.Command{
		Use:   "user-add",
		Short: "Create user for an organization. Only an autheticated user with a valid token and permission can create a user for an organization",
		Run: func(cmd *cobra.Command, args []string) {
			if username == "" {
				fmt.Println("Username (--username or -u) is required")
				os.Exit(1)
			}
			if email == "" {
				fmt.Println("Email address (--email or -e) is required")
				os.Exit(1)
			}

			data := map[string]any{
				"username":     username,
				"email":        email,
				"phone_number": phone,
				"first_name":   fName,
				"middle_name":  mName,
				"last_name":    lName,
			}

			res, err := request(http.MethodPost, apiURL, "/s/user/create", data, *token)
			if err != nil {
				fmt.Println("Request failed:", err)
				os.Exit(1)
			}
			defer res.Body.Close()

			output, _ := io.ReadAll(res.Body)
			if res.StatusCode != http.StatusCreated {
				fmt.Printf("Error: status %d\n%s\n", res.StatusCode, string(output))
				return
			}
			fmt.Println("User created successfully but inative. Activate the user setting a password")

		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	cmd.Flags().StringVarP(&email, "email", "e", "", "Email address")
	// cmd.Flags().StringVarP(&orgSlug, "organization-slug", "o", "", "Organization Slug")

	cmd.Flags().StringVarP(&fName, "first-name", "f", "", "First name")
	cmd.Flags().StringVarP(&mName, "middle-name", "m", "", "Middle name")
	cmd.Flags().StringVarP(&lName, "last-name", "l", "", "Last name")

	cmd.Flags().StringVarP(&phone, "phone", "p", "", "Phone number")

	return cmd
}
