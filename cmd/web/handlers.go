package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/narinderv/snipText/pkg/models"
)

func (config *configuration) homePageHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		config.infoLog.Println("Handled by main...")
		config.notFoundError(w)
		return
	}

	snips, err := config.snips.GetLatest()
	if err != nil {
		config.serverError(w, err)
		return
	}

	config.renderTemplate(w, r, "home.page.tmpl", &templateData{AllSnips: snips})
}

func (config *configuration) showSnip(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

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

func (config *configuration) createSnip(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")

		config.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "New Snip"
	content := "This is a new snip from the code"
	expiry := "10"

	id, err := config.snips.Insert(title, content, expiry)
	if err != nil {
		config.serverError(w, err)
	}

	// If successful, redirect to the page showing this new snip
	http.Redirect(w, r, fmt.Sprintf("/sniptext?id=%d", id), http.StatusSeeOther)
}
