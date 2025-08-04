package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// routes creates and configures the HTTP router
func routes() *mux.Router {
	router := mux.NewRouter()

	// Add middleware
	router.Use(loggingMiddleware)
	router.Use(corsMiddleware)

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/file/upload", mosaicHandler).Methods("POST")
	
	// Health check endpoint
	api.HandleFunc("/health", healthHandler).Methods("GET")

	// Serve static files (frontend build)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("frontend/build")))

	return router
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Call the next handler
		next.ServeHTTP(w, r)
		
		// Log the request
		logrus.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"duration":   time.Since(start),
			"user_agent": r.UserAgent(),
		}).Info("HTTP request processed")
	})
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// healthHandler provides a health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "service": "mosaic-app"}`))
}
