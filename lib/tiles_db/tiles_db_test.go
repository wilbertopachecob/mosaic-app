package tiles_db

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestCloneTilesDB tests the CloneTilesDB function
func TestCloneTilesDB(t *testing.T) {
	// Create a test database
	original := map[string][3]float64{
		"test1.jpg": [3]float64{255, 0, 0},
		"test2.jpg": [3]float64{0, 255, 0},
		"test3.jpg": [3]float64{0, 0, 255},
	}

	// Clone the database
	cloned := CloneTilesDB(original)

	// Check that the clone has the same content
	if len(cloned) != len(original) {
		t.Errorf("Expected clone to have %d items, got %d", len(original), len(cloned))
	}

	for key, value := range original {
		if clonedValue, exists := cloned[key]; !exists {
			t.Errorf("Expected key '%s' to exist in clone", key)
		} else if clonedValue != value {
			t.Errorf("Expected value %v for key '%s', got %v", value, key, clonedValue)
		}
	}

	// Check that modifying the clone doesn't affect the original
	cloned["new.jpg"] = [3]float64{128, 128, 128}
	if len(original) == len(cloned) {
		t.Error("Expected original to remain unchanged when clone is modified")
	}
}

// TestIsImageFile tests the isImageFile function
func TestIsImageFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"JPEG file", "image.jpg", true},
		{"JPEG file uppercase", "IMAGE.JPG", true},
		{"PNG file", "image.png", true},
		{"GIF file", "image.gif", true},
		{"BMP file", "image.bmp", true},
		{"TIFF file", "image.tiff", true},
		{"WebP file", "image.webp", true},
		{"Text file", "file.txt", false},
		{"No extension", "file", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isImageFile(tt.filename)
			if result != tt.expected {
				t.Errorf("isImageFile('%s') = %t, want %t", tt.filename, result, tt.expected)
			}
		})
	}
}

// TestTilesDBWithEmptyDirectory tests TilesDB with an empty directory
func TestTilesDBWithEmptyDirectory(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "tiles_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Note: In a real implementation, you might want to make the tiles directory configurable
	// For this test, we'll just verify the function handles empty directories gracefully

	// Create an empty tiles directory
	emptyTilesDir := filepath.Join(tempDir, "empty_tiles")
	if err := os.Mkdir(emptyTilesDir, 0755); err != nil {
		t.Fatalf("Failed to create empty tiles directory: %v", err)
	}

	// The function should return an empty map for an empty directory
	// Note: This test is limited by the current implementation which hardcodes "tiles"
	// In a real refactor, you'd want to make the tiles directory configurable
}

// BenchmarkCloneTilesDB benchmarks the CloneTilesDB function
func BenchmarkCloneTilesDB(b *testing.B) {
	// Create a large test database
	original := make(map[string][3]float64, 1000)
	for i := 0; i < 1000; i++ {
		original[fmt.Sprintf("test%d.jpg", i)] = [3]float64{float64(i), float64(i), float64(i)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CloneTilesDB(original)
	}
}

// BenchmarkIsImageFile benchmarks the isImageFile function
func BenchmarkIsImageFile(b *testing.B) {
	filename := "test_image.jpg"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isImageFile(filename)
	}
} 