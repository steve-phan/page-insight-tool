package handlers

import (
	"net/http"
	"page-insight-tool/internal/middleware"
	analyzer "page-insight-tool/internal/services/analyzer"
	"page-insight-tool/internal/validation"

	"github.com/gin-gonic/gin"
)

// AnalyzeHandler handles URL analysis requests with clean error handling
func AnalyzeHandler(analyzerService *analyzer.AnalyzerService, errorHandler *middleware.ErrorHandler, urlValidator *validation.URLValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract and validate URL parameter
		rawURL := c.Query("url")

		// Validate URL using dedicated validator
		if err := urlValidator.ValidateURL(rawURL); err != nil {
			errorHandler.HandleError(c, err)
			return
		}

		// Perform analysis using the pre-configured analyzer service
		response, err := analyzerService.Analyze(c.Request.Context(), rawURL)
		if err != nil {
			// Let the error handler deal with mapping and response
			errorHandler.HandleError(c, err)
			return
		}

		// Success response
		c.JSON(http.StatusOK, response)
	}
}
