package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"wilbertopachecob/mosaic/config"
	"wilbertopachecob/mosaic/lib/mosaic"
	"wilbertopachecob/mosaic/middleware"
)

func main() {
	cfg := config.Load()
	cliCfg := parseCLIFlags(cfg.TilesDir)
	if cliCfg.inputPath != "" {
		if err := runCLI(cliCfg); err != nil {
			log.Fatal(err)
		}
		return
	}

	fmt.Println("Initializing mosaic generator...")
	generator := mosaic.NewGenerator()

	if err := generator.LoadTiles(cfg.TilesDir); err != nil {
		log.Printf("Warning: Failed to load tiles from %q: %v", cfg.TilesDir, err)
	} else {
		fmt.Printf("Loaded %d tiles successfully\n", generator.TileCount())
	}

	configureServer(cfg, generator)

	rateLimiter := middleware.NewIPRateLimiter(cfg.RateLimitRequests, cfg.RateLimitWindow)
	generationSem := middleware.NewSemaphore(cfg.MaxConcurrentGenerations)

	var upload http.Handler = http.HandlerFunc(uploadHandler)
	upload = rateLimiter.Middleware(upload)
	upload = generationSem.Middleware(upload)
	upload = middleware.Timeout(cfg.RequestTimeout, upload)

	http.Handle("/api/file/upload", upload)
	http.HandleFunc("/api/health", healthHandler)

	fs := http.FileServer(http.Dir("dist/build"))
	http.Handle("/", fs)

	addr := ":" + cfg.ServerPort
	fmt.Printf("Mosaic server starting on http://localhost%s\n", addr)
	if cfg.Production {
		fmt.Println("Production mode enabled (generic server errors)")
	}
	fmt.Printf("Security: max upload %d bytes, image cap %dx%d, tile size %d-%d, rate limit %d/%s, concurrency %d, timeout %s\n",
		cfg.MaxFileSize,
		cfg.MaxImageWidth, cfg.MaxImageHeight,
		cfg.MinTileSize, cfg.MaxTileSize,
		cfg.RateLimitRequests, cfg.RateLimitWindow,
		cfg.MaxConcurrentGenerations,
		cfg.RequestTimeout,
	)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

type cliConfig struct {
	inputPath  string
	outputPath string
	tilesDir   string
	tileSize   int
	quality    int
	blend      float64
}

func parseCLIFlags(defaultTilesDir string) cliConfig {
	var cfg cliConfig
	flag.StringVar(&cfg.tilesDir, "tiles", defaultTilesDir, "directory containing tile images")
	flag.IntVar(&cfg.tileSize, "tile-size", 20, "mosaic tile size in pixels")
	flag.StringVar(&cfg.outputPath, "output", "", "output image path; defaults to <input>_mosaic.jpg")
	flag.IntVar(&cfg.quality, "quality", 92, "JPEG quality from 1 to 100")
	flag.Float64Var(&cfg.blend, "blend", mosaic.DefaultOptions().SourceBlend, "original image blend from 0.0 to 1.0")
	flag.Parse()

	if flag.NArg() > 0 {
		cfg.inputPath = flag.Arg(0)
	}
	return cfg
}

func runCLI(cfg cliConfig) error {
	if cfg.tileSize < 1 {
		return fmt.Errorf("tile-size must be greater than 0")
	}
	if cfg.quality < 1 {
		cfg.quality = 1
	}
	if cfg.quality > 100 {
		cfg.quality = 100
	}
	if cfg.outputPath == "" {
		cfg.outputPath = defaultOutputPath(cfg.inputPath)
	}

	in, err := os.Open(cfg.inputPath)
	if err != nil {
		return fmt.Errorf("open input image: %w", err)
	}
	defer func() { _ = in.Close() }()

	target, _, err := image.Decode(in)
	if err != nil {
		return fmt.Errorf("decode input image: %w", err)
	}

	g := mosaic.NewGenerator()
	if err := g.LoadTiles(cfg.tilesDir); err != nil {
		return fmt.Errorf("load tiles: %w", err)
	}
	if g.TileCount() == 0 {
		return mosaic.ErrNoTiles
	}

	out, err := g.GenerateWithOptions(target, cfg.tileSize, mosaic.Options{SourceBlend: cfg.blend})
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cfg.outputPath), 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := writeImage(cfg.outputPath, out, cfg.quality); err != nil {
		return err
	}

	abs, err := filepath.Abs(cfg.outputPath)
	if err != nil {
		abs = cfg.outputPath
	}
	fmt.Println(abs)
	return nil
}

func defaultOutputPath(inputPath string) string {
	ext := filepath.Ext(inputPath)
	if ext == "" {
		ext = ".jpg"
	}
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	return filepath.Join(filepath.Dir(inputPath), base+"_mosaic"+ext)
}

func writeImage(path string, img image.Image, jpegQuality int) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create output image: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close output image: %w", cerr)
		}
	}()

	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
		if err := png.Encode(f, img); err != nil {
			return fmt.Errorf("encode png: %w", err)
		}
	case ".jpg", ".jpeg", "":
		if err := jpeg.Encode(f, img, &jpeg.Options{Quality: jpegQuality}); err != nil {
			return fmt.Errorf("encode jpeg: %w", err)
		}
	default:
		return fmt.Errorf("unsupported output extension %q; use .jpg, .jpeg, or .png", filepath.Ext(path))
	}

	return nil
}
