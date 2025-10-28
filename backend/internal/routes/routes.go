package routes

import (
	"net/http"

	"page-insight-tool/internal/config"
	"page-insight-tool/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(cfg *config.Config) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// Basic hardening can be added later (security headers, body size limits)

	// API routes
	setupAPIRoutes(router, cfg)

	// Handle OPTIONS requests globally
	router.OPTIONS("/*path", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	return router
}

// setupAPIRoutes configures API v1 routes
func setupAPIRoutes(router *gin.Engine, cfg *config.Config) {
	api := router.Group("/api/v1")
	{
		api.GET("/health", handlers.HealthHandler(cfg))
		api.GET("/analyze", handlers.AnalyzeHandler(cfg))
	}
}
