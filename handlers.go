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

	"github.com/sirupsen/logrus"
	"wilbertopachecob/mosaic/lib/img"
	"wilbertopachecob/mosaic/lib/tiles_db"
	"wilbertopachecob/mosaic/models"
)

// mosaicHandler handles the mosaic generation request
// It processes an uploaded image and creates a mosaic using tiles from the database
func mosaicHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	
	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Parse multipart form with 10MB limit
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		logrus.WithError(err).Error("Failed to parse multipart form")
		sendErrorResponse(w, http.StatusBadRequest, "Invalid form data", err.Error())
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("imgUpload")
	if err != nil {
		logrus.WithError(err).Error("Failed to get uploaded file")
		sendErrorResponse(w, http.StatusBadRequest, "No file uploaded", err.Error())
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > 10<<20 { // 10MB limit
		sendErrorResponse(w, http.StatusBadRequest, "File too large", "File size exceeds 10MB limit")
		return
	}

	// Get tile size parameter
	tileSizeStr := r.FormValue("tileSize")
	tileSize, err := strconv.Atoi(tileSizeStr)
	if err != nil || tileSize <= 0 {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid tile size", "Tile size must be a positive integer")
		return
	}

	// Validate tile size range
	if tileSize < 5 || tileSize > 200 {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid tile size", "Tile size must be between 5 and 200 pixels")
		return
	}

	// Decode the original image
	original, format, err := image.Decode(file)
	if err != nil {
		logrus.WithError(err).Error("Failed to decode image")
		sendErrorResponse(w, http.StatusBadRequest, "Invalid image format", err.Error())
		return
	}

	logrus.WithFields(logrus.Fields{
		"format":    format,
		"tileSize":  tileSize,
		"fileSize":  header.Size,
		"fileName":  header.Filename,
	}).Info("Processing mosaic request")

	// Generate mosaic
	mosaicImg, err := generateMosaic(original, tileSize)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate mosaic")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to generate mosaic", err.Error())
		return
	}

	// Calculate processing time
	duration := math.Round(time.Since(startTime).Seconds()*100) / 100

	// Create response
	response := models.MosaicResponse{
		MosaicImg: mosaicImg,
		Duration:  duration,
	}

	// Send success response
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.WithError(err).Error("Failed to encode response")
	}
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

	// Process image tile by tile
	for y := bounds.Min.Y; y < bounds.Max.Y; y += tileSize {
		for x := bounds.Min.X; x < bounds.Max.X; x += tileSize {
			// Get color from original image at this position
			r, g, b, _ := original.At(x, y).RGBA()
			color := [3]float64{float64(r), float64(g), float64(b)}
			
			// Find nearest tile by color
			nearestFileByColor := img.Nearest(color, &db)
			
			// Open and process the tile
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
	resizedTile := img.Resize(tileImg, tileSize)
	tile := resizedTile.SubImage(resizedTile.Bounds())
	
	// Define tile bounds
	tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
	
	// Draw tile onto mosaic
	draw.Draw(newImage, tileBounds, tile, sourcePoint, draw.Src)
	
	return nil
}

// encodeImageToBase64 encodes an image to base64 string
func encodeImageToBase64(img image.Image) (string, error) {
	var buf bytes.Buffer
	
	// Encode as JPEG with quality 90
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}
	
	// Convert to base64
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// sendErrorResponse sends a standardized error response
func sendErrorResponse(w http.ResponseWriter, statusCode int, message, details string) {
	response := models.ErrorResponse{
		Error:   message,
		Message: details,
		Code:    statusCode,
	}
	
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.WithError(err).Error("Failed to encode error response")
	}
}
