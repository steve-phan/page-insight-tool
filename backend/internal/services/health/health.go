package health

import (
	"os"
	"github.com/steve-phan/page-insight-tool/internal/config"
	"time"
)

type HealthService struct {
	cfg *config.Config
}

func NewHealthService(cfg *config.Config) *HealthService {
	return &HealthService{cfg: cfg}
}
func (hs *HealthService) CheckHealth() map[string]string {
	return map[string]string{
		"status":      "healthy",
		"timestamp":   time.Now().UTC().String(),
		"version":     getVersion(),
		"build_date":  getBuildDate(),
		"environment": hs.cfg.App.Environment,
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
