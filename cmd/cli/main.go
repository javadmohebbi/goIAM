package main

import (
	"fmt"
	"os"

	"github.com/javadmohebbi/goIAM/cmd/cli/cmds"
	"github.com/spf13/cobra"
)

var (
	apiURL string
	token  string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "goiam",
		Short: "goIAM CLI - Manage IAM users and 2FA",
	}

	rootCmd.PersistentFlags().StringVar(&apiURL, "api", "http://localhost:8080", "Base API URL")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "JWT token for authenticated routes")

	cmds.RegisterCommands(rootCmd, &apiURL, &token)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
