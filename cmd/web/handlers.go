package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		log.Println("Handled by main...")
		http.NotFound(w, r)
		return
	}

	templateFiles := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	tmplate, err := template.ParseFiles(templateFiles...)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmplate.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func showSnip(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Show the snippet having id = %d", id)
}

func createSnip(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")

		/*
			Cleaner way to do the below tasks
			w.WriteHeader(405)
			w.Write([]byte("Method Not Allowed"))
		*/
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new Snippet"))
}
