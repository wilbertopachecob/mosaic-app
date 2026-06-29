package main

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultOutputPath(t *testing.T) {
	assert.Equal(t, "vacation_mosaic.jpg", defaultOutputPath("vacation.jpg"))
	assert.Equal(t, filepath.Join("photos", "vacation_mosaic.png"), defaultOutputPath(filepath.Join("photos", "vacation.png")))
	assert.Equal(t, "image_mosaic.jpg", defaultOutputPath("image"))
}

func TestWriteImageJPEG(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.jpg")
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255})

	require.NoError(t, writeImage(path, img, 90))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}

func TestWriteImagePNG(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.png")
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))

	require.NoError(t, writeImage(path, img, 90))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}

func TestWriteImageUnsupportedExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.gif")
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))

	err := writeImage(path, img, 90)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported output extension")
}

func TestRunCLIClampsQuality(t *testing.T) {
	tilesDir := t.TempDir()
	inputDir := t.TempDir()
	inputPath := filepath.Join(inputDir, "input.png")

	f, err := os.Create(inputPath)
	require.NoError(t, err)
	_, err = f.WriteString("not-a-real-image")
	require.NoError(t, err)
	require.NoError(t, f.Close())

	err = runCLI(cliConfig{
		inputPath:  inputPath,
		outputPath: filepath.Join(t.TempDir(), "out.jpg"),
		tilesDir:   tilesDir,
		tileSize:   20,
		quality:    200,
		blend:      0.5,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode input image")
}
