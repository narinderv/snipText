package main

import (
	"fmt"
	"net/http"
)

func addSecureHeaders(nxtHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode:block")
		w.Header().Set("X-Framw-Options", "deny")

		nxtHandler.ServeHTTP(w, r)
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
