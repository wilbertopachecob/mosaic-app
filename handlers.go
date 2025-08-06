package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	imgpkg "wilbertopachecob/mosaic/lib/img"
	"wilbertopachecob/mosaic/lib/tiles_db"

	"github.com/sirupsen/logrus"
)

// mosaicHandler handles the mosaic generation request
func mosaicHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		sendErrorResponse(w, "Failed to parse form", err, http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("imgUpload")
	if err != nil {
		sendErrorResponse(w, "Failed to get uploaded file", err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get tile size parameter
	tileSizeStr := r.FormValue("tileSize")
	tileSize, err := strconv.Atoi(tileSizeStr)
	if err != nil || tileSize <= 0 {
		tileSize = 20 // Default tile size
	}

	// Log request details
	logrus.WithFields(logrus.Fields{
		"fileName": header.Filename,
		"fileSize": header.Size,
		"tileSize": tileSize,
	}).Info("Processing mosaic request")

	// Decode original image
	original, format, err := image.Decode(file)
	if err != nil {
		sendErrorResponse(w, "Failed to decode image", err, http.StatusBadRequest)
		return
	}

	// Generate mosaic
	mosaicImg, err := generateMosaic(original, tileSize)
	if err != nil {
		sendErrorResponse(w, "Failed to generate mosaic", err, http.StatusInternalServerError)
		return
	}

	// Calculate duration
	duration := math.Round(time.Since(t0).Seconds()*100) / 100

	// Send response
	response := struct {
		MosaicImg string  `json:"mosaicImg"`
		Duration  float64 `json:"duration"`
		Format    string  `json:"format"`
	}{
		MosaicImg: mosaicImg,
		Duration:  duration,
		Format:    format,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// generateMosaic creates a mosaic from the original image using tiles from the database
func generateMosaic(original image.Image, tileSize int) (string, error) {
	bounds := original.Bounds()

	// Create new image for the mosaic
	newImage := image.NewNRGBA(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))

	// Clone tiles database to avoid concurrent access issues
	db := tiles_db.CloneTilesDB(tilesDB)

	// Source point for drawing
	sourcePoint := image.Point{0, 0}

	// Process image tile by tile - using original algorithm with single pixel sampling
	for y := bounds.Min.Y; y < bounds.Max.Y; y += tileSize {
		for x := bounds.Min.X; x < bounds.Max.X; x += tileSize {
			// Get color from original image at this position (single pixel sampling)
			r, g, b, _ := original.At(x, y).RGBA()
			color := [3]float64{float64(r), float64(g), float64(b)}

			// Find nearest tile by color
			nearestFileByColor := imgpkg.Nearest(color, &db)

			// If no tile found (database empty), refill it
			if nearestFileByColor == "" {
				db = tiles_db.CloneTilesDB(tilesDB)
				if len(db) > 0 {
					nearestFileByColor = imgpkg.Nearest(color, &db)
				}
			}

			// Process the tile
			if err := processTile(nearestFileByColor, newImage, x, y, tileSize, sourcePoint); err != nil {
				logrus.WithError(err).WithField("tile", nearestFileByColor).Warn("Failed to process tile")
			}
		}
	}

	// Encode the mosaic image to base64
	return encodeImageToBase64(newImage)
}

// processTile processes a single tile and draws it onto the mosaic
func processTile(tilePath string, newImage *image.NRGBA, x, y, tileSize int, sourcePoint image.Point) error {
	if tilePath == "" {
		// If no tile found, fill with black
		for py := y; py < y+tileSize; py++ {
			for px := x; px < x+tileSize; px++ {
				newImage.Set(px, py, image.Black)
			}
		}
		return nil
	}

	file, err := os.Open(tilePath)
	if err != nil {
		return fmt.Errorf("failed to open tile file: %w", err)
	}
	defer file.Close()

	tileImg, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode tile image: %w", err)
	}

	// Resize tile to match tile size
	resizedTile := imgpkg.Resize(tileImg, tileSize)
	tile := resizedTile.SubImage(resizedTile.Bounds())

	// Define tile bounds
	tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)

	// Draw tile onto mosaic
	draw.Draw(newImage, tileBounds, tile, sourcePoint, draw.Src)

	return nil
}

// encodeImageToBase64 encodes an image to base64 string
func encodeImageToBase64(img image.Image) (string, error) {
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// sendErrorResponse sends a JSON error response
func sendErrorResponse(w http.ResponseWriter, message string, err error, statusCode int) {
	logrus.WithError(err).Error(message)

	response := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
