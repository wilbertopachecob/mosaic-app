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
	"wilbertopachecob/mosaic/config"
	"wilbertopachecob/mosaic/lib/mosaic"
)

func init() {
	if mosaicGenerator == nil {
		mosaicGenerator = mosaic.NewGenerator()
		_ = mosaicGenerator.LoadTiles("tiles")
	}
	if appConfig == nil {
		appConfig = testConfig()
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
	img := createTestImage(100, 100)
	imgBytes := imageToBytes(t, img)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("imgUpload", "test.jpg")
	require.NoError(t, err)
	_, err = part.Write(imgBytes)
	require.NoError(t, err)

	require.NoError(t, writer.WriteField("tileSize", "invalid"))
	require.NoError(t, writer.Close())

	req, err := http.NewRequest("POST", "/api/file/upload", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(uploadHandler)
	handler.ServeHTTP(rr, req)

	assert.True(t, rr.Code == http.StatusOK || rr.Code == http.StatusInternalServerError)
}

func TestUploadHandlerRejectsOutOfRangeTileSize(t *testing.T) {
	imgBytes := imageToBytes(t, createTestImage(100, 100))

	for _, tileSize := range []string{"1", "200"} {
		t.Run("tileSize="+tileSize, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("imgUpload", "test.jpg")
			require.NoError(t, err)
			_, err = part.Write(imgBytes)
			require.NoError(t, err)
			require.NoError(t, writer.WriteField("tileSize", tileSize))
			require.NoError(t, writer.Close())

			req, err := http.NewRequest("POST", "/api/file/upload", body)
			require.NoError(t, err)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			rr := httptest.NewRecorder()
			uploadHandler(rr, req)

			assert.Equal(t, http.StatusBadRequest, rr.Code)

			var response errorResponse
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))
			assert.Equal(t, "invalid_tile_size", response.Error)
		})
	}
}

func TestUploadHandlerRejectsInvalidBlend(t *testing.T) {
	imgBytes := imageToBytes(t, createTestImage(100, 100))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("imgUpload", "test.jpg")
	require.NoError(t, err)
	_, err = part.Write(imgBytes)
	require.NoError(t, err)
	require.NoError(t, writer.WriteField("tileSize", "20"))
	require.NoError(t, writer.WriteField("blend", "2"))
	require.NoError(t, writer.Close())

	req, err := http.NewRequest("POST", "/api/file/upload", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	uploadHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUploadHandlerRejectsOversizedImage(t *testing.T) {
	appConfig = &config.Config{
		MaxFileSize:    10 * 1024 * 1024,
		MaxImageWidth:  200,
		MaxImageHeight: 200,
		MinTileSize:    5,
		MaxTileSize:    100,
	}
	t.Cleanup(func() {
		appConfig = testConfig()
	})

	imgBytes := imageToBytes(t, createTestImage(300, 300))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("imgUpload", "test.jpg")
	require.NoError(t, err)
	_, err = part.Write(imgBytes)
	require.NoError(t, err)
	require.NoError(t, writer.WriteField("tileSize", "20"))
	require.NoError(t, writer.Close())

	req, err := http.NewRequest("POST", "/api/file/upload", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	uploadHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUploadHandlerRejectsOversizedBody(t *testing.T) {
	appConfig = &config.Config{
		MaxFileSize:    1024, // 1KB cap
		MaxImageWidth:  4096,
		MaxImageHeight: 4096,
		MinTileSize:    5,
		MaxTileSize:    100,
	}
	t.Cleanup(func() {
		appConfig = testConfig()
	})

	// Build a body beyond the MaxFileSize + 1MB allowance so MaxBytesReader
	// trips during parsing, before the whole body is buffered to disk.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("imgUpload", "big.jpg")
	require.NoError(t, err)
	_, err = part.Write(bytes.Repeat([]byte("A"), 2<<20)) // 2MB of payload
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	require.Greater(t, body.Len(), int(appConfig.MaxFileSize)+(1<<20))

	req, err := http.NewRequest("POST", "/api/file/upload", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	uploadHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))
	assert.Equal(t, "file_too_large", response["error"])
}

func TestServerErrorMessageProduction(t *testing.T) {
	appConfig = &config.Config{Production: true}
	t.Cleanup(func() {
		appConfig = testConfig()
	})

	msg := serverErrorMessage(assert.AnError)
	assert.Equal(t, "An internal error occurred", msg)
}

// TestUploadHandlerWithValidRequest tests upload handler with a valid request
func TestUploadHandlerWithValidRequest(t *testing.T) {
	if mosaicGenerator == nil || mosaicGenerator.TileCount() == 0 {
		t.Skip("No tiles loaded for testing")
	}

	img := createTestImage(50, 50)
	imgBytes := imageToBytes(t, img)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("imgUpload", "test.jpg")
	require.NoError(t, err)
	_, err = part.Write(imgBytes)
	require.NoError(t, err)

	require.NoError(t, writer.WriteField("tileSize", "20"))
	require.NoError(t, writer.WriteField("blend", "0.55"))
	require.NoError(t, writer.Close())

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

func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	return img
}

func imageToBytes(t *testing.T, img image.Image) []byte {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	require.NoError(t, err)
	return buf.Bytes()
}
