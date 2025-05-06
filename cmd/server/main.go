// Package server starts the goIAM HTTP API service.
package server

import (
	"fmt"
	"os"

	"github.com/javadmohebbi/goIAM/internal/api"
	"github.com/javadmohebbi/goIAM/internal/config"
	"github.com/javadmohebbi/goIAM/internal/db"
	flag "github.com/spf13/pflag"
)

// Main is the entry point for starting the goIAM API server.
// It parses command-line flags, loads configuration from a YAML file,
// initializes the database, and starts the HTTP server.
//
// Flags:
//
//	-c, --config: path to the configuration YAML file (default: "config.yaml")
//	-p, --port: port to bind the HTTP server to (default: 8080)
//	-d, --debug: enable verbose debug output (default: false)
func Main() {
	// Parse command-line flags
	configPath := flag.StringP("config", "c", "config.yaml", "Path to YAML config")
	port := flag.IntP("port", "p", 8080, "Port to run the HTTP server on")
	debug := flag.BoolP("debug", "d", false, "Enable debug output")
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
	db.Init(cfg.Database, cfg.DatabaseDSN)

	// Start API server
	fmt.Printf("[goIAM] Starting on port %d with DB [%s]...\n", cfg.Port, cfg.Database)
	if err := api.StartServer(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
