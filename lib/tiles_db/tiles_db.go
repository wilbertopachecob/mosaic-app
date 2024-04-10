package tiles_db

import (
	"fmt"
	"image"
	"os"

	imgpkg "wilbertopachecob/mosaic/lib/img"
)

func TilesDB() map[string][3]float64 {
	fmt.Println("start populating tiles db")
	db := make(map[string][3]float64)
	files, _ := os.ReadDir("tiles")
	for _, f := range files {
		name := "tiles/" + f.Name()
		file, err := os.Open(name)
		if err == nil {
			img, _, err := image.Decode(file)
			if err == nil {
				db[name] = imgpkg.AverageColor(img)
			} else {
				fmt.Println("error populating TILEDB:", err, name)
			}

		} else {
			fmt.Println("can not open file", name, err)
		}
		file.Close()
	}
	fmt.Println("finished populating tiles db")
	return db
}

func CloneTilesDB(tilesDB map[string][3]float64) map[string][3]float64 {
	db := make(map[string][3]float64)
	for k, v := range tilesDB {
		db[k] = v
	}
	return db
}
