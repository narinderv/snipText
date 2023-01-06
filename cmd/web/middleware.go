package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/narinderv/snipText/pkg/models"
)

func noSurf(nxtHandler http.Handler) http.Handler {
	csrfHandler := nosurf.New(nxtHandler)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func addSecureHeaders(nxtHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode:block")
		w.Header().Set("X-Framw-Options", "deny")

		nxtHandler.ServeHTTP(w, r)
	})
}

func (config *configuration) requireAuthenticatedUser(nxtHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if config.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}

		nxtHandler.ServeHTTP(w, r)
	})
}

func (config *configuration) authenticate(nxtHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user ID already exists in the session or not
		exists := config.sessionManager.Exists(r, "userID")
		if !exists {
			nxtHandler.ServeHTTP(w, r)
			return
		}

		// If userID exists, check if this is a valid user
		user, err := config.users.Get(config.sessionManager.GetInt(r, "userID"))
		// User is not valid. Remove userID from session
		if err == models.ErrNoRecord {
			config.sessionManager.Remove(r, "userID")
			nxtHandler.ServeHTTP(w, r)
			return
		} else if err != nil {
			config.serverError(w, err)
			return
		}

		// User is valid. Add the user details to the context
		ctxt := context.WithValue(r.Context(), contextKeyUser, user)
		nxtHandler.ServeHTTP(w, r.WithContext(ctxt))
	})
}

func (config *configuration) logRequests(nxtHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.infoLog.Printf("Req: %s - %s|%s|%s", r.RemoteAddr, r.Proto, r.Method, r.URL)

		nxtHandler.ServeHTTP(w, r)
	})
}

func (config *configuration) recoverFromPanic(nxtHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				config.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		nxtHandler.ServeHTTP(w, r)
	})
}
