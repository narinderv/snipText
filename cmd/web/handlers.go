package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/narinderv/snipText/pkg/forms"
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

	config.renderTemplate(w, r, "create.page.tmpl", &templateData{
		Form: *forms.NewForm(nil),
	})
}

func (config *configuration) createSnip(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		config.clientError(w, http.StatusBadRequest)
		return
	}

	form := *forms.NewForm(r.PostForm)

	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.Permittedvalues("expires", "365", "7", "1")

	if !form.IsValid() {
		config.renderTemplate(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})

		return
	}

	id, err := config.snips.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		config.serverError(w, err)
	}

	// Snip creation confirmation message
	config.sessionManager.Put(r, "flash", "Snip saved successfully")

	// If successful, redirect to the page showing this new snip. Using semantic URLs now
	http.Redirect(w, r, fmt.Sprintf("/sniptext/%d", id), http.StatusSeeOther)
}
