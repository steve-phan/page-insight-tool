package services

import (
	"page-insight-tool/internal/config"
	analyzer "page-insight-tool/internal/services/analyzer"
)

// Services holds all application services
type Services struct {
	Config   *config.Config
	Analyzer *analyzer.AnalyzerService
}
