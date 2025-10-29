package handlers

import (
	"page-insight-tool/internal/config"
	"page-insight-tool/internal/middleware"
	"page-insight-tool/internal/services"
	analyzer "page-insight-tool/internal/services/analyzer"
	"page-insight-tool/internal/services/health"
	"page-insight-tool/internal/validation"

	"github.com/gin-gonic/gin"
)

// HandlerFactory encapsulates all handler dependencies
// This follows the factory pattern and provides clean dependency injection
type HandlerFactory struct {
	config       *config.Config
	analyzer     *analyzer.AnalyzerService
	health       *health.HealthService
	errorHandler *middleware.ErrorHandler
	urlValidator *validation.URLValidator
}

// NewHandlerFactory creates a new handler factory with dependencies
func NewHandlerFactory(services *services.Services) *HandlerFactory {
	return &HandlerFactory{
		config:       services.Config,
		analyzer:     services.Analyzer,
		health:       services.Health,
		errorHandler: middleware.NewErrorHandler(),
		urlValidator: validation.NewURLValidator(),
	}
}

// ErrorHandler returns the error handling middleware
func (hf *HandlerFactory) ErrorHandler() *middleware.ErrorHandler {
	return hf.errorHandler
}

// HealthHandler returns the health check handler
func (hf *HandlerFactory) HealthHandler() gin.HandlerFunc {
	return HealthHandler(hf.health, hf.config)
}

// AnalyzeHandler returns the analyze handler with error handling and validation
func (hf *HandlerFactory) AnalyzeHandler() gin.HandlerFunc {
	return AnalyzeHandler(hf.analyzer, hf.errorHandler, hf.urlValidator)
}
