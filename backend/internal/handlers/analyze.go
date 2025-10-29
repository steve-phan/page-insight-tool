package handlers

import (
	"net/http"

	"github.com/steve-phan/page-insight-tool/internal/middleware"
	analyzer "github.com/steve-phan/page-insight-tool/internal/services/analyzer"
	"github.com/steve-phan/page-insight-tool/internal/validation"

	"github.com/gin-gonic/gin"
)

// AnalyzeHandler handles URL analysis requests with clean error handling
// @Summary      Analyze a web page
// @Description  Analyzes a web page and extracts HTML version, title, headings, links, login forms, and CSR detection information
// @Tags         Analysis
// @Accept       json
// @Produce      json
// @Param        url   query     string  true  "URL of the web page to analyze"  example(https://example.com)
// @Success      200   {object}  models.AnalysisResponse
// @Failure      400   {object}  models.HTTPError  "Invalid URL"
// @Failure      422   {object}  models.HTTPError  "HTML parsing error"
// @Failure      429   {object}  models.HTTPError  "Rate limit exceeded"
// @Failure      500   {object}  models.HTTPError  "Internal server error"
// @Router       /analyze [get]
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
			errorHandler.HandleError(c, err)
			return
		}

		// Success response
		c.JSON(http.StatusOK, response)
	}
}
