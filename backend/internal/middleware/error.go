package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/steve-phan/page-insight-tool/internal/errors"
	"github.com/steve-phan/page-insight-tool/internal/models"

	"github.com/gin-gonic/gin"
)

// ErrorHandler provides centralized error handling middleware for Gin
type ErrorHandler struct {
	mapper *errors.ErrorMapper
}

// NewErrorHandler creates a new error handling middleware
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		mapper: errors.NewErrorMapper(),
	}
}

// Middleware returns the Gin middleware function for error handling
func (eh *ErrorHandler) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process the request
		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			// Get the last error (most recent)
			err := c.Errors.Last().Err

			// Map to HTTP error
			httpError := eh.mapper.MapToHTTPError(err)
			if httpError == nil {
				// Fallback for nil errors
				httpError = &models.HTTPError{
					Message:   "Unknown error occurred",
					Type:      string(errors.ErrorTypeInternal),
					Code:      http.StatusInternalServerError,
					Timestamp: time.Now(),
				}
			}

			// Add request context
			eh.enrichErrorContext(c, httpError)

			// Log the error (structured logging would be better)
			eh.logError(c, err, httpError)

			// Send error response
			eh.sendErrorResponse(c, httpError)
		}
	}
}

// HandleError is a utility to manually trigger error handling from handlers
func (eh *ErrorHandler) HandleError(c *gin.Context, err error) {
	if err != nil {
		// Add error to Gin's error collection
		_ = c.Error(err)

		// Map and send response immediately
		httpError := eh.mapper.MapToHTTPError(err)
		eh.enrichErrorContext(c, httpError)
		eh.logError(c, err, httpError)
		eh.sendErrorResponse(c, httpError)
	}
}

// enrichErrorContext adds request context to the error
func (eh *ErrorHandler) enrichErrorContext(c *gin.Context, httpError *models.HTTPError) {
	// Add request ID if available (from request ID middleware)
	if requestID := c.GetString("request_id"); requestID != "" {
		httpError.RequestID = requestID
	}

	// Add timestamp if not set
	if httpError.Timestamp.IsZero() {
		httpError.Timestamp = time.Now()
	}

	// Add additional context to details if needed
	if httpError.Details == nil {
		httpError.Details = make(map[string]interface{})
	}

	// Add request path for debugging (in development)
	if gin.Mode() != gin.ReleaseMode {
		httpError.Details["path"] = c.Request.URL.Path
		httpError.Details["method"] = c.Request.Method
	}
}

// logError logs the error with appropriate context
func (eh *ErrorHandler) logError(c *gin.Context, originalErr error, httpError *models.HTTPError) {
	// In production, this should use structured logging (e.g., logrus, zap)
	logLevel := eh.determineLogLevel(httpError.Code)

	requestID := httpError.RequestID
	if requestID == "" {
		requestID = "unknown"
	}

	switch logLevel {
	case "error":
		log.Printf("[ERROR] [%s] %s %s - %v (mapped to: %s)",
			requestID, c.Request.Method, c.Request.URL.Path, originalErr, httpError.Message)
	case "warn":
		log.Printf("[WARN] [%s] %s %s - %v (mapped to: %s)",
			requestID, c.Request.Method, c.Request.URL.Path, originalErr, httpError.Message)
	default:
		log.Printf("[INFO] [%s] %s %s - %v (mapped to: %s)",
			requestID, c.Request.Method, c.Request.URL.Path, originalErr, httpError.Message)
	}
}

// determineLogLevel determines appropriate log level based on HTTP status code
func (eh *ErrorHandler) determineLogLevel(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "error"
	case statusCode >= 400:
		return "warn"
	default:
		return "info"
	}
}

// sendErrorResponse sends the standardized error response
func (eh *ErrorHandler) sendErrorResponse(c *gin.Context, httpError *models.HTTPError) {
	// Don't double-send responses
	if c.Writer.Written() {
		return
	}

	// Create standardized response
	response := models.NewHTTPErrorResponse(httpError)

	// Send JSON response with appropriate status code
	c.JSON(httpError.Code, response)

	// Abort to prevent further processing
	c.Abort()
}

// Recovery middleware for panic handling
func (eh *ErrorHandler) Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Convert panic to error
		err := &errors.DomainError{
			Type:       errors.ErrorTypeInternal,
			Message:    "Internal server error - panic occurred",
			StatusCode: http.StatusInternalServerError,
			Details: map[string]interface{}{
				"panic": recovered,
			},
		}

		// Handle the error through normal error handling
		eh.HandleError(c, err)
	})
}
