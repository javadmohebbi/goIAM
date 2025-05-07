package cmds

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// UpdateProfileCmd creates a Cobra CLI command that sends a PATCH request to update the user's profile.
//
// Flags:
//
//	--first-name   First name to update
//	--last-name    Last name to update
//	--middle-name  Middle name to update
//	--address      Address to update
//
// Only the provided flags are included in the request body.
// Requires a valid bearer token and authenticated session.
func UpdateProfileCmd(apiURL *string, token *string) *cobra.Command {
	var firstName, lastName, middleName, address string

	cmd := &cobra.Command{
		Use:   "update-profile",
		Short: "Update your user profile",
		Run: func(cmd *cobra.Command, args []string) {
			// Prepare the update payload
			payload := make(map[string]any)

			if firstName != "" {
				payload["first_name"] = firstName
			}
			if lastName != "" {
				payload["last_name"] = lastName
			}
			if middleName != "" {
				payload["middle_name"] = middleName
			}
			if address != "" {
				payload["address"] = address
			}

			if len(payload) == 0 {
				fmt.Println("No fields to update. Use flags to specify profile fields.")
				return
			}

			if token == nil || *token == "" {
				fmt.Println("Error: --token is required to update your profile.")
				return
			}

			res, err := request(
				http.MethodPatch,
				apiURL,
				"/secure/auth/profile",
				payload,
				*token,
			)
			if err != nil {
				fmt.Println("Request failed:", err)
				return
			}
			defer res.Body.Close()
			output, _ := io.ReadAll(res.Body)

			if res.StatusCode != http.StatusAccepted {
				fmt.Printf("Error: status %d\n%s\n", res.StatusCode, string(output))
				return
			}
			fmt.Println("Profile updated successfully.")
		},
	}

	cmd.Flags().StringVar(&firstName, "first-name", "", "First name")
	cmd.Flags().StringVar(&lastName, "last-name", "", "Last name")
	cmd.Flags().StringVar(&middleName, "middle-name", "", "Middle name")
	cmd.Flags().StringVar(&address, "address", "", "Address")

	return cmd
}
