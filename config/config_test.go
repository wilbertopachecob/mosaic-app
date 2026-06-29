package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("SERVER_PORT", "")
	t.Setenv("MAX_FILE_SIZE", "")
	t.Setenv("TILES_DIR", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("MAX_IMAGE_WIDTH", "")
	t.Setenv("MAX_IMAGE_HEIGHT", "")
	t.Setenv("MIN_TILE_SIZE", "")
	t.Setenv("MAX_TILE_SIZE", "")
	t.Setenv("MAX_CONCURRENT_GENERATIONS", "")
	t.Setenv("RATE_LIMIT_REQUESTS", "")
	t.Setenv("RATE_LIMIT_WINDOW", "")
	t.Setenv("REQUEST_TIMEOUT", "")
	t.Setenv("PRODUCTION", "")

	cfg := Load()

	assert.Equal(t, "8080", cfg.ServerPort)
	assert.Equal(t, int64(10*1024*1024), cfg.MaxFileSize)
	assert.Equal(t, "tiles", cfg.TilesDir)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, 4096, cfg.MaxImageWidth)
	assert.Equal(t, 4096, cfg.MaxImageHeight)
	assert.Equal(t, 5, cfg.MinTileSize)
	assert.Equal(t, 100, cfg.MaxTileSize)
	assert.Equal(t, 2, cfg.MaxConcurrentGenerations)
	assert.Equal(t, 10, cfg.RateLimitRequests)
	assert.Equal(t, time.Hour, cfg.RateLimitWindow)
	assert.Equal(t, 120*time.Second, cfg.RequestTimeout)
	assert.False(t, cfg.Production)
}

func TestLoadFromEnvironment(t *testing.T) {
	t.Setenv("SERVER_PORT", "3000")
	t.Setenv("MAX_FILE_SIZE", "2048")
	t.Setenv("TILES_DIR", "custom-tiles")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("MAX_IMAGE_WIDTH", "800")
	t.Setenv("MAX_IMAGE_HEIGHT", "600")
	t.Setenv("MIN_TILE_SIZE", "10")
	t.Setenv("MAX_TILE_SIZE", "50")
	t.Setenv("MAX_CONCURRENT_GENERATIONS", "4")
	t.Setenv("RATE_LIMIT_REQUESTS", "20")
	t.Setenv("RATE_LIMIT_WINDOW", "60")
	t.Setenv("REQUEST_TIMEOUT", "30")
	t.Setenv("PRODUCTION", "true")

	cfg := Load()

	assert.Equal(t, "3000", cfg.ServerPort)
	assert.Equal(t, int64(2048), cfg.MaxFileSize)
	assert.Equal(t, "custom-tiles", cfg.TilesDir)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, 800, cfg.MaxImageWidth)
	assert.Equal(t, 600, cfg.MaxImageHeight)
	assert.Equal(t, 10, cfg.MinTileSize)
	assert.Equal(t, 50, cfg.MaxTileSize)
	assert.Equal(t, 4, cfg.MaxConcurrentGenerations)
	assert.Equal(t, 20, cfg.RateLimitRequests)
	assert.Equal(t, time.Minute, cfg.RateLimitWindow)
	assert.Equal(t, 30*time.Second, cfg.RequestTimeout)
	assert.True(t, cfg.Production)
}

func TestEnvHelpersIgnoreInvalidValues(t *testing.T) {
	t.Setenv("MAX_FILE_SIZE", "not-a-number")
	t.Setenv("MAX_IMAGE_WIDTH", "bad")
	t.Setenv("PRODUCTION", "maybe")

	cfg := Load()

	assert.Equal(t, int64(10*1024*1024), cfg.MaxFileSize)
	assert.Equal(t, 4096, cfg.MaxImageWidth)
	assert.False(t, cfg.Production)
}

func TestGetEnvWithDefault(t *testing.T) {
	key := "MOSAIC_TEST_STRING"
	require.NoError(t, os.Unsetenv(key))
	assert.Equal(t, "fallback", getEnvWithDefault(key, "fallback"))

	require.NoError(t, os.Setenv(key, "value"))
	t.Cleanup(func() { _ = os.Unsetenv(key) })
	assert.Equal(t, "value", getEnvWithDefault(key, "fallback"))
}
