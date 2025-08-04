package img

import (
	"image"
	"image/color"
	"math"
)

// AverageColor calculates the average color of an image
func AverageColor(img image.Image) [3]float64 {
	bounds := img.Bounds()
	r, g, b := 0.0, 0.0, 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for i := bounds.Min.X; i < bounds.Max.X; i++ {
			r1, g1, b1, _ := img.At(i, y).RGBA()
			r, g, b = r+float64(r1), g+float64(g1), b+float64(b1)
		}
	}
	totalPixels := float64(bounds.Max.X * bounds.Max.Y)
	return [3]float64{r / totalPixels, g / totalPixels, b / totalPixels}
}

// Resize resizes an image to a new width while maintaining aspect ratio
// This matches the original working implementation from the MVP commit
func Resize(in image.Image, newWidth int) image.NRGBA {
	bounds := in.Bounds()
	
	// Handle edge cases to prevent division by zero
	if newWidth <= 0 {
		newWidth = 1
	}
	if bounds.Dx() <= 0 {
		return *image.NewNRGBA(image.Rect(0, 0, 1, 1))
	}
	
	ratio := bounds.Dx() / newWidth
	if ratio <= 0 {
		ratio = 1
	}
	
	// Calculate new dimensions
	newHeight := bounds.Dy() / ratio
	if newHeight <= 0 {
		newHeight = 1
	}
	
	out := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Resize by sampling pixels - this is the key difference from the broken version
	for y, j := bounds.Min.Y, 0; y < bounds.Max.Y && j < newHeight; y, j = y+ratio, j+1 {
		for x, i := bounds.Min.X, 0; x < bounds.Max.X && i < newWidth; x, i = x+ratio, i+1 {
			r, g, b, a := in.At(x, y).RGBA()
			out.SetNRGBA(i, j, color.NRGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}
	return *out
}

// Nearest finds the nearest color match in the database and removes it
func Nearest(target [3]float64, db *map[string][3]float64) string {
	var filename string
	smallest := 1000000.0
	for k, v := range *db {
		dist := Distance(target, v)
		if dist < smallest {
			filename, smallest = k, dist
		}
	}
	delete(*db, filename)
	return filename
}

// Distance calculates the Euclidean distance between two color points
func Distance(p1 [3]float64, p2 [3]float64) float64 {
	return math.Sqrt(Sq(p2[0]-p1[0]) + Sq(p2[1]-p1[1]) + Sq(p2[2]-p1[2]))
}

// Sq calculates the square of a number
func Sq(n float64) float64 {
	return n * n
}
