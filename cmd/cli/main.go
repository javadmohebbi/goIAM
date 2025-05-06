// Package main provides the entry point for the goIAM CLI application.
// It uses Cobra for command parsing and supports optional terminal state restoration and signal handling.
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
	apiURL string // Base URL of the goIAM API
	token  string // JWT token for authenticated routes
)

// main initializes terminal handling (for clean stdin restoration on exit),
// sets up root Cobra command, and registers all subcommands.
func main() {
	fd := int(os.Stdin.Fd())

	var oldState *term.State
	var err error

	// Save the terminal state and restore it on exit if running in terminal
	if term.IsTerminal(fd) {
		oldState, err = term.GetState(fd)
		if err != nil {
			panic(err)
		}

		defer func() {
			_ = term.Restore(fd, oldState)
		}()

		// Gracefully handle interrupt (Ctrl+C)
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigChan
			fmt.Println("\nProgram interrupted.")
			_ = term.Restore(fd, oldState)
			os.Exit(0)
		}()
	}

	// Define root command for goIAM CLI
	rootCmd := &cobra.Command{
		Use:   "goiam",
		Short: "goIAM CLI - Manage IAM users and 2FA",
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&apiURL, "api", "http://localhost:8080", "Base API URL")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "JWT token for authenticated routes")

	// Register subcommands
	cmds.RegisterCommands(rootCmd, &apiURL, &token)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
