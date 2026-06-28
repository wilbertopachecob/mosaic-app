package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"strconv"
	"time"

	"wilbertopachecob/mosaic/lib/mosaic"
)

// Response represents the API response
type Response struct {
	MosaicImg string  `json:"mosaicImg"`
	Duration  float64 `json:"duration"`
	Format    string  `json:"format"`
}

// generateMosaic creates a mosaic from the original image and returns base64-encoded JPEG.
func generateMosaic(original image.Image, tileSize int, blend float64) (string, error) {
	mosaicImg, err := mosaicGenerator.GenerateWithOptions(original, tileSize, mosaic.Options{
		SourceBlend: blend,
	})
	if err != nil {
		return "", fmt.Errorf("mosaic generation: %w", err)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, mosaicImg, &jpeg.Options{Quality: 90}); err != nil {
		return "", fmt.Errorf("failed to encode mosaic: %v", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// uploadHandler handles file upload and mosaic generation
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// #region agent log
	mosaic.DebugLog("H0", "uploadHandler entered", map[string]interface{}{"method": r.Method})
	// #endregion
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, _, err := r.FormFile("imgUpload")
	if err != nil {
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get tile size
	tileSizeStr := r.FormValue("tileSize")
	tileSize, err := strconv.Atoi(tileSizeStr)
	if err != nil {
		tileSize = 20 // Default tile size
	}

	blend := mosaic.DefaultOptions().SourceBlend
	if blendStr := r.FormValue("blend"); blendStr != "" {
		parsedBlend, err := strconv.ParseFloat(blendStr, 64)
		if err == nil {
			blend = parsedBlend
		}
	}

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Failed to decode image", http.StatusBadRequest)
		return
	}

	// Generate mosaic
	// #region agent log
	mosaic.DebugLog("H0", "before generateMosaic", map[string]interface{}{"tileSize": tileSize, "blend": blend, "imgBounds": img.Bounds().String()})
	// #endregion
	start := time.Now()
	mosaicBase64, err := generateMosaic(img, tileSize, blend)
	duration := time.Since(start).Seconds()
	if err != nil {
		// #region agent log
		mosaic.DebugLog("H0", "generateMosaic error", map[string]interface{}{"error": err.Error()})
		// #endregion
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Failed to generate mosaic",
			"message": err.Error(),
		})
		return
	}

	// #region agent log
	mosaic.DebugLog("H0", "generateMosaic success", map[string]interface{}{"base64Len": len(mosaicBase64)})
	// #endregion
	response := Response{
		MosaicImg: mosaicBase64,
		Duration:  duration,
		Format:    "jpeg",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// healthHandler provides health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "healthy",
		"service": "mosaic-app",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
