package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	uiFileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", uiFileServer))

	dataFileServer := http.FileServer(http.Dir("./downloads/"))
	router.Handler(http.MethodGet, "/downloads/*filepath", http.StripPrefix("/downloads", dataFileServer))
	
	dynamic := alice.New(noSurf, app.sessionManager.LoadAndSave)

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/ssh-client", app.sshClient)
	router.HandlerFunc(http.MethodGet, "/config-parser-cli", app.configParserCli)
	router.Handler(http.MethodGet, "/config-parser", dynamic.ThenFunc(app.configParserHome))
    router.Handler(http.MethodPost, "/config-parser/upload", dynamic.ThenFunc(app.configUpload))
	router.Handler(http.MethodGet, "/config-parser/download", dynamic.ThenFunc(app.configDownload))

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	standard := alice.New(app.maxRequestSize, app.recoverPanic, app.logRequest, secureHeaders)
	
	return standard.Then(router)
}