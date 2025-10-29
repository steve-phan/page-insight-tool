package handlers

import (
	"page-insight-tool/internal/config"
	"page-insight-tool/internal/services"
	analyzer "page-insight-tool/internal/services/analyzer"

	"github.com/gin-gonic/gin"
)

// HandlerFactory encapsulates all handler dependencies
// This follows the factory pattern and provides clean dependency injection
type HandlerFactory struct {
	config   *config.Config
	analyzer *analyzer.AnalyzerService
}

// NewHandlerFactory creates a new handler factory with dependencies
func NewHandlerFactory(services *services.Services) *HandlerFactory {
	return &HandlerFactory{
		config:   services.Config,
		analyzer: services.Analyzer,
	}
}

// HealthHandler returns the health check handler
func (hf *HandlerFactory) HealthHandler() gin.HandlerFunc {
	return HealthHandler(hf.config)
}

// AnalyzeHandler returns the analyze handler
func (hf *HandlerFactory) AnalyzeHandler() gin.HandlerFunc {
	return AnalyzeHandler(hf.analyzer)
}
