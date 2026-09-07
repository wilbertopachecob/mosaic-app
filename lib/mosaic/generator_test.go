package mosaic

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// createTestImage creates a simple RGBA image of the given size.
func createTestImage(width, height int, c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func TestGenerator_TileCount_Empty(t *testing.T) {
	g := NewGenerator()
	if g.TileCount() != 0 {
		t.Errorf("TileCount() = %d, want 0", g.TileCount())
	}
}

func TestGenerator_Generate_NoTiles(t *testing.T) {
	g := NewGenerator()
	target := createTestImage(100, 100, color.RGBA{128, 128, 128, 255})

	_, err := g.Generate(target, 20)
	if err != ErrNoTiles {
		t.Errorf("Generate with no tiles: got err %v, want ErrNoTiles", err)
	}
}

func TestGenerator_LoadTiles_InvalidDir(t *testing.T) {
	g := NewGenerator()
	err := g.LoadTiles("/nonexistent/path/12345")
	if err == nil {
		t.Error("LoadTiles with invalid dir: expected error")
	}
	if g.TileCount() != 0 {
		t.Errorf("TileCount after failed load = %d, want 0", g.TileCount())
	}
}

func TestGenerator_LoadTiles_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	g := NewGenerator()
	err := g.LoadTiles(dir)
	if err != nil {
		t.Fatalf("LoadTiles(empty dir): %v", err)
	}
	if g.TileCount() != 0 {
		t.Errorf("TileCount(empty dir) = %d, want 0", g.TileCount())
	}
}

func writeJPEGFile(t *testing.T, path string, img image.Image) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create %s: %v", path, err)
	}
	if err := jpeg.Encode(f, img, nil); err != nil {
		_ = f.Close()
		t.Fatalf("encode %s: %v", path, err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close %s: %v", path, err)
	}
}

func writePNGFile(t *testing.T, path string, img image.Image) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create %s: %v", path, err)
	}
	if err := png.Encode(f, img); err != nil {
		_ = f.Close()
		t.Fatalf("encode %s: %v", path, err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close %s: %v", path, err)
	}
}

func TestGenerator_LoadTiles_WithImages(t *testing.T) {
	dir := t.TempDir()

	writeJPEGFile(t, filepath.Join(dir, "red.jpg"), createTestImage(10, 10, color.RGBA{255, 0, 0, 255}))
	writePNGFile(t, filepath.Join(dir, "blue.png"), createTestImage(10, 10, color.RGBA{0, 0, 255, 255}))

	g := NewGenerator()
	err := g.LoadTiles(dir)
	if err != nil {
		t.Fatalf("LoadTiles: %v", err)
	}
	if g.TileCount() != 2 {
		t.Errorf("TileCount = %d, want 2", g.TileCount())
	}
}

func TestGenerator_LoadTiles_PrefersPhotoTilesOverColorSwatches(t *testing.T) {
	dir := t.TempDir()

	writeJPEGFile(t, filepath.Join(dir, "color_red.jpg"), createTestImage(10, 10, color.RGBA{255, 0, 0, 255}))

	for i := 0; i < 8; i++ {
		photoTile := image.NewRGBA(image.Rect(0, 0, 10, 10))
		for y := 0; y < 10; y++ {
			for x := 0; x < 10; x++ {
				if (x+y+i)%2 == 0 {
					photoTile.Set(x, y, color.RGBA{255, uint8(i * 20), 0, 255})
				} else {
					photoTile.Set(x, y, color.RGBA{0, uint8(255 - i*20), 255, 255})
				}
			}
		}
		writeJPEGFile(t, filepath.Join(dir, fmt.Sprintf("photo_%d.jpg", i)), photoTile)
	}

	g := NewGenerator()
	if err := g.LoadTiles(dir); err != nil {
		t.Fatalf("LoadTiles: %v", err)
	}
	if g.TileCount() != 8 {
		t.Errorf("TileCount = %d, want only 8 photo tiles", g.TileCount())
	}
}

func TestGenerator_Generate_WithTiles(t *testing.T) {
	dir := t.TempDir()

	for i, c := range []color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
	} {
		writeJPEGFile(t, filepath.Join(dir, fmt.Sprintf("tile%d.jpg", i)), createTestImage(20, 20, c))
	}

	g := NewGenerator()
	if err := g.LoadTiles(dir); err != nil {
		t.Fatalf("LoadTiles: %v", err)
	}

	target := createTestImage(40, 40, color.RGBA{200, 10, 10, 255}) // Reddish
	mosaic, err := g.Generate(target, 20)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	bounds := mosaic.Bounds()
	if bounds.Dx() != 40 || bounds.Dy() != 40 {
		t.Errorf("mosaic size = %dx%d, want 40x40", bounds.Dx(), bounds.Dy())
	}
}

func TestGenerator_GenerateWithOptions_SourceBlend(t *testing.T) {
	dir := t.TempDir()
	writeJPEGFile(t, filepath.Join(dir, "blue.jpg"), createTestImage(10, 10, color.RGBA{0, 0, 255, 255}))

	g := NewGenerator()
	if err := g.LoadTiles(dir); err != nil {
		t.Fatalf("LoadTiles: %v", err)
	}

	target := createTestImage(10, 10, color.RGBA{200, 10, 10, 255})
	pureTile, err := g.GenerateWithOptions(target, 10, Options{SourceBlend: 0})
	if err != nil {
		t.Fatalf("GenerateWithOptions pure tile: %v", err)
	}
	fullBlend, err := g.GenerateWithOptions(target, 10, Options{SourceBlend: 1})
	if err != nil {
		t.Fatalf("GenerateWithOptions full blend: %v", err)
	}

	pureColor := color.NRGBAModel.Convert(pureTile.At(0, 0)).(color.NRGBA)
	fullColor := color.NRGBAModel.Convert(fullBlend.At(0, 0)).(color.NRGBA)
	if fullColor.R != 200 || fullColor.G != 10 || fullColor.B != 10 {
		t.Errorf("full blend pixel = %v, want target color", fullColor)
	}
	if pureColor == fullColor {
		t.Errorf("pure tile output should differ from full source blend")
	}
}

func TestGenerator_Generate_TileSizeClamped(t *testing.T) {
	dir := t.TempDir()
	writeJPEGFile(t, filepath.Join(dir, "t.jpg"), createTestImage(5, 5, color.RGBA{128, 128, 128, 255}))

	g := NewGenerator()
	if err := g.LoadTiles(dir); err != nil {
		t.Fatalf("LoadTiles: %v", err)
	}
	target := createTestImage(10, 10, color.RGBA{128, 128, 128, 255})

	// tileSize 0 should be clamped to 1
	_, err := g.Generate(target, 0)
	if err != nil {
		t.Errorf("Generate(tileSize=0): %v", err)
	}
}
