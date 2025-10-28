package main

import (
	"log"

	"page-insight-tool/internal/config"
	"page-insight-tool/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start server
	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for shutdown signal
	srv.WaitForShutdown()
}
