package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type config struct {
	Server struct {
		Host string
		Port string
		ReadTimeout int
		WriteTimeout int
		IdleTimeout int
		maxUpload int
	}
}

type application struct {
	// infoLog *log.Logger
	// errorLog *log.Logger
	// snippets *models.SnippetModel
	// users *models.UserModel
	templateCache map[string]*template.Template
	// formDecoder *form.Decoder
	// sessionManager *scs.SessionManager
	config config
}

var appConfig config

func main() {
	readConfig(&appConfig)
	
	templateCache, err := newTemplateCache()
	if err != nil {
        log.Fatal(err)
    }

	app := &application{
		templateCache: templateCache,
		config: appConfig,
	}
	
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port),
		Handler: app.routes(),
		IdleTimeout: time.Duration(app.config.Server.IdleTimeout) * time.Second,
		ReadTimeout: time.Duration(app.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(app.config.Server.WriteTimeout) * time.Second,
	}

	log.Infof("starting server on port %s...\n", app.config.Server.Port)

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func readConfig(cfg *config) {
	f, err := os.Open("./config/config.yml")
	if err != nil {
		log.Fatalf("Cannot read app config file because of: %s", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatalf("Cannot parse app config file because of: %s", err)
	}
}