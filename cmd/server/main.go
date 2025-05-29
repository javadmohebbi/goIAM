// Package main starts the goIAM HTTP API service.
package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/javadmohebbi/goIAM/internal/api"
	"github.com/javadmohebbi/goIAM/internal/config"
	"github.com/javadmohebbi/goIAM/internal/db"
	flag "github.com/spf13/pflag"
)

// main is the entry point for launching the goIAM API server.
//
// It performs the following tasks:
//  1. Parses command-line flags: config path, HTTP port, and debug mode.
//  2. Loads configuration from a YAML file (with optional override via IAM_CONFIG_PATH).
//  3. Applies environment variable overrides for configuration values.
//  4. Applies runtime overrides from CLI flags (port, debug).
//  5. Initializes the database connection using GORM.
//  6. Starts the HTTP server with the loaded configuration.
//
// Environment Variables:
//   - IAM_CONFIG_PATH: override config file location
//   - IAM_PORT: override server port
//   - IAM_DATABASE: override database engine
//   - IAM_DATABASE_DSN: override connection string
//   - IAM_AUTH_PROVIDER: override authentication providers (comma-separated)
//
// Flags:
//
//	-c, --config: path to the configuration YAML file (default: "config.yaml")
//	             (can be overridden by IAM_CONFIG_PATH env var)
//	-p, --port: port to bind the HTTP server to (default: 8080)
//	           (can be overridden by IAM_PORT env var)
//	-d, --debug: enable verbose debug output (default: false)
func Main() {
	// Parse command-line flags
	// These can be overridden by environment variables:
	//   - IAM_CONFIG_PATH for configPath
	//   - IAM_PORT for port
	//   - IAM_DEBUG for debug (not implemented yet via env)
	configPath := flag.StringP("config", "c", "config.yaml", "Path to YAML config (overridable via IAM_CONFIG_PATH)")
	port := flag.IntP("port", "p", 8080, "Port to run the HTTP server on (overridable via IAM_PORT)")
	debug := flag.BoolP("debug", "d", false, "Enable debug output (flag only)")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Apply runtime overrides
	cfg.Port = *port
	cfg.Debug = *debug

	// Initialize database
	// db.Init(cfg.Database, cfg.DatabaseDSN)
	_db := db.Init(cfg.Database, cfg.DatabaseDSN)

	// creating new API server instance
	_api := api.New(cfg, _db)

	// Signal handeling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		//syscall.SIGKILL, // "always fatal", "SIGKILL and SIGSTOP may not be caught by a program"
		syscall.SIGHUP, // "terminal is disconnected"
	)

	go func() {
		// Start API server
		fmt.Printf("[%s] Starting on port %d with DB [%s]...\n", cfg.AppName, cfg.Port, cfg.Database)
		if err := _api.StartServer(); err != nil {
			// fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			log.Println("Server error:", err)
			// os.Exit(1)
			sigCh <- syscall.SIGTERM
		}
	}()

	sig := <-sigCh
	log.Println("\nsignal recieved: ", sig.String())

	// Close API for now
	// but later will be done automatically
	_api.StopAndClose()

}
