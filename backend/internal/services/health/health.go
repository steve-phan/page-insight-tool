package health

import (
	"os"
	"time"

	"github.com/steve-phan/page-insight-tool/internal/config"
)

type HealthService struct {
	cfg *config.Config
}

func NewHealthService(cfg *config.Config) *HealthService {
	return &HealthService{cfg: cfg}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status      string `json:"status" example:"healthy"`
	Timestamp   string `json:"timestamp" example:"2025-10-29T22:00:00Z"`
	Version     string `json:"version" example:"1.0.0"`
	BuildDate   string `json:"build_date" example:"2025-10-29T22:00:00Z"`
	Environment string `json:"environment" example:"development"`
}

func (hs *HealthService) CheckHealth() HealthResponse {
	environment := "unknown"
	if hs != nil && hs.cfg != nil {
		environment = hs.cfg.App.Environment
		if environment == "" {
			environment = "unknown"
		}
	}

	return HealthResponse{
		Status:      "healthy",
		Timestamp:   time.Now().UTC().String(),
		Version:     getVersion(),
		BuildDate:   getBuildDate(),
		Environment: environment,
	}
}

// Version information - these will be set by the build system
func getVersion() string {
	if version := os.Getenv("VERSION"); version != "" {
		return version
	}
	return "dev"
}

func getBuildDate() string {
	if buildDate := os.Getenv("BUILD_DATE"); buildDate != "" {
		return buildDate
	}
	return "unknown"
}
