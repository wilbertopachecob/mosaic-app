package img

import (
	"image"
	"image/color"
	"testing"
)

// TestAverageColor tests the AverageColor function
func TestAverageColor(t *testing.T) {
	// Create a simple test image (2x2 red pixels)
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	red := color.RGBA{255, 0, 0, 255}
	img.Set(0, 0, red)
	img.Set(0, 1, red)
	img.Set(1, 0, red)
	img.Set(1, 1, red)

	avg := AverageColor(img)

	// Check that the average is red (high R, low G and B)
	if avg[0] < 25000 { // R should be high
		t.Errorf("Expected high red value, got %f", avg[0])
	}
	if avg[1] > 1000 { // G should be low
		t.Errorf("Expected low green value, got %f", avg[1])
	}
	if avg[2] > 1000 { // B should be low
		t.Errorf("Expected low blue value, got %f", avg[2])
	}
}

// TestDistance tests the Distance function
func TestDistance(t *testing.T) {
	tests := []struct {
		name     string
		p1       [3]float64
		p2       [3]float64
		expected float64
	}{
		{
			name:     "Same point",
			p1:       [3]float64{0, 0, 0},
			p2:       [3]float64{0, 0, 0},
			expected: 0,
		},
		{
			name:     "Distance 1",
			p1:       [3]float64{0, 0, 0},
			p2:       [3]float64{1, 0, 0},
			expected: 1,
		},
		{
			name:     "Distance sqrt(3)",
			p1:       [3]float64{0, 0, 0},
			p2:       [3]float64{1, 1, 1},
			expected: 1.7320508075688772, // sqrt(3)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Distance(tt.p1, tt.p2)
			if result != tt.expected {
				t.Errorf("Distance(%v, %v) = %f, want %f", tt.p1, tt.p2, result, tt.expected)
			}
		})
	}
}

// TestSq tests the Sq function
func TestSq(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{1, 1},
		{2, 4},
		{-2, 4},
		{3.5, 12.25},
	}

	for _, tt := range tests {
		result := Sq(tt.input)
		if result != tt.expected {
			t.Errorf("Sq(%f) = %f, want %f", tt.input, result, tt.expected)
		}
	}
}

// TestNearest tests the Nearest function
func TestNearest(t *testing.T) {
	// Create a test database
	db := map[string][3]float64{
		"red.jpg":   [3]float64{255, 0, 0},
		"green.jpg": [3]float64{0, 255, 0},
		"blue.jpg":  [3]float64{0, 0, 255},
	}

	// Test finding nearest to red
	target := [3]float64{250, 10, 10}
	nearest := Nearest(target, &db)

	if nearest != "red.jpg" {
		t.Errorf("Expected 'red.jpg', got '%s'", nearest)
	}

	// Check that the found tile was removed from the database
	if _, exists := db["red.jpg"]; exists {
		t.Error("Expected 'red.jpg' to be removed from database")
	}

	// Check that other tiles are still in the database
	if _, exists := db["green.jpg"]; !exists {
		t.Error("Expected 'green.jpg' to remain in database")
	}
	if _, exists := db["blue.jpg"]; !exists {
		t.Error("Expected 'blue.jpg' to remain in database")
	}
}

// TestResize tests the Resize function
func TestResize(t *testing.T) {
	// Create a test image (4x4)
	original := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			original.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	// Resize to 2x2
	resized := Resize(original, 2)

	// Check dimensions
	bounds := resized.Bounds()
	if bounds.Dx() != 2 || bounds.Dy() != 2 {
		t.Errorf("Expected 2x2 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// Check that the image is still red
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			if r < 25000 || g > 1000 || b > 1000 {
				t.Errorf("Expected red pixel at (%d, %d), got R:%d G:%d B:%d", x, y, r, g, b)
			}
		}
	}
}

// BenchmarkAverageColor benchmarks the AverageColor function
func BenchmarkAverageColor(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AverageColor(img)
	}
}

// BenchmarkDistance benchmarks the Distance function
func BenchmarkDistance(b *testing.B) {
	p1 := [3]float64{0, 0, 0}
	p2 := [3]float64{255, 255, 255}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Distance(p1, p2)
	}
} 