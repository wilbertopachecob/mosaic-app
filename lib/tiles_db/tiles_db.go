package tiles_db

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	imgpkg "wilbertopachecob/mosaic/lib/img"
)

// TilesDB initializes and populates the tiles database
// Scans the tiles directory for image files and calculates their average colors
// Returns a map of filename to average color [R, G, B]
func TilesDB() map[string][3]float64 {
	logrus.Info("Starting tiles database population")
	
	db := make(map[string][3]float64)
	tilesDir := "tiles"
	
	// Check if tiles directory exists
	if _, err := os.Stat(tilesDir); os.IsNotExist(err) {
		logrus.Warnf("Tiles directory '%s' does not exist", tilesDir)
		return db
	}
	
	// Read tiles directory
	files, err := os.ReadDir(tilesDir)
	if err != nil {
		logrus.WithError(err).Error("Failed to read tiles directory")
		return db
	}
	
	// Process each file in the tiles directory
	for _, file := range files {
		if file.IsDir() {
			continue // Skip subdirectories
		}
		
		filename := file.Name()
		filePath := filepath.Join(tilesDir, filename)
		
		// Check if file is an image
		if !isImageFile(filename) {
			logrus.Debugf("Skipping non-image file: %s", filename)
			continue
		}
		
		// Process the image file
		if err := processImageFile(filePath, db); err != nil {
			logrus.WithError(err).WithField("file", filePath).Error("Failed to process image file")
		}
	}
	
	logrus.WithField("tileCount", len(db)).Info("Tiles database population completed")
	return db
}

// CloneTilesDB creates a deep copy of the tiles database
// This is necessary to avoid concurrent access issues during mosaic generation
func CloneTilesDB(tilesDB map[string][3]float64) map[string][3]float64 {
	db := make(map[string][3]float64, len(tilesDB))
	for k, v := range tilesDB {
		db[k] = v
	}
	return db
}

// isImageFile checks if a filename has an image extension
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp"}
	
	for _, imgExt := range imageExtensions {
		if ext == imgExt {
			return true
		}
	}
	return false
}

// processImageFile processes a single image file and adds it to the database
func processImageFile(filePath string, db map[string][3]float64) error {
	// Open the image file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	// Decode the image
	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}
	
	// Calculate average color
	avgColor := imgpkg.AverageColor(img)
	
	// Add to database
	db[filePath] = avgColor
	
	logrus.WithFields(logrus.Fields{
		"file":   filePath,
		"format": format,
		"color":  avgColor,
	}).Debug("Added tile to database")
	
	return nil
}
