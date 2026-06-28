package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"wilbertopachecob/mosaic/config"
)

func testConfig() *config.Config {
	return &config.Config{
		MaxFileSize:    10 * 1024 * 1024,
		MaxImageWidth:  4096,
		MaxImageHeight: 4096,
		MinTileSize:    5,
		MaxTileSize:    100,
	}
}

func TestParseTileSize(t *testing.T) {
	cfg := testConfig()

	size, err := parseTileSize("", cfg)
	require.NoError(t, err)
	assert.Equal(t, defaultTileSize, size)

	size, err = parseTileSize("invalid", cfg)
	require.NoError(t, err)
	assert.Equal(t, defaultTileSize, size)

	size, err = parseTileSize("20", cfg)
	require.NoError(t, err)
	assert.Equal(t, 20, size)

	_, err = parseTileSize("1", cfg)
	assert.Error(t, err)

	_, err = parseTileSize("200", cfg)
	assert.Error(t, err)
}

func TestParseBlend(t *testing.T) {
	blend, err := parseBlend("")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, blend, 0.0)
	assert.LessOrEqual(t, blend, 1.0)

	blend, err = parseBlend("0.55")
	require.NoError(t, err)
	assert.Equal(t, 0.55, blend)

	_, err = parseBlend("not-a-number")
	assert.Error(t, err)

	_, err = parseBlend("1.5")
	assert.Error(t, err)
}

func TestValidateImageDimensions(t *testing.T) {
	cfg := testConfig()

	err := validateImageDimensions(createTestImage(100, 100), cfg)
	assert.NoError(t, err)

	err = validateImageDimensions(createTestImage(5000, 100), cfg)
	assert.Error(t, err)

	err = validateImageDimensions(createTestImage(100, 5000), cfg)
	assert.Error(t, err)
}
