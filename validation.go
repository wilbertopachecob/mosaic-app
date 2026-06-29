package main

import (
	"fmt"
	"image"
	"strconv"

	"wilbertopachecob/mosaic/config"
	"wilbertopachecob/mosaic/lib/mosaic"
)

const defaultTileSize = 20

func parseTileSize(raw string, cfg *config.Config) (int, error) {
	if raw == "" {
		return defaultTileSize, nil
	}

	tileSize, err := strconv.Atoi(raw)
	if err != nil {
		return defaultTileSize, nil
	}

	if tileSize < cfg.MinTileSize || tileSize > cfg.MaxTileSize {
		return 0, fmt.Errorf("tileSize must be between %d and %d", cfg.MinTileSize, cfg.MaxTileSize)
	}

	return tileSize, nil
}

func parseBlend(raw string) (float64, error) {
	blend := mosaic.DefaultOptions().SourceBlend
	if raw == "" {
		return blend, nil
	}

	parsed, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("blend must be a number between 0.0 and 1.0")
	}

	if parsed < 0 || parsed > 1 {
		return 0, fmt.Errorf("blend must be between 0.0 and 1.0")
	}

	return parsed, nil
}

func validateImageDimensions(img image.Image, cfg *config.Config) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width <= 0 || height <= 0 {
		return fmt.Errorf("image has invalid dimensions")
	}
	if width > cfg.MaxImageWidth || height > cfg.MaxImageHeight {
		return fmt.Errorf("image dimensions exceed maximum of %dx%d pixels", cfg.MaxImageWidth, cfg.MaxImageHeight)
	}

	return nil
}
