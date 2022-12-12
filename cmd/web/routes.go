package main

import (
	"net/http"
)

func (config *configuration) routes() *http.ServeMux {
	// HTTP Request Handlers
	mux := http.NewServeMux()

	// HTTP Handlers
	mux.HandleFunc("/sniptext", config.showSnip)
	mux.HandleFunc("/sniptext/create", config.createSnip)
	mux.HandleFunc("/", config.homePageHandler)

	// Static File server
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
