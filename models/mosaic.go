package models

// MosaicRequest represents the request structure for mosaic generation
type MosaicRequest struct {
	TileSize int `json:"tileSize"`
}

// MosaicResponse represents the response structure for mosaic generation
type MosaicResponse struct {
	MosaicImg string  `json:"mosaicImg"`
	Duration  float64 `json:"duration"`
	Error     string  `json:"error,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Color represents RGB color values
type Color [3]float64

// Tile represents a tile in the mosaic database
type Tile struct {
	Filename string
	Color    Color
}
