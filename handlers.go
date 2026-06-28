package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"time"

	"wilbertopachecob/mosaic/config"
	"wilbertopachecob/mosaic/lib/mosaic"
)

// Response represents the API response
type Response struct {
	MosaicImg string  `json:"mosaicImg"`
	Duration  float64 `json:"duration"`
	Format    string  `json:"format"`
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

var (
	appConfig       *config.Config
	mosaicGenerator *mosaic.Generator
)

func configureServer(cfg *config.Config, generator *mosaic.Generator) {
	appConfig = cfg
	mosaicGenerator = generator
}

func writeJSONError(w http.ResponseWriter, status int, errKey, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{
		Error:   errKey,
		Message: message,
	})
}

func serverErrorMessage(err error) string {
	log.Printf("upload error: %v", err)
	if appConfig != nil && appConfig.Production {
		return "An internal error occurred"
	}
	return err.Error()
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cfg := appConfig
	if cfg == nil {
		cfg = config.Load()
	}

	maxMemory := cfg.MaxFileSize + (1 << 20) // file size plus ~1MB for form fields
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "Failed to parse upload")
		return
	}

	file, header, err := r.FormFile("imgUpload")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_request", "Missing image upload")
		return
	}
	defer file.Close()

	if header.Size > cfg.MaxFileSize {
		writeJSONError(w, http.StatusBadRequest, "file_too_large", fmt.Sprintf("File exceeds maximum size of %d bytes", cfg.MaxFileSize))
		return
	}

	tileSize, err := parseTileSize(r.FormValue("tileSize"), cfg)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_tile_size", err.Error())
		return
	}

	blend, err := parseBlend(r.FormValue("blend"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_blend", err.Error())
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_image", "Failed to decode image")
		return
	}

	if err := validateImageDimensions(img, cfg); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_image", err.Error())
		return
	}

	start := time.Now()
	mosaicBase64, err := generateMosaic(img, tileSize, blend)
	duration := time.Since(start).Seconds()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "generation_failed", serverErrorMessage(err))
		return
	}

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
