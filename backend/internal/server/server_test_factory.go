package server

import (
	"net/http"

	"github.com/steve-phan/page-insight-tool/internal/config"
	"github.com/steve-phan/page-insight-tool/internal/handlers"
	"github.com/steve-phan/page-insight-tool/internal/routes"
	"github.com/steve-phan/page-insight-tool/internal/services"

	"github.com/gin-gonic/gin"
)

// NewForTesting creates a server instance using test factory for infrastructure testing
// This bypasses the fail-fast validation to focus on server infrastructure tests
func NewForTesting(cfg *config.Config) (*Server, error) {
	// Use test factory for infrastructure tests
	serviceFactory := services.NewTestServiceFactory(cfg)
	srvs, err := serviceFactory.CreateServices()
	if err != nil {
		return nil, err
	}

	// Set Gin mode based on environment
	if srvs.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create handler factory
	handlerFactory := handlers.NewHandlerFactory(srvs)

	// Setup routes with handlers
	router := routes.SetupRoutes(handlerFactory)

	httpSrv := &http.Server{
		Addr:           srvs.Config.GetAddress(),
		Handler:        router,
		ReadTimeout:    srvs.Config.Server.ReadTimeout,
		WriteTimeout:   srvs.Config.Server.WriteTimeout,
		IdleTimeout:    srvs.Config.Server.IdleTimeout,
		MaxHeaderBytes: srvs.Config.Server.MaxHeaderBytes,
	}

	return &Server{
		services: srvs,
		handlers: handlerFactory,
		router:   router,
		httpSrv:  httpSrv,
	}, nil
}
