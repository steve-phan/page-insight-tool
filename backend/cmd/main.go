package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Version information set during build
var (
	version   = "dev"     // set by build system
	buildDate = "unknown" // set by build system
	gitCommit = "unknown" // set by build system
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Basic middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check endpoint with version info
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":     "healthy",
			"timestamp":  time.Now().UTC(),
			"version":    version,
			"build_date": buildDate,
			"git_commit": gitCommit,
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		api.POST("/analyze", analyzeHandler)
		api.GET("/stats", statsHandler)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Page Insight Tool v%s on port %s", version, port)
	log.Printf("Build date: %s, Git commit: %s", buildDate, gitCommit)
	
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Placeholder handlers for now
func analyzeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Analysis endpoint - to be implemented",
		"status":  "pending",
	})
}

func statsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Stats endpoint - to be implemented",
		"status":  "pending",
	})
}
