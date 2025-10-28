package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"page-insight-tool/internal/config"
	"page-insight-tool/internal/routes"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	config  *config.Config
	router  *gin.Engine
	httpSrv *http.Server
}

// New creates a new server instance
func New(cfg *config.Config) *Server {
	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := routes.SetupRoutes(cfg)

	httpSrv := &http.Server{
		Addr:           cfg.GetAddress(),
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	return &Server{
		config:  cfg,
		router:  router,
		httpSrv: httpSrv,
	}
}

// Start starts the server
func (s *Server) Start() error {
	// Start server in a goroutine
	go func() {
		log.Printf("Starting %s v%s on %s", s.config.App.Name, getVersion(), s.config.GetAddress())
		log.Printf("Environment: %s", s.config.App.Environment)
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
