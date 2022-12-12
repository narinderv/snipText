package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (config *configuration) serverError(w http.ResponseWriter, err error) {
	stackTrace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	// Log the actual file and line where the error occured instead of this function
	config.errorLog.Output(2, stackTrace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (config *configuration) clientError(w http.ResponseWriter, status int) {

	http.Error(w, http.StatusText(status), status)
}

func (config *configuration) notFoundError(w http.ResponseWriter) {

	config.clientError(w, http.StatusNotFound)
}
