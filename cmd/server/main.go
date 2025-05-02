package server

import (
	"fmt"
	"os"

	"github.com/javadmohebbi/goIAM/internal/api"
	"github.com/javadmohebbi/goIAM/internal/config"
	flag "github.com/spf13/pflag"
)

func Main() {
	// Define flags (short + long)
	configPath := flag.StringP("config", "c", "config.yaml", "Path to configuration YAML file")
	port := flag.IntP("port", "p", 8080, "Port to run the HTTP server on")
	debug := flag.BoolP("debug", "d", false, "Enable debug output")

	flag.Parse()

	// Load config
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}
	cfg.Port = *port
	cfg.Debug = *debug

	fmt.Printf("[goIAM] Starting on port %d (debug: %v)\n", cfg.Port, cfg.Debug)

	if err := api.StartServer(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
