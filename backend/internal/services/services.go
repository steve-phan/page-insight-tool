package services

import (
	"github.com/steve-phan/page-insight-tool/internal/config"
	analyzer "github.com/steve-phan/page-insight-tool/internal/services/analyzer"
	"github.com/steve-phan/page-insight-tool/internal/services/health"
	"github.com/steve-phan/page-insight-tool/internal/services/redis"
)

// Services holds all application services
type Services struct {
	Config   *config.Config
	Analyzer *analyzer.AnalyzerService
	Health   *health.HealthService
	Redis    *redis.RedisService
}
