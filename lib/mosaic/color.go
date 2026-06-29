package mosaic

import (
	"image"
	"math"
)

// RGB represents a color as a 3-tuple of float64 values (R, G, B).
// Uses float64 for precise Euclidean distance calculations during tile matching.
type RGB [3]float64

// AverageColor computes the average RGB of an image by summing all pixel values
// and dividing by the pixel count. Uses raw 16-bit RGBA values for accuracy
// (image.RGBA returns 0-65535, not 0-255).
func AverageColor(img image.Image) RGB {
	avg, _ := AverageColorAndVariance(img)
	return avg
}

// AverageColorAndVariance computes the average RGB and total channel variance
// for an image. Variance is used to distinguish real photo tiles from flat
// color swatches that technically match colors but make mosaics look synthetic.
func AverageColorAndVariance(img image.Image) (RGB, float64) {
	bounds := img.Bounds()
	var r, g, b, rSq, gSq, bSq float64
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			rv := float64(r1)
			gv := float64(g1)
			bv := float64(b1)
			r += rv
			g += gv
			b += bv
			rSq += rv * rv
			gSq += gv * gv
			bSq += bv * bv
		}
	}
	n := float64(bounds.Dx() * bounds.Dy())
	if n <= 0 {
		return RGB{0, 0, 0}, 0
	}
	avg := RGB{r / n, g / n, b / n}
	variance := (rSq / n) - (avg[0] * avg[0]) +
		(gSq / n) - (avg[1] * avg[1]) +
		(bSq / n) - (avg[2] * avg[2])
	return avg, variance
}

// TopLeftPixelColor returns the RGB of the top-left pixel of a region.
// Matches the Sau Sheong article: "pick the top left pixel and assume that's the average color".
func TopLeftPixelColor(img image.Image, x, y int) RGB {
	r1, g1, b1, _ := img.At(x, y).RGBA()
	return RGB{float64(r1), float64(g1), float64(b1)}
}

// RegionAverageColor computes the average RGB of a rectangular region.
// Used to determine which tile best matches each region of the target image.
func RegionAverageColor(img image.Image, minX, minY, maxX, maxY int) RGB {
	var r, g, b float64
	n := 0
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r += float64(r1)
			g += float64(g1)
			b += float64(b1)
			n++
		}
	}
	if n <= 0 {
		return RGB{0, 0, 0}
	}
	nf := float64(n)
	return RGB{r / nf, g / nf, b / nf}
}

// Distance returns the Euclidean distance between two colors in RGB space.
// Smaller distance means better color match for tile selection.
func Distance(a, b RGB) float64 {
	return math.Sqrt(sq(b[0]-a[0]) + sq(b[1]-a[1]) + sq(b[2]-a[2]))
}

func sq(x float64) float64 { return x * x }
