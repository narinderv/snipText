package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// HTTP Handlers
	mux.HandleFunc("/sniptext", showSnip)
	mux.HandleFunc("/sniptext/create", createSnip)
	mux.HandleFunc("/", homePageHandler)

	// Static File server
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting the server on port 8888")
	err := http.ListenAndServe(":8888", mux)
	log.Fatal(err)
}
