package cmds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

func LoginCmd(apiURL *string) *cobra.Command {
	var username, password string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login and get JWT",
		Run: func(cmd *cobra.Command, args []string) {
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
			fmt.Println(string(output))
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username (required)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password (required)")
	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")

	return cmd
}
