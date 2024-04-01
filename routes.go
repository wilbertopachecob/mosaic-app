package main

import "net/http"

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	// serve static assests
	mux.Handle("/", http.FileServer(http.Dir("frontend/build")))
	// api routes
	mux.HandleFunc("/api/file/upload", mosaic)
	return mux
}
