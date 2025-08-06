package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	ServerPort  string
	MaxFileSize int64
	TilesDir    string
	LogLevel    string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	config := &Config{
		ServerPort:  getEnvWithDefault("SERVER_PORT", "8080"),
		MaxFileSize: getEnvAsInt64WithDefault("MAX_FILE_SIZE", 10*1024*1024), // 10MB default
		TilesDir:    getEnvWithDefault("TILES_DIR", "tiles"),
		LogLevel:    getEnvWithDefault("LOG_LEVEL", "info"),
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
