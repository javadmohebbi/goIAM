package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/javadmohebbi/goIAM/cmd/cli/cmds"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	apiURL string
	token  string
)

func main() {

	fd := int(os.Stdin.Fd())

	var oldState *term.State
	var err error

	if term.IsTerminal(fd) {
		oldState, err = term.GetState(fd)
		if err != nil {
			panic(err)
		}

		defer func() {
			_ = term.Restore(fd, oldState)
		}()

		// Handle interrupt signals (Ctrl+C)
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigChan
			fmt.Println("\nProgram interrupted.")
			_ = term.Restore(fd, oldState)
			os.Exit(0)
		}()
	}

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
