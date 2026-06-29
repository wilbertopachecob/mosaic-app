// Package mosaic provides a library for generating photo mosaics in Go.
//
// A photo mosaic is an image composed of many smaller tile images, where each
// tile is chosen to match the average color of the corresponding region in
// the target image. This implementation is based on the algorithm described
// in https://medium.com/sausheong/creating-a-photo-mosaic-web-app-in-go-887c8c502f82
// with improvements for accuracy and robustness.
//
// # Usage
//
//	g := mosaic.NewGenerator()
//	if err := g.LoadTiles("path/to/tiles"); err != nil {
//	    log.Fatal(err)
//	}
//	out, err := g.Generate(targetImage, 20) // 20px tile size
//
// # Algorithm
//
// 1. Load tile images from a directory and compute each tile's average RGB.
// 2. For each tile-sized region of the target image, compute the region's average color.
// 3. Find the tile whose average color is closest (Euclidean distance in RGB space).
// 4. Resize the tile to fit the region and draw it onto the output.
//
// The package uses only the Go standard library (image, image/draw, math, os, path/filepath).
package mosaic
