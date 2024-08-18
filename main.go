package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"wilbertopachecob/mosaic/lib/tiles_db"
)

var TILESDB map[string][3]float64

func main() {
	serverPort := getEnvVar("SERVER_PORT")

	mux := routes()

	TILESDB = tiles_db.TilesDB()
	fmt.Println("Mosaic server started at http://localhost:" + serverPort)
	http.ListenAndServe(":"+serverPort, mux)
}

func getEnvVar(envVar string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(envVar)
}
