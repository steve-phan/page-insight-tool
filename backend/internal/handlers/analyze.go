package handlers

import (
	"net/http"
	"page-insight-tool/internal/config"
	"page-insight-tool/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func AnalyzeHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Query("url")
		if url == "" || !validateUrl(url) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid URL",
			})
			return
		}

		response, err := services.Analyze(url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, response)

	}
}

func validateUrl(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
