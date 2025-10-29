package handlers

import (
	"net/http"

	"github.com/steve-phan/page-insight-tool/internal/config"
	"github.com/steve-phan/page-insight-tool/internal/services/health"

	"github.com/gin-gonic/gin"
)

// HealthHandler returns the health check handler
// @Summary      Health check endpoint
// @Description  Returns the health status of the API including version, environment, and server status
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {object}  health.HealthResponse
// @Router       /health [get]
func HealthHandler(healthService *health.HealthService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, healthService.CheckHealth())
	}
}
