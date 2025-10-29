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

	// Create server with dependency injection and fail-fast validation
	// The server handles: Config → ServiceFactory → Services → HandlerFactory → Routes
	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Start server
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for shutdown signal
	srv.WaitForShutdown()
}
