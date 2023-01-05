package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (config *configuration) authenticatedUser(r *http.Request) int {
	return config.sessionManager.GetInt(r, "userID")
}

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

func (config *configuration) addCommonData(data *templateData, r *http.Request) *templateData {

	if data == nil {
		data = &templateData{}
	}

	data.CurrentYear = time.Now().Year()

	data.Flash = config.sessionManager.PopString(r, "flash")

	data.AuthenticatedUser = config.authenticatedUser(r)
	data.CSRFToken = nosurf.Token(r)

	return data
}

func (config *configuration) renderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data *templateData) {

	tmplate, ok := config.templateCache[templateName]
	if !ok {
		config.serverError(w, fmt.Errorf("the template %s, does not exist", templateName))
		return
	}

	buffer := new(bytes.Buffer)

	err := tmplate.Execute(buffer, config.addCommonData(data, r))
	if err != nil {
		config.serverError(w, err)
		return
	}

	buffer.WriteTo(w)
}
