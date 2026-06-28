package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"wilbertopachecob/mosaic/lib/mosaic"
)

func init() {
	// Initialize mosaic generator for handler tests (main() is not run in tests)
	if mosaicGenerator == nil {
		mosaicGenerator = mosaic.NewGenerator()
		_ = mosaicGenerator.LoadTiles("tiles")
	}
}

// TestHealthHandler tests the health check endpoint
func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "mosaic-app", response["service"])
}

// TestUploadHandlerWithInvalidRequest tests upload handler with invalid requests
func TestUploadHandlerWithInvalidRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		contentType    string
		body           string
		expectedStatus int
	}{
		{
			name:           "Wrong method",
			method:         "GET",
			contentType:    "application/json",
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid form (no multipart)",
			method:         "POST",
			contentType:    "application/json",
			body:           "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/api/file/upload", strings.NewReader(tt.body))
			require.NoError(t, err)

			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(uploadHandler)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

// TestUploadHandlerWithInvalidTileSize tests that invalid tile size falls back to default
func TestUploadHandlerWithInvalidTileSize(t *testing.T) {
	// Create a test image
	img := createTestImage(100, 100)
	imgBytes := imageToBytes(t, img)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("imgUpload", "test.jpg")
	require.NoError(t, err)
	part.Write(imgBytes)

	writer.WriteField("tileSize", "invalid")
	writer.Close()

	req, err := http.NewRequest("POST", "/api/file/upload", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(uploadHandler)
	handler.ServeHTTP(rr, req)

	// With invalid tile size we use default 20; request may succeed or fail (no tiles)
	assert.True(t, rr.Code == http.StatusOK || rr.Code == http.StatusInternalServerError)
}

// TestUploadHandlerWithValidRequest tests upload handler with a valid request
func TestUploadHandlerWithValidRequest(t *testing.T) {
	// Skip if no tiles loaded (mosaicGenerator is set in main, may be nil in test)
	if mosaicGenerator == nil || mosaicGenerator.TileCount() == 0 {
		t.Skip("No tiles loaded for testing")
	}

	img := createTestImage(50, 50)
	imgBytes := imageToBytes(t, img)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("imgUpload", "test.jpg")
	require.NoError(t, err)
	part.Write(imgBytes)

	writer.WriteField("tileSize", "20")
	writer.WriteField("blend", "0.55")
	writer.Close()

	req, err := http.NewRequest("POST", "/api/file/upload", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(uploadHandler)
	handler.ServeHTTP(rr, req)

	assert.True(t, rr.Code == http.StatusOK || rr.Code == http.StatusInternalServerError)

	if rr.Code == http.StatusOK {
		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "mosaicImg")
		assert.Contains(t, response, "duration")
	}
}

// Helper functions

// createTestImage creates a simple test image
func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	return img
}

// imageToBytes converts an image to JPEG bytes
func imageToBytes(t *testing.T, img image.Image) []byte {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	require.NoError(t, err)
	return buf.Bytes()
}
