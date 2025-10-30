package services

import (
	"fmt"

	"github.com/steve-phan/page-insight-tool/internal/config"
	analyzer "github.com/steve-phan/page-insight-tool/internal/services/analyzer"
	"github.com/steve-phan/page-insight-tool/internal/services/analyzer/extractors"
	"github.com/steve-phan/page-insight-tool/internal/services/health"
	"github.com/steve-phan/page-insight-tool/internal/services/redis"
)

// ServiceFactory handles creation and validation of all application services
type ServiceFactory struct {
	config *config.Config
}

// NewServiceFactory creates a new service factory
func NewServiceFactory(config *config.Config) *ServiceFactory {
	return &ServiceFactory{
		config: config,
	}
}

// CreateServices creates and validates all application services
// This is where we implement fail-fast validation
func (sf *ServiceFactory) CreateServices() (*Services, error) {
	// Create Redis service
	redisService, err := redis.NewRedisService(sf.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis service: %w", err)
	}

	// Create analyzer service with configured extractors
	analyzerService, err := analyzer.NewAnalyzerService(sf.config,
		analyzer.WithExtractors(
			&extractors.TitleExtractor{},
			&extractors.HeadingsExtractor{},
			&extractors.LinksExtractor{},
			&extractors.LoginFormExtractor{},
			&extractors.VersionExtractor{},
		))

	if err != nil {
		return nil, fmt.Errorf("failed to create analyzer service: %w", err)
	}

	// Create health service
	healthService := health.NewHealthService(sf.config)

	return &Services{
		Config:   sf.config,
		Analyzer: analyzerService,
		Health:   healthService,
		Redis:    redisService,
	}, nil
}
