package server

import (
	"fmt"
	"os"

	"github.com/javadmohebbi/goIAM/internal/api"
	"github.com/javadmohebbi/goIAM/internal/config"
	"github.com/javadmohebbi/goIAM/internal/db"
	flag "github.com/spf13/pflag"
)

func Main() {
	// CLI flags
	configPath := flag.StringP("config", "c", "config.yaml", "Path to YAML config")
	port := flag.IntP("port", "p", 8080, "Port to run the HTTP server on")
	debug := flag.BoolP("debug", "d", false, "Enable debug output")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	cfg.Port = *port
	cfg.Debug = *debug

	db.Init(cfg.Database, cfg.DatabaseDSN)

	fmt.Printf("[goIAM] Starting on port %d with DB [%s]...\n", cfg.Port, cfg.Database)
	if err := api.StartServer(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
