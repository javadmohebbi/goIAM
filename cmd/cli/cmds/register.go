package cmds

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"

    "github.com/spf13/cobra"
)

func RegisterCmd(apiURL *string) *cobra.Command {
    var username, password, email, phone, first, middle, last, address string
    cmd := &cobra.Command{
        Use:   "register",
        Short: "Register a new user",
        Run: func(cmd *cobra.Command, args []string) {
            data := map[string]string{
                "username": username, "password": password,
                "email": email, "phone_number": phone,
                "first_name": first, "middle_name": middle,
                "last_name": last, "address": address,
            }
            post(apiURL, "/auth/register", data, "")
        },
    }

    cmd.Flags().StringVar(&username, "username", "", "Username")
    cmd.Flags().StringVar(&password, "password", "", "Password")
    cmd.Flags().StringVar(&email, "email", "", "Email")
    cmd.Flags().StringVar(&phone, "phone", "", "Phone")
    cmd.Flags().StringVar(&first, "first", "", "First name")
    cmd.Flags().StringVar(&middle, "middle", "", "Middle name")
    cmd.Flags().StringVar(&last, "last", "", "Last name")
    cmd.Flags().StringVar(&address, "address", "", "Address")

    cmd.MarkFlagRequired("username")
    cmd.MarkFlagRequired("password")
    return cmd
}

func post(apiURL *string, path string, data map[string]string, token string) {
    body, _ := json.Marshal(data)
    req, _ := http.NewRequest("POST", *apiURL+path, bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    if token != "" {
        req.Header.Set("Authorization", "Bearer " + token)
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
