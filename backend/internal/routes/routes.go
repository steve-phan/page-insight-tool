package routes

import (
	"net/http"

	"page-insight-tool/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes with handler factory
// This follows clean dependency flow: Services → Handlers → Routes
func SetupRoutes(handlerFactory *handlers.HandlerFactory) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// Basic hardening can be added later (security headers, body size limits)

	// API routes
	setupAPIRoutes(router, handlerFactory)

	// Handle OPTIONS requests globally
	router.OPTIONS("/*path", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	return router
}

// setupAPIRoutes configures API v1 routes with handler factory
func setupAPIRoutes(router *gin.Engine, handlerFactory *handlers.HandlerFactory) {
	api := router.Group("/api/v1")
	{
		api.GET("/health", handlerFactory.HealthHandler())
		api.GET("/analyze", handlerFactory.AnalyzeHandler())
	}
}
