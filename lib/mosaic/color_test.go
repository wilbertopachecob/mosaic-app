package mosaic

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func TestAverageColor(t *testing.T) {
	// Create a 2x2 red image
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	red := color.RGBA{255, 0, 0, 255}
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			img.Set(x, y, red)
		}
	}

	avg := AverageColor(img)

	// RGBA returns 16-bit values; red=255 becomes 255*257 = 65535 in high byte
	// Average should be high R, low G and B
	if avg[0] < 25000 {
		t.Errorf("Expected high red value, got %f", avg[0])
	}
	if avg[1] > 1000 {
		t.Errorf("Expected low green value, got %f", avg[1])
	}
	if avg[2] > 1000 {
		t.Errorf("Expected low blue value, got %f", avg[2])
	}
}

func TestRegionAverageColor(t *testing.T) {
	// 4x4 image, check 2x2 region
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	green := color.RGBA{0, 255, 0, 255}
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			img.Set(x, y, green)
		}
	}

	avg := RegionAverageColor(img, 0, 0, 2, 2)
	if avg[1] < 25000 {
		t.Errorf("Expected high green in region, got %f", avg[1])
	}
}

func TestRegionAverageColor_EmptyRegion(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	avg := RegionAverageColor(img, 2, 2, 2, 2) // zero-size region
	if avg[0] != 0 || avg[1] != 0 || avg[2] != 0 {
		t.Errorf("Expected zero color for empty region, got %v", avg)
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		name     string
		a        RGB
		b        RGB
		expected float64
	}{
		{"same point", RGB{0, 0, 0}, RGB{0, 0, 0}, 0},
		{"unit distance", RGB{0, 0, 0}, RGB{1, 0, 0}, 1},
		{"sqrt(3)", RGB{0, 0, 0}, RGB{1, 1, 1}, math.Sqrt(3)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Distance(tt.a, tt.b)
			if math.Abs(got-tt.expected) > 1e-9 {
				t.Errorf("Distance(%v, %v) = %f, want %f", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestDistance_Symmetry(t *testing.T) {
	a, b := RGB{1, 2, 3}, RGB{4, 5, 6}
	if Distance(a, b) != Distance(b, a) {
		t.Error("Distance should be symmetric")
	}
}
