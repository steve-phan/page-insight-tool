package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/steve-phan/page-insight-tool/internal/config"
	"github.com/steve-phan/page-insight-tool/internal/handlers"
	"github.com/steve-phan/page-insight-tool/internal/memcach"
	"github.com/steve-phan/page-insight-tool/internal/routes"
	"github.com/steve-phan/page-insight-tool/internal/services"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	services *services.Services
	handlers *handlers.HandlerFactory
	router   *gin.Engine
	httpSrv  *http.Server
}

// New creates a new server instance following proper dependency flow:
// Config → ServiceFactory → Services → HandlerFactory → Handlers → Routes → Server
func New(cfg *config.Config) (*Server, error) {
	// Create service factory
	serviceFactory := services.NewServiceFactory(cfg)

	// Create services with fail-fast validation
	appServices, err := serviceFactory.CreateServices()
	if err != nil {
		return nil, err
	}

	// Set Gin mode based on environment
	if appServices.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize handlers with services (dependency injection)
	handlerFactory := handlers.NewHandlerFactory(appServices)

	// Initialize routes with handlers (clean dependency flow)
	router := routes.SetupRoutes(handlerFactory)

	httpSrv := &http.Server{
		Addr:           appServices.Config.GetAddress(),
		Handler:        router,
		ReadTimeout:    appServices.Config.Server.ReadTimeout,
		WriteTimeout:   appServices.Config.Server.WriteTimeout,
		IdleTimeout:    appServices.Config.Server.IdleTimeout,
		MaxHeaderBytes: appServices.Config.Server.MaxHeaderBytes,
	}

	return &Server{
		services: appServices,
		handlers: handlerFactory,
		router:   router,
		httpSrv:  httpSrv,
	}, nil
}

// Start starts the server
func (s *Server) Start() error {
	// Start server in a goroutine
	go func() {
		log.Printf("Starting %s v%s on %s", s.services.Config.App.Name, getVersion(), s.services.Config.GetAddress())
		log.Printf("Environment: %s", s.services.Config.App.Environment)
		log.Printf("Build date: %s", getBuildDate())

		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	return nil
}

// Stop gracefully stops the server
func (s *Server) Stop() error {
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpSrv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
		return err
	}

	log.Println("Server exited")
	return nil
}

// WaitForShutdown waits for interrupt signal and shuts down gracefully
func (s *Server) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	memcach.GetMemCache().Stop()

	if err := s.Stop(); err != nil {
		log.Fatalf("Failed to stop server: %v", err)
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
	return time.Now().Local().Format(time.RFC1123)
}
