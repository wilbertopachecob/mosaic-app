package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	imgpkg "wilbertopachecob/mosaic/lib/img"

	"wilbertopachecob/mosaic/lib/tiles_db"
)

var TILESDB map[string][3]float64

func main() {
	serverPort := getEnvVar("SERVER_PORT")

	mux := routes()

	TILESDB = tiles_db.TilesDB()
	fmt.Println("Mosaic server started")
	http.ListenAndServe(":"+serverPort, mux)
}

func getEnvVar(envVar string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(envVar)
}

func mosaic(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	r.ParseMultipartForm(10485760)

	file, _, err := r.FormFile("imgUpload")
	if err != nil {
		fmt.Println("error with uploaded file: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		data := struct {
			Error string `json:"error"`
		}{Error: err.Error()}

		json.NewEncoder(w).Encode(data)
	}

	defer file.Close()
	tileSize, _ := strconv.Atoi(r.FormValue("tileSize"))
	original, _, _ := image.Decode(file)
	bounds := original.Bounds()

	newImage := image.NewNRGBA(image.Rect(bounds.Min.X, bounds.Min.X, bounds.Max.X, bounds.Max.Y))
	db := tiles_db.CloneTilesDB(TILESDB)

	sourcePoint := image.Point{0, 0}

	for y := bounds.Min.Y; y < bounds.Max.Y; y = y + tileSize {
		for x := bounds.Min.X; x < bounds.Max.X; x = x + tileSize {
			r, g, b, _ := original.At(x, y).RGBA()
			color := [3]float64{float64(r), float64(g), float64(b)}
			nearestFileByColor := imgpkg.Nearest(color, &db)
			file, err := os.Open(nearestFileByColor)
			if err == nil {
				img, _, err := image.Decode(file)
				if err == nil {
					t := imgpkg.Resize(img, tileSize)
					tile := t.SubImage(t.Bounds())
					tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
					draw.Draw(newImage, tileBounds, tile, sourcePoint, draw.Src)
				} else {
					fmt.Println("error: ", err, nearestFileByColor)
				}
			} else {
				fmt.Println("error: ", nearestFileByColor)
			}
			file.Close()
		}
	}

	buf2 := new(bytes.Buffer)
	jpeg.Encode(buf2, newImage, nil)
	mosaicImg := base64.StdEncoding.EncodeToString(buf2.Bytes())
	t1 := time.Now()
	secondsDiff := t1.Sub(t0).Seconds()
	duration := math.Round(secondsDiff*100) / 100

	data := struct {
		MosaicImg string  `json:"mosaicImg"`
		Duration  float64 `json:"duration"`
	}{MosaicImg: mosaicImg, Duration: duration}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(data)
}
