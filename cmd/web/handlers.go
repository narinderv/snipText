package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/narinderv/snipText/pkg/models"
)

func (config *configuration) homePageHandler(w http.ResponseWriter, r *http.Request) {

	snips, err := config.snips.GetLatest()
	if err != nil {
		config.serverError(w, err)
		return
	}

	config.renderTemplate(w, r, "home.page.tmpl", &templateData{AllSnips: snips})
}

func (config *configuration) showSnip(w http.ResponseWriter, r *http.Request) {

	// pat Module parses the query string along with the ':' value itself
	idStr := r.URL.Query().Get(":id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		config.notFoundError(w)
		return
	}

	snipInfo, err := config.snips.Get(id)
	if err == models.ErrNoRecord {
		config.notFoundError(w)
		return
	} else if err != nil {
		config.serverError(w, err)
		return
	}

	config.renderTemplate(w, r, "snip.page.tmpl", &templateData{Snip: snipInfo})
}

func (config *configuration) createSnipForm(w http.ResponseWriter, r *http.Request) {

	config.renderTemplate(w, r, "create.page.tmpl", nil)
}

func (config *configuration) createSnip(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		config.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expiry := r.PostForm.Get("expires")

	// Validation: Create a map of field to error and validate all the fields
	errors := make(map[string]string)

	// Title
	if strings.TrimSpace(title) == "" {
		errors["title"] = "Title cannot be empty"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "Title length cannot be more than 100 characters"
	}

	// Content
	if strings.TrimSpace(content) == "" {
		errors["content"] = "Content cannot be empty"
	}

	// Expiry
	if strings.TrimSpace(expiry) == "" {
		errors["expiry"] = "Please select an expiry duration"
	} else if expiry != "1" && expiry != "7" && expiry != "365" {
		errors["expiry"] = "Invalid expiry value"
	}

	id, err := config.snips.Insert(title, content, expiry)
	if err != nil {
		config.serverError(w, err)
	}

	// Check if any validations have failed
	if len(errors) > 0 {
		// Errors have occured. Re-display the form with the error messages and prefilled values
		config.renderTemplate(w, r, "create.page.tmpl", &templateData{
			FormData: r.PostForm,
			Errors:   errors,
		})

		return
	}

	// If successful, redirect to the page showing this new snip. Using semantic URLs now
	http.Redirect(w, r, fmt.Sprintf("/sniptext/%d", id), http.StatusSeeOther)
}
