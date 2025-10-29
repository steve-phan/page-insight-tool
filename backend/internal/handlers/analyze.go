package handlers

import (
	"net/http"
	"net/url"
	"page-insight-tool/internal/config"
	"page-insight-tool/internal/models"
	analyzer "page-insight-tool/internal/services/analyzer"
	"page-insight-tool/internal/services/analyzer/extractors"
	"strings"

	"github.com/gin-gonic/gin"
)

// AnalyzeHandler handles URL analysis requests with configurable extractors
func AnalyzeHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Query("url")
		if url == "" {
			errorResponse := models.NewInvalidURLError("URL parameter is required")
			c.JSON(errorResponse.Code, errorResponse)
			return
		}

		if !validateUrl(url) {
			errorResponse := models.NewInvalidURLError("Invalid URL format")
			c.JSON(errorResponse.Code, errorResponse)
			return
		}

		service := analyzer.NewAnalyzerService(cfg, analyzer.WithExtractors(
			&extractors.TitleExtractor{},
			&extractors.HeadingsExtractor{},
			&extractors.LinksExtractor{},
			&extractors.LoginFormExtractor{},
			&extractors.VersionExtractor{},
		))

		response, err := service.Analyze(c.Request.Context(), url)
		if err != nil {
			// Determine error type based on the error message
			errorResponse := determineErrorType(err)
			c.JSON(errorResponse.Code, errorResponse)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func validateUrl(rawURL string) bool {
	// Parse the URL using Go's built-in URL parser
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// Check if the URL has a valid scheme (http or https)
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	// Check if the URL has a valid host
	if u.Host == "" {
		return false
	}

	return true
}

// Error pattern constants for mapping
const (
	ErrPatternHTMLParse         = "failed to parse HTML"
	ErrPatternNonOKStatus       = "non-200 status code"
	ErrPatternBodyTooLarge      = "response body too large"
	ErrPatternTimeout           = "timeout"
	ErrPatternContextDeadline   = "context deadline exceeded"
	ErrPatternNoSuchHost        = "no such host"
	ErrPatternConnectionRefused = "connection refused"
	ErrPatternCertificate       = "certificate"
)

// Error mapping for different error scenarios
var errorMappings = map[string]*models.AnalysisError{
	ErrPatternHTMLParse:         models.NewHTMLParseError("Failed to parse HTML content"),
	ErrPatternNonOKStatus:       models.NewInternalServerError("Server returned non-200 status code"),
	ErrPatternBodyTooLarge:      models.NewInternalServerError("Response body exceeds maximum size limit"),
	ErrPatternTimeout:           models.NewTimeoutError("Request timed out"),
	ErrPatternContextDeadline:   models.NewTimeoutError("Request timed out"),
	ErrPatternNoSuchHost:        models.NewPageNotFoundError("Host not found"),
	ErrPatternConnectionRefused: models.NewPageNotFoundError("Connection refused - host may be down"),
	ErrPatternCertificate:       models.NewForbiddenError("SSL certificate error"),
}

// determineErrorType maps analyzer errors to structured error responses using mapping
func determineErrorType(err error) *models.AnalysisError {
	errMsg := err.Error()

	// Check each mapping pattern
	for pattern, errorResponse := range errorMappings {
		if strings.Contains(errMsg, pattern) {
			// Handle special case for non-200 status codes
			if pattern == "non-200 status code" {
				if strings.Contains(errMsg, "403") {
					return models.NewForbiddenError("Access to the URL is forbidden")
				}
				if strings.Contains(errMsg, "404") {
					return models.NewPageNotFoundError("Page not found")
				}
				return models.NewInternalServerError("Server returned non-200 status code")
			}
			return errorResponse
		}
	}

	// Default fallback
	return models.NewInternalServerError("Internal server error occurred")
}
