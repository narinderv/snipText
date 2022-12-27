package main

import (
	"html/template"
	"path/filepath"

	"github.com/narinderv/snipText/pkg/models"
)

type templateData struct {
	Snip     *models.SnipText
	AllSnips []*models.SnipText
}

func templateCache(baseDir string) (map[string]*template.Template, error) {
	// Create a new cache to return
	templateCache := map[string]*template.Template{}

	// Get all the page templates
	pages, err := filepath.Glob(filepath.Join(baseDir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Build the templates for all pages
	for _, page := range pages {
		// Page name
		name := filepath.Base(page)

		// Get the page template
		tmpl, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the base layout files to the template
		tmpl, err = tmpl.ParseGlob(filepath.Join(baseDir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the partial template files to the template
		tmpl, err = tmpl.ParseGlob(filepath.Join(baseDir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the final template to the cache
		templateCache[name] = tmpl
	}

	return templateCache, nil
}
