package routes

import (
	"net/http"

	"github.com/steve-phan/page-insight-tool/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes with handler factory
// This follows clean dependency flow: Services → Handlers → Routes
func SetupRoutes(handlerFactory *handlers.HandlerFactory) *gin.Engine {
	router := gin.New()

	// Global middleware - ORDER MATTERS!
	router.Use(gin.Logger())

	// Use our custom recovery middleware instead of Gin's default
	router.Use(handlerFactory.ErrorHandler().Recovery())

	// Add our error handling middleware (should be last)
	router.Use(handlerFactory.ErrorHandler().Middleware())

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
