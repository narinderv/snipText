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
	// Enable the session Manager for the handlers
	mux.Get("/", config.sessionManager.Enable(config.authenticate(noSurf(http.HandlerFunc(config.homePageHandler)))))
	mux.Get("/sniptext/create", config.sessionManager.Enable(config.authenticate(config.requireAuthenticatedUser(noSurf(http.HandlerFunc(config.createSnipForm))))))
	mux.Post("/sniptext/create", config.sessionManager.Enable(config.authenticate(config.requireAuthenticatedUser(noSurf(http.HandlerFunc(config.createSnip))))))

	// Semantic URLs will be used e.g., http://ip:port/base/query-string
	mux.Get("/sniptext/:id", config.sessionManager.Enable(config.authenticate(noSurf(http.HandlerFunc(config.showSnip)))))

	// SignUp
	mux.Get("/user/signup", config.sessionManager.Enable(config.authenticate(noSurf(http.HandlerFunc(config.signupForm)))))
	mux.Post("/user/signup", config.sessionManager.Enable(config.authenticate(noSurf(http.HandlerFunc(config.signup)))))

	// Login
	mux.Get("/user/login", config.sessionManager.Enable(config.authenticate(noSurf(http.HandlerFunc(config.loginForm)))))
	mux.Post("/user/login", config.sessionManager.Enable(config.authenticate(noSurf(http.HandlerFunc(config.login)))))

	// Logout
	mux.Post("/user/logout", config.sessionManager.Enable(config.authenticate(config.requireAuthenticatedUser(noSurf(http.HandlerFunc(config.logout))))))

	// Ping Request
	mux.Get("/ping", http.HandlerFunc(config.pingResponse))

	// Static File server
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return config.recoverFromPanic(config.logRequests(addSecureHeaders(mux)))
}
