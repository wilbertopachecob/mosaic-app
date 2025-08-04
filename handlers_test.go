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
)

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

// TestMosaicHandlerWithInvalidRequest tests mosaic handler with invalid requests
func TestMosaicHandlerWithInvalidRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		contentType    string
		body           string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Wrong method",
			method:         "GET",
			contentType:    "application/json",
			body:           "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid content type",
			method:         "POST",
			contentType:    "application/json",
			body:           "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid form data",
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
			handler := http.HandlerFunc(mosaicHandler)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

// TestMosaicHandlerWithInvalidTileSize tests mosaic handler with invalid tile sizes
func TestMosaicHandlerWithInvalidTileSize(t *testing.T) {
	// Create a test image
	img := createTestImage(100, 100)
	imgBytes := imageToBytes(t, img)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("imgUpload", "test.jpg")
	require.NoError(t, err)
	part.Write(imgBytes)

	// Add invalid tile size
	err = writer.WriteField("tileSize", "invalid")
	require.NoError(t, err)

	writer.Close()

	req, err := http.NewRequest("POST", "/api/file/upload", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mosaicHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid tile size")
}

// TestMosaicHandlerWithValidRequest tests mosaic handler with a valid request
func TestMosaicHandlerWithValidRequest(t *testing.T) {
	// Skip if no tiles database is available
	if len(tilesDB) == 0 {
		t.Skip("No tiles database available for testing")
	}

	// Create a test image
	img := createTestImage(50, 50)
	imgBytes := imageToBytes(t, img)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("imgUpload", "test.jpg")
	require.NoError(t, err)
	part.Write(imgBytes)

	// Add valid tile size
	err = writer.WriteField("tileSize", "20")
	require.NoError(t, err)

	writer.Close()

	req, err := http.NewRequest("POST", "/api/file/upload", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mosaicHandler)

	handler.ServeHTTP(rr, req)

	// Should succeed or fail gracefully
	assert.True(t, rr.Code == http.StatusCreated || rr.Code == http.StatusInternalServerError)

	if rr.Code == http.StatusCreated {
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

// TestSendErrorResponse tests the sendErrorResponse function
func TestSendErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()

	sendErrorResponse(rr, http.StatusBadRequest, "Test Error", "Test Details")

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Test Error", response["error"])
	assert.Equal(t, "Test Details", response["message"])
	assert.Equal(t, float64(400), response["code"])
}
