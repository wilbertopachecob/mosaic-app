package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	ServerPort  string
	MaxFileSize int64
	TilesDir    string
	LogLevel    string

	// Image validation
	MaxImageWidth  int
	MaxImageHeight int

	// Mosaic parameter bounds (must match frontend tileValues)
	MinTileSize int
	MaxTileSize int

	// Abuse prevention
	MaxConcurrentGenerations int
	RateLimitRequests        int
	RateLimitWindow          time.Duration
	RequestTimeout           time.Duration

	// When true, API errors omit internal details
	Production bool
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	rateLimitWindowSec := getEnvAsIntWithDefault("RATE_LIMIT_WINDOW", 3600)

	config := &Config{
		ServerPort:               getEnvWithDefault("SERVER_PORT", "8080"),
		MaxFileSize:              getEnvAsInt64WithDefault("MAX_FILE_SIZE", 10*1024*1024),
		TilesDir:                 getEnvWithDefault("TILES_DIR", "tiles"),
		LogLevel:                 getEnvWithDefault("LOG_LEVEL", "info"),
		MaxImageWidth:            getEnvAsIntWithDefault("MAX_IMAGE_WIDTH", 4096),
		MaxImageHeight:           getEnvAsIntWithDefault("MAX_IMAGE_HEIGHT", 4096),
		MinTileSize:              getEnvAsIntWithDefault("MIN_TILE_SIZE", 5),
		MaxTileSize:              getEnvAsIntWithDefault("MAX_TILE_SIZE", 100),
		MaxConcurrentGenerations: getEnvAsIntWithDefault("MAX_CONCURRENT_GENERATIONS", 2),
		RateLimitRequests:        getEnvAsIntWithDefault("RATE_LIMIT_REQUESTS", 10),
		RateLimitWindow:          time.Duration(rateLimitWindowSec) * time.Second,
		RequestTimeout:           time.Duration(getEnvAsIntWithDefault("REQUEST_TIMEOUT", 120)) * time.Second,
		Production:               getEnvAsBoolWithDefault("PRODUCTION", false),
	}

	return config
}

// getEnvWithDefault gets an environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt64WithDefault gets an environment variable as int64 with a default value
func getEnvAsInt64WithDefault(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBoolWithDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
