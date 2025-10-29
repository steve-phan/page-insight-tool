package services

import (
	"page-insight-tool/internal/config"
	analyzer "page-insight-tool/internal/services/analyzer"
	"page-insight-tool/internal/services/health"
)

// Services holds all application services
type Services struct {
	Config   *config.Config
	Analyzer *analyzer.AnalyzerService
	Health   *health.HealthService
}
