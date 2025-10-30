package services

import (
	"github.com/steve-phan/page-insight-tool/internal/config"
	analyzer "github.com/steve-phan/page-insight-tool/internal/services/analyzer"
	"github.com/steve-phan/page-insight-tool/internal/services/analyzer/extractors"
	"github.com/steve-phan/page-insight-tool/internal/services/health"
	"github.com/steve-phan/page-insight-tool/internal/services/redis"
)

// TestServiceFactory creates services suitable for testing infrastructure components
// This factory bypasses fail-fast validation and provides minimal, working services
type TestServiceFactory struct {
	config *config.Config
}

// NewTestServiceFactory creates a test-oriented service factory
func NewTestServiceFactory(config *config.Config) *TestServiceFactory {
	return &TestServiceFactory{
		config: config,
	}
}

// CreateServices creates minimal services for testing infrastructure
// Uses a single basic extractor to avoid fail-fast validation while keeping tests focused
func (tsf *TestServiceFactory) CreateServices() (*Services, error) {
	// Create test Redis service (no actual connection for tests)
	redisService := redis.NewTestRedisService(tsf.config)

	// For infrastructure tests, we just need a working analyzer with minimal extractors
	// Use only TitleExtractor as it's the simplest and most reliable
	analyzerService, err := analyzer.NewAnalyzerService(tsf.config,
		analyzer.WithExtractors(
			&extractors.TitleExtractor{}, // Minimal, reliable extractor for testing
		))

	if err != nil {
		return nil, err
	}

	// Create health service
	healthService := health.NewHealthService(tsf.config)

	return &Services{
		Config:   tsf.config,
		Analyzer: analyzerService,
		Health:   healthService,
		Redis:    redisService,
	}, nil
}
