package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	router.Handler(http.MethodGet, "/", standard.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/config-parser", standard.ThenFunc(app.configParser))
    router.Handler(http.MethodPost, "/config-parser/upload", standard.ThenFunc(app.configUpload))

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	
	return standard.Then(router)
}