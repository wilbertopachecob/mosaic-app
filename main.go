package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wilbertopachecob/mosaic/config"
	"wilbertopachecob/mosaic/lib/tiles_db"
)

// Global tiles database - initialized at startup
var tilesDB map[string][3]float64

// main is the entry point of the application
func main() {
	// Load configuration
	cfg := config.Load()
	
	// Initialize tiles database
	log.Println("Initializing tiles database...")
	tilesDB = tiles_db.TilesDB()
	log.Printf("Tiles database initialized with %d tiles", len(tilesDB))

	// Create router
	router := routes()

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Mosaic server starting on http://localhost:%s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
