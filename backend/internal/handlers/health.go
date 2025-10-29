package handlers

import (
	"net/http"

	"page-insight-tool/internal/config"
	"page-insight-tool/internal/services/health"

	"github.com/gin-gonic/gin"
)

// HealthHandler returns the health check handler
func HealthHandler(healthService *health.HealthService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, healthService.CheckHealth())
	}
}
