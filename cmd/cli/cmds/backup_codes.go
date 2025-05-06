// Package cmds provides CLI commands to interact with the goIAM backend.
package cmds

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// RegenBackupCodesCmd returns the `backup-codes` Cobra command,
// which regenerates a new set of backup codes for the authenticated user.
//
// Requires:
//   - A valid JWT token (passed via --token flag).
//   - The goIAM server running and reachable at the given API URL.
//
// Sends:
//   - POST /secure/auth/backup-codes/regenerate
//
// Prints:
//   - A list of new backup codes or an error response from the server.
func RegenBackupCodesCmd(apiURL *string, token *string) *cobra.Command {
	return &cobra.Command{
		Use:   "backup-codes",
		Short: "Regenerate backup codes",
		Run: func(cmd *cobra.Command, args []string) {
			// Construct POST request
			req, _ := http.NewRequest("POST", *apiURL+"/secure/auth/backup-codes/regenerate", nil)
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
}
