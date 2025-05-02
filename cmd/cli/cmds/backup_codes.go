package cmds

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

func RegenBackupCodesCmd(apiURL *string, token *string) *cobra.Command {
	return &cobra.Command{
		Use:   "backup-codes",
		Short: "Regenerate backup codes",
		Run: func(cmd *cobra.Command, args []string) {
			req, _ := http.NewRequest("POST", *apiURL+"/secure/auth/backup-codes/regenerate", nil)
			req.Header.Set("Authorization", "Bearer "+*token)

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer res.Body.Close()

			output, _ := io.ReadAll(res.Body)
			fmt.Println(string(output))
		},
	}
}
