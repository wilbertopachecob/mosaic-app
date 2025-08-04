package img

import (
	"image"
	"image/color"
	"math"
)

// AverageColor calculates the average RGB color of an image
// Returns a [3]float64 array representing [R, G, B] values
func AverageColor(img image.Image) [3]float64 {
	bounds := img.Bounds()
	r, g, b := 0.0, 0.0, 0.0
	
	// Sum all pixel values
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r += float64(r1)
			g += float64(g1)
			b += float64(b1)
		}
	}
	
	// Calculate average
	totalPixels := float64(bounds.Max.X * bounds.Max.Y)
	return [3]float64{r / totalPixels, g / totalPixels, b / totalPixels}
}

// Resize resizes an image to a new width while maintaining aspect ratio
// Returns a new NRGBA image with the specified width
func Resize(in image.Image, newWidth int) image.NRGBA {
	bounds := in.Bounds()
	
	// Calculate resize ratio
	ratio := bounds.Dx() / newWidth
	if ratio <= 0 {
		ratio = 1 // Prevent division by zero
	}
	
	// Calculate new dimensions
	newHeight := bounds.Dy() / ratio
	out := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Resize by sampling pixels
	for y, j := bounds.Min.Y, 0; y < bounds.Max.Y && j < newHeight; y, j = y+ratio, j+1 {
		for x, i := bounds.Min.X, 0; x < bounds.Max.X && i < newWidth; x, i = x+ratio, i+1 {
			r, g, b, a := in.At(x, y).RGBA()
			out.SetNRGBA(i, j, color.NRGBA{
				uint8(r >> 8),
				uint8(g >> 8),
				uint8(b >> 8),
				uint8(a >> 8),
			})
		}
	}
	
	return *out
}

// Nearest finds the tile with the closest color match to the target color
// Removes the found tile from the database to avoid reuse
// Returns the filename of the nearest matching tile
func Nearest(target [3]float64, db *map[string][3]float64) string {
	var filename string
	smallest := math.MaxFloat64
	
	for k, v := range *db {
		dist := Distance(target, v)
		if dist < smallest {
			filename, smallest = k, dist
		}
	}
	
	// Remove the selected tile from database to avoid reuse
	if filename != "" {
		delete(*db, filename)
	}
	
	return filename
}

// Distance calculates the Euclidean distance between two RGB color points
// Returns the distance as a float64
func Distance(p1 [3]float64, p2 [3]float64) float64 {
	return math.Sqrt(Sq(p2[0]-p1[0]) + Sq(p2[1]-p1[1]) + Sq(p2[2]-p1[2]))
}

// Sq calculates the square of a number
// Helper function for distance calculations
func Sq(n float64) float64 {
	return n * n
}
