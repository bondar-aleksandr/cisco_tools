package main

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	// "os"
	"time"
	"html/template"
)

type application struct {
	// infoLog *log.Logger
	// errorLog *log.Logger
	// snippets *models.SnippetModel
	// users *models.UserModel
	templateCache map[string]*template.Template
	// formDecoder *form.Decoder
	// sessionManager *scs.SessionManager
}

const PORT = ":4000"

func main() {

	templateCache, err := newTemplateCache()
	if err != nil {
        log.Fatal(err)
    }

	app := &application{
		templateCache: templateCache,
	}
	
	srv := &http.Server{
		Addr: PORT,
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Infof("starting server on port %s...\n", PORT)

	err = srv.ListenAndServe()
	log.Fatal(err)
}