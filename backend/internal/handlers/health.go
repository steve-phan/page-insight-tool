package handlers

import (
	"net/http"
	"os"
	"time"

	"page-insight-tool/internal/config"

	"github.com/gin-gonic/gin"
)

// HealthHandler returns the health check handler
func HealthHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":        "healthy",
			"timestamp":     time.Now().UTC(),
			"version":       getVersion(),
			"build_date":    getBuildDate(),
			"environment":   cfg.App.Environment,
			"config_loaded": true,
		})
	}
}

// Version information - these will be set by the build system
func getVersion() string {
	if version := os.Getenv("VERSION"); version != "" {
		return version
	}
	return "dev"
}

func getBuildDate() string {
	if buildDate := os.Getenv("BUILD_DATE"); buildDate != "" {
		return buildDate
	}
	return "unknown"
}
