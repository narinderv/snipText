package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (config *configuration) routes() http.Handler {
	// HTTP Request Handlers
	// We are now using a different HTTP router (pat) than the default ServMux.
	// This provides additional functionalities of Method based routing (GET, POST, etc.)
	// and also semantic URL processing where query string parameters can be made part of the URL itself.
	// This helps us in implementing RESTful routes
	mux := pat.New()

	// HTTP Handlers
	mux.Get("/", http.HandlerFunc(config.homePageHandler))
	mux.Get("/sniptext/create", http.HandlerFunc(config.createSnipForm))
	mux.Post("/sniptext/create", http.HandlerFunc(config.createSnip))

	// Semantic URLs will be used e.g., http://ip:port/base/query-string
	mux.Get("/sniptext/:id", http.HandlerFunc(config.showSnip))

	// Static File server
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return config.recoverFromPanic(config.logRequests(addSecureHeaders(mux)))
}
