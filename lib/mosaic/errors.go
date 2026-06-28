package mosaic

import "errors"

// ErrNoTiles is returned by Generate when no tiles have been loaded.
var ErrNoTiles = errors.New("mosaic: no tiles loaded; call LoadTiles before Generate")
