// cmd/server/main.go
package server

import (
	"flag"
	"fmt"
	"log"

	"github.com/javadmohebbi/goIAM/internal/api"
	"github.com/javadmohebbi/goIAM/internal/config"
)

func Main() {
	// CLI flags
	configPath := flag.String("config", "config.yaml", "Path to configuration YAML file")
	port := flag.Int("port", 8080, "Port to run the HTTP server on")
	debug := flag.Bool("debug", false, "Enable debug output")
	flag.Parse()

	// Load YAML config
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if *port != 0 {
		cfg.Port = *port
	}
	cfg.Debug = *debug

	fmt.Printf("[goIAM] Starting on port %d (debug: %v)\n", cfg.Port, cfg.Debug)

	if err := api.StartServer(cfg); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
