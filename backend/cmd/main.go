package main

import (
	"log"

	"github.com/steve-phan/page-insight-tool/docs/api"
	"github.com/swaggo/swag"

	"github.com/steve-phan/page-insight-tool/internal/config"
	"github.com/steve-phan/page-insight-tool/internal/server"
)

// @title           Page Insight Tool API
// @version         1.0
// @description     A production-grade web page analysis tool that analyzes HTML pages and extracts insights including headings, links, login forms, and more.
// @description     Supports both static HTML analysis and Client-Side Rendered (CSR) website detection.

// @contact.name   API Support
// @contact.email  support@pageinsighttool.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes   http https

// @tag.name  Analysis
// @tag.description  Web page analysis endpoints

// @tag.name  Health
// @tag.description  Health check and monitoring endpoints

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

	swag.Register(api.SwaggerInfo.InstanceName(), api.SwaggerInfo)

	// Wait for shutdown signal
	srv.WaitForShutdown()
}
