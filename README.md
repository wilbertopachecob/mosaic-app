# Mosaic Generator

A Go and React photo mosaic app. Upload an image in the browser or run the CLI, and the backend builds a mosaic from a local tile library.

The current generator uses average-color matching, bilinear tile resizing, per-tile color correction, and an EasyMoza-style source-image blend so the result still reads like the original image from a distance.

## Features

- Browser upload flow with original-image preview and generated mosaic download.
- Terminal/CLI generation from the same binary.
- Recursive tile loading from `tiles/` and nested folders such as `tiles/downloaded/`.
- Flat color swatches are ignored when enough real photo tiles are available.
- Configurable tile size, JPEG quality, output path, tile directory, and original-image blend.
- Go tests for the mosaic generator, handlers, and supporting packages.

## Quick Start

### Prerequisites

- Go 1.23+
- Node.js 18+ and npm
- A tile image library in `tiles/`

### Install

```bash
git clone https://github.com/wilbertopachecob/mosaic-app.git
cd mosaic-app

go mod tidy

cd frontend
npm install
npm run build
cd ..
```

Create `.env` if needed:

```bash
SERVER_PORT=8080
MAX_FILE_SIZE=10485760
TILES_DIR=tiles
LOG_LEVEL=info
```

## Tile Library

Put source tile images under `tiles/`. The loader scans recursively, so this works:

```text
tiles/
  downloaded/
    image_001.jpg
    image_002.jpg
  portraits/
    face_001.jpg
```

More tiles produce better mosaics. For service-like output, use hundreds or thousands of images. A tile set that matches the subject matter also helps: portraits work better for portrait mosaics than random landscapes.

Supported tile formats:

- `.jpg`
- `.jpeg`
- `.png`

## Web App

Run the server:

```bash
go run .
```

Open:

```text
http://localhost:8080
```

The web API uses the default generator style, including source-image blending.

## CLI

Build the binary:

```bash
go build -o server .
```

Generate a mosaic:

```bash
./server input.jpg
```

The command prints the generated image path. By default, output is written next to the input as `<input>_mosaic.<ext>`.

Useful flags:

```bash
./server \
  -tiles tiles \
  -tile-size 8 \
  -blend 0.48 \
  -quality 95 \
  -output tmp/output.jpg \
  input.jpg
```

| Flag | Default | Description |
| --- | --- | --- |
| `-tiles` | `TILES_DIR` or `tiles` | Directory containing tile images. Scanned recursively. |
| `-tile-size` | `20` | Tile size in pixels. Smaller values preserve more detail. |
| `-blend` | `0.42` | Original image blend from `0.0` to `1.0`. Use `0` for pure tile output. |
| `-quality` | `92` | JPEG quality from `1` to `100`. |
| `-output` | `<input>_mosaic.<ext>` | Output file path. Supports `.jpg`, `.jpeg`, and `.png`. |

Recommended EasyMoza-like settings:

```bash
./server -tile-size 6 -blend 0.45 -output output.jpg input.jpg
```

## API

### Health Check

```http
GET /api/health
```

### Generate Mosaic

```http
POST /api/file/upload
Content-Type: multipart/form-data
```

Parameters:

- `imgUpload`: image file
- `tileSize`: tile size in pixels
- `blend`: optional original image blend from `0.0` to `1.0`; defaults to `0.42`

Response:

```json
{
  "mosaicImg": "base64_encoded_image",
  "duration": 2.45,
  "format": "jpeg"
}
```

## Project Structure

```text
mosaic-app/
├── config/             # Environment configuration
├── frontend/           # React frontend
├── lib/
│   ├── mosaic/         # Active mosaic generator
│   ├── img/            # Legacy image helpers/tests
│   └── tiles_db/       # Legacy tile DB helpers/tests
├── models/
├── tiles/              # Local tile library, ignored by git
├── handlers.go         # HTTP handlers
├── main.go             # Web server and CLI entry point
└── README.md
```

## Development

Run tests:

```bash
go test ./...
```

Build the CLI/server binary:

```bash
go build -o server .
```

Frontend development:

```bash
cd frontend
npm start
```

Production frontend build:

```bash
cd frontend
npm run build
```

## Notes On Quality

If the generated image looks worse than online mosaic services, check:

- **Tile count**: use hundreds or thousands of tiles.
- **Tile relevance**: portraits for faces, products for product mosaics, etc.
- **Tile size**: smaller tile sizes preserve more detail but take longer.
- **Blend**: increase `-blend` for stronger resemblance to the source image.
- **Output resolution**: higher-resolution input images give the mosaic more room to read.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE).
