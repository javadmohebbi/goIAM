package cmds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

func Disable2FACmd(apiURL *string, token *string) *cobra.Command {
	var code string

	cmd := &cobra.Command{
		Use:   "2fa-disable",
		Short: "Disable 2FA using TOTP code",
		Run: func(cmd *cobra.Command, args []string) {
			data := map[string]string{"code": code}
			body, _ := json.Marshal(data)

			req, _ := http.NewRequest("POST", *apiURL+"/secure/auth/2fa/disable", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
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
	cmd.Flags().StringVar(&code, "code", "", "Your current 2FA TOTP code")
	cmd.MarkFlagRequired("code")
	return cmd
}
