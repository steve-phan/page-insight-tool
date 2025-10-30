package routes

import (
	"net/http"
	"time"

	_ "github.com/steve-phan/page-insight-tool/docs/api" // Swagger docs
	"github.com/steve-phan/page-insight-tool/internal/handlers"
	"github.com/steve-phan/page-insight-tool/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	setupAPIRoutes(router, handlerFactory)

	// Handle OPTIONS requests globally
	router.OPTIONS("/*path", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	return router
}

// setupAPIRoutes configures API v1 routes with handler factory
func setupAPIRoutes(router *gin.Engine, handlerFactory *handlers.HandlerFactory) {
	api := router.Group("/api/v1")

	rateLimiter := middleware.NewRedisRateLimiter(handlerFactory.RedisService().GetClient())

	{
		// Health endpoint: More lenient rate limit (100 requests per minute)
		healthGroup := api.Group("/health")
		healthGroup.Use(rateLimiter.RateLimit(100, time.Minute))
		healthGroup.GET("", handlerFactory.HealthHandler())

		// Analyze endpoint: Stricter rate limit (5 requests per 10 seconds)
		analyzeGroup := api.Group("/analyze")
		analyzeGroup.Use(rateLimiter.RateLimit(5, 10*time.Second))
		analyzeGroup.GET("", handlerFactory.AnalyzeHandler())
	}
}
