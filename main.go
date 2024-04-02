package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var TILESDB map[string][3]float64

func main() {
	serverPort := getEnvVar("SERVER_PORT")

	mux := routes()

	TILESDB = tilesDB()
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
	db := cloneTilesDB()

	sourcePoint := image.Point{0, 0}

	for y := bounds.Min.Y; y < bounds.Max.Y; y = y + tileSize {
		for x := bounds.Min.X; x < bounds.Max.X; x = x + tileSize {
			r, g, b, _ := original.At(x, y).RGBA()
			color := [3]float64{float64(r), float64(g), float64(b)}
			nearestFileByColor := nearest(color, &db)
			file, err := os.Open(nearestFileByColor)
			if err == nil {
				img, _, err := image.Decode(file)
				if err == nil {
					t := resize(img, tileSize)
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

func averageColor(img image.Image) [3]float64 {
	bounds := img.Bounds()
	r, g, b := 0.0, 0.0, 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for i := bounds.Min.X; i < bounds.Max.X; i++ {
			r1, g1, b1, _ := img.At(i, y).RGBA()
			r, g, b = r+float64(r1), g+float64(g1), b+float64(b1)
		}
	}
	totalPixels := float64(bounds.Max.X * bounds.Max.Y)
	return [3]float64{r / totalPixels, g / totalPixels, b / totalPixels}
}

func resize(in image.Image, newWidth int) image.NRGBA {
	bounds := in.Bounds()
	ratio := bounds.Dx() / newWidth
	out := image.NewNRGBA(image.Rect(bounds.Min.X/ratio, bounds.Min.X/ratio, bounds.Max.X/ratio, bounds.Max.Y/ratio))

	for y, j := bounds.Min.Y, bounds.Min.Y; y < bounds.Max.Y; y, j = y+ratio, j+1 {
		for x, i := bounds.Min.X, bounds.Min.X; i < bounds.Max.X; x, i = x+ratio, i+1 {
			r, g, b, a := in.At(x, y).RGBA()
			out.SetNRGBA(i, j, color.NRGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}
	return *out
}

func tilesDB() map[string][3]float64 {
	fmt.Println("start populating tiles db")
	db := make(map[string][3]float64)
	files, _ := os.ReadDir("tiles")
	for _, f := range files {
		name := "tiles/" + f.Name()
		file, err := os.Open(name)
		if err == nil {
			img, _, err := image.Decode(file)
			if err == nil {
				db[name] = averageColor(img)
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

func nearest(target [3]float64, db *map[string][3]float64) string {
	var filename string
	smallest := 1000000.0
	for k, v := range *db {
		dist := distance(target, v)
		if dist < smallest {
			filename, smallest = k, dist
		}
	}
	delete(*db, filename)
	return filename
}

// finds the Euclidian distance between two points
func distance(p1 [3]float64, p2 [3]float64) float64 {
	return math.Sqrt(sq(p2[0]-p1[0]) + sq(p2[1]-p1[1]) + sq(p2[2]-p1[2]))
}

func sq(n float64) float64 {
	return n * n
}

func cloneTilesDB() map[string][3]float64 {
	db := make(map[string][3]float64)
	for k, v := range TILESDB {
		db[k] = v
	}
	return db
}
