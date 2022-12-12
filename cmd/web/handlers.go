package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (config *configuration) homePageHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		config.infoLog.Println("Handled by main...")
		config.notFoundError(w)
		return
	}

	templateFiles := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	tmplate, err := template.ParseFiles(templateFiles...)

	if err != nil {
		config.serverError(w, err)
		return
	}

	err = tmplate.Execute(w, nil)
	if err != nil {
		config.serverError(w, err)
		return
	}
}

func (config *configuration) showSnip(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		config.notFoundError(w)
		return
	}

	fmt.Fprintf(w, "Show the snippet having id = %d", id)
}

func (config *configuration) createSnip(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")

		config.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new Snippet"))
}
