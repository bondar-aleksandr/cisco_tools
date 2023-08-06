package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "home")

	data := &templateData{}
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) configParser(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "config-parser page")
	data := &templateData{}
	app.render(w, http.StatusOK, "config-parser.tmpl", data)
}

func (app *application) configUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "config upload")
}