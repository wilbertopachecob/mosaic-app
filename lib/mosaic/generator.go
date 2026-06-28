package mosaic

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg" // Register JPEG decoder for image.Decode
	_ "image/png"  // Register PNG decoder for image.Decode
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Generator creates photo mosaics from a target image and a set of tile images.
// Tiles are matched by average color using Euclidean distance in RGB space.
type Generator struct {
	mu    sync.RWMutex
	tiles []tileEntry
}

// Options controls the visual style of generated mosaics.
type Options struct {
	// SourceBlend controls how much of the original image is blended over the
	// tile mosaic. Services such as EasyMoza use this kind of compositing to
	// preserve faces, edges, and shadows while keeping tiles visible up close.
	// Use 0 for pure tile output and 1 for the original image.
	SourceBlend float64
}

// DefaultOptions returns the production default mosaic style.
func DefaultOptions() Options {
	return Options{SourceBlend: 0.42}
}

// tileEntry holds a preloaded tile image and its average color for fast matching.
type tileEntry struct {
	path     string
	img      image.Image
	color    RGB
	variance float64
}

// NewGenerator returns a new mosaic generator with no tiles loaded.
func NewGenerator() *Generator {
	return &Generator{tiles: make([]tileEntry, 0)}
}

// LoadTiles scans dir for image files (.jpg, .jpeg, .png), loads each one,
// computes its average color, and adds it to the tile database.
// Returns an error if the directory cannot be read.
// Tiles that fail to decode are skipped (no error).
func (g *Generator) LoadTiles(dir string) error {
	var allTiles []tileEntry
	var photoTiles []tileEntry
	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			if path == dir {
				return err
			}
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return nil
		}

		img, err := loadImage(path)
		if err != nil {
			return nil // Skip unreadable files
		}

		avg, variance := AverageColorAndVariance(img)
		tile := tileEntry{
			path:     path,
			img:      img,
			color:    avg,
			variance: variance,
		}
		allTiles = append(allTiles, tile)

		// Prefer real photo tiles when available. The repository contains
		// color_* swatches for tests/demos, but those make production mosaics
		// look like flat posterization instead of photo mosaics.
		if !strings.HasPrefix(strings.ToLower(entry.Name()), "color_") && variance > 300 {
			photoTiles = append(photoTiles, tile)
		}
		return nil
	})
	if err != nil {
		return err
	}

	tiles := allTiles
	if len(photoTiles) >= 8 {
		tiles = photoTiles
	}

	g.mu.Lock()
	g.tiles = tiles
	g.mu.Unlock()
	return nil
}

// TileCount returns the number of loaded tiles. Returns 0 if none loaded.
func (g *Generator) TileCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.tiles)
}

// #region agent log
// DebugLog writes an NDJSON line to the session log file for debugging.
func DebugLog(hypothesisId, message string, data map[string]interface{}) {
	payload := map[string]interface{}{
		"hypothesisId": hypothesisId,
		"location":     "generator.go:Generate",
		"message":      message,
		"data":         data,
		"timestamp":    time.Now().UnixMilli(),
	}
	b, _ := json.Marshal(payload)
	logPath := "/Users/wilbertopachecobatista/Projects/mosaic-app/.cursor/debug-f1a3d5.log"
	os.MkdirAll(filepath.Dir(logPath), 0755)
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		f.Write(append(b, '\n'))
		f.Close()
	}
}

// #endregion

// Generate creates a photo mosaic from target by dividing it into tile-sized
// regions, finding the best-matching tile for each region, and drawing it.
// tileSize is the side length of each square tile in pixels.
// Returns an error if no tiles are loaded or generation fails.
func (g *Generator) Generate(target image.Image, tileSize int) (image.Image, error) {
	return g.GenerateWithOptions(target, tileSize, DefaultOptions())
}

// GenerateWithOptions creates a photo mosaic using explicit visual options.
func (g *Generator) GenerateWithOptions(target image.Image, tileSize int, opts Options) (image.Image, error) {
	g.mu.RLock()
	tiles := g.tiles
	g.mu.RUnlock()

	if len(tiles) == 0 {
		return nil, ErrNoTiles
	}

	// Clamp tile size to avoid degenerate regions
	if tileSize < 1 {
		tileSize = 1
	}
	opts.SourceBlend = clamp01(opts.SourceBlend)

	bounds := target.Bounds()
	out := image.NewNRGBA(bounds)
	// Pre-fill with white so undrawn pixels encode as white in JPEG (default NRGBA is transparent black → black in JPEG)
	draw.Draw(out, bounds, &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	// #region agent log
	DebugLog("H5", "Generate start", map[string]interface{}{
		"tileCount": len(tiles),
		"tileSize":  tileSize,
		"bounds":    map[string]int{"minX": bounds.Min.X, "minY": bounds.Min.Y, "maxX": bounds.Max.X, "maxY": bounds.Max.Y},
		"width":     bounds.Dx(),
		"height":    bounds.Dy(),
		"blend":     opts.SourceBlend,
	})
	// #endregion

	regionCount := 0
	var firstRect, lastRect image.Rectangle
	var firstResizedW, firstResizedH int
	var lastTilePath string
	var lastRegionColor RGB

	// Process each tile-sized region
	for y := bounds.Min.Y; y < bounds.Max.Y; y += tileSize {
		for x := bounds.Min.X; x < bounds.Max.X; x += tileSize {
			endX := x + tileSize
			endY := y + tileSize
			if endX > bounds.Max.X {
				endX = bounds.Max.X
			}
			if endY > bounds.Max.Y {
				endY = bounds.Max.Y
			}

			// Use region average for robust matching (top-left pixel can misrepresent varied regions)
			regionColor := RegionAverageColor(target, x, y, endX, endY)

			// Find best matching tile (nearest color in RGB space)
			best := findNearest(regionColor, tiles)

			// Resize tile to fit region and draw
			w, h := endX-x, endY-y
			resized := resizeTile(best.img, w, h)
			adjusted := colorCorrectTile(resized, best.color, regionColor)
			rect := image.Rect(x, y, endX, endY)
			// Author uses draw.Src to fully replace pixels; draw.Over can leave transparent gaps
			draw.Draw(out, rect, adjusted, image.Point{}, draw.Src)

			// #region agent log
			regionCount++
			if regionCount == 1 {
				firstRect = rect
				firstResizedW, firstResizedH = w, h
			}
			lastRect = rect
			lastTilePath = best.path
			lastRegionColor = regionColor
			// #endregion
		}
	}

	if opts.SourceBlend > 0 {
		blendWithTarget(out, target, opts.SourceBlend)
	}

	// #region agent log
	samplePixel := func(ox, oy int) map[string]interface{} {
		c := out.NRGBAAt(ox+bounds.Min.X, oy+bounds.Min.Y)
		return map[string]interface{}{"x": ox, "y": oy, "R": c.R, "G": c.G, "B": c.B, "A": c.A}
	}
	DebugLog("H2", "Generate end", map[string]interface{}{
		"regionCount":     regionCount,
		"firstRect":       map[string]int{"minX": firstRect.Min.X, "minY": firstRect.Min.Y, "maxX": firstRect.Max.X, "maxY": firstRect.Max.Y},
		"lastRect":        map[string]int{"minX": lastRect.Min.X, "minY": lastRect.Min.Y, "maxX": lastRect.Max.X, "maxY": lastRect.Max.Y},
		"firstResizedW":   firstResizedW,
		"firstResizedH":   firstResizedH,
		"lastTilePath":    filepath.Base(lastTilePath),
		"lastRegionColor": []float64{lastRegionColor[0], lastRegionColor[1], lastRegionColor[2]},
		"pixelTL":         samplePixel(0, 0),
		"pixelCenter":     samplePixel(bounds.Dx()/2, bounds.Dy()/2),
		"pixelBR":         samplePixel(bounds.Dx()-1, bounds.Dy()-1),
	})
	// #endregion

	return out, nil
}

// findNearest returns the tile with the smallest color distance to target.
// Does not modify the tile list (allows tile reuse for larger mosaics).
func findNearest(target RGB, tiles []tileEntry) tileEntry {
	best := tiles[0]
	minDist := Distance(target, best.color)

	for i := 1; i < len(tiles); i++ {
		d := Distance(target, tiles[i].color)
		if d < minDist {
			minDist = d
			best = tiles[i]
		}
	}
	return best
}

// loadImage decodes an image from path. Supports JPEG and PNG.
func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// image.Decode returns (Image, format, error)
	decoded, _, err := image.Decode(f)
	return decoded, err
}

// resizeTile scales img to width x height using bilinear sampling.
func resizeTile(img image.Image, width, height int) image.Image {
	// #region agent log
	if width <= 0 || height <= 0 {
		DebugLog("H3", "resizeTile invalid dims", map[string]interface{}{"width": width, "height": height})
		return img
	}
	// #endregion

	src := img.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, width, height))
	srcW, srcH := src.Dx(), src.Dy()
	if srcW <= 0 || srcH <= 0 {
		return out
	}

	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			srcX := scaleCoord(dx, width, srcW)
			srcY := scaleCoord(dy, height, srcH)
			out.SetNRGBA(dx, dy, bilinearAt(img, src, srcX, srcY))
		}
	}
	return out
}

func scaleCoord(dst, dstSize, srcSize int) float64 {
	if dstSize <= 1 || srcSize <= 1 {
		return 0
	}
	return float64(dst) * float64(srcSize-1) / float64(dstSize-1)
}

func bilinearAt(img image.Image, bounds image.Rectangle, x, y float64) color.NRGBA {
	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	x1 := min(x0+1, bounds.Dx()-1)
	y1 := min(y0+1, bounds.Dy()-1)
	fx := x - float64(x0)
	fy := y - float64(y0)

	c00 := nrgbaAt(img, bounds.Min.X+x0, bounds.Min.Y+y0)
	c10 := nrgbaAt(img, bounds.Min.X+x1, bounds.Min.Y+y0)
	c01 := nrgbaAt(img, bounds.Min.X+x0, bounds.Min.Y+y1)
	c11 := nrgbaAt(img, bounds.Min.X+x1, bounds.Min.Y+y1)

	return color.NRGBA{
		R: interpolateChannel(c00.R, c10.R, c01.R, c11.R, fx, fy),
		G: interpolateChannel(c00.G, c10.G, c01.G, c11.G, fx, fy),
		B: interpolateChannel(c00.B, c10.B, c01.B, c11.B, fx, fy),
		A: interpolateChannel(c00.A, c10.A, c01.A, c11.A, fx, fy),
	}
}

func interpolateChannel(c00, c10, c01, c11 uint8, fx, fy float64) uint8 {
	top := float64(c00)*(1-fx) + float64(c10)*fx
	bottom := float64(c01)*(1-fx) + float64(c11)*fx
	return clamp8(top*(1-fy) + bottom*fy)
}

func nrgbaAt(img image.Image, x, y int) color.NRGBA {
	r, g, b, a := img.At(x, y).RGBA()
	if a == 0 {
		return color.NRGBA{}
	}
	return color.NRGBA{
		R: uint8((r * 0xff / a) & 0xff),
		G: uint8((g * 0xff / a) & 0xff),
		B: uint8((b * 0xff / a) & 0xff),
		A: uint8(a >> 8),
	}
}

func colorCorrectTile(tile image.Image, tileColor, targetColor RGB) image.Image {
	bounds := tile.Bounds()
	out := image.NewNRGBA(bounds)
	const correctionStrength = 0.72
	const targetBlend = 0.18

	target := rgbTo8Bit(targetColor)
	tileAvg := rgbTo8Bit(tileColor)
	dr := (target[0] - tileAvg[0]) * correctionStrength
	dg := (target[1] - tileAvg[1]) * correctionStrength
	db := (target[2] - tileAvg[2]) * correctionStrength

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := nrgbaAt(tile, x, y)
			r := float64(c.R) + dr
			g := float64(c.G) + dg
			b := float64(c.B) + db

			out.SetNRGBA(x-bounds.Min.X, y-bounds.Min.Y, color.NRGBA{
				R: clamp8(r*(1-targetBlend) + target[0]*targetBlend),
				G: clamp8(g*(1-targetBlend) + target[1]*targetBlend),
				B: clamp8(b*(1-targetBlend) + target[2]*targetBlend),
				A: c.A,
			})
		}
	}
	return out
}

func blendWithTarget(out *image.NRGBA, target image.Image, alpha float64) {
	bounds := out.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			mosaicColor := out.NRGBAAt(x, y)
			targetColor := nrgbaAt(target, x, y)

			out.SetNRGBA(x, y, color.NRGBA{
				R: blendChannel(mosaicColor.R, targetColor.R, alpha),
				G: blendChannel(mosaicColor.G, targetColor.G, alpha),
				B: blendChannel(mosaicColor.B, targetColor.B, alpha),
				A: blendChannel(mosaicColor.A, targetColor.A, alpha),
			})
		}
	}
}

func blendChannel(base, overlay uint8, alpha float64) uint8 {
	return clamp8(float64(base)*(1-alpha) + float64(overlay)*alpha)
}

func rgbTo8Bit(c RGB) [3]float64 {
	return [3]float64{c[0] / 257, c[1] / 257, c[2] / 257}
}

func clamp01(v float64) float64 {
	if v <= 0 {
		return 0
	}
	if v >= 1 {
		return 1
	}
	return v
}

func clamp8(v float64) uint8 {
	if v <= 0 {
		return 0
	}
	if v >= 255 {
		return 255
	}
	return uint8(v + 0.5)
}
