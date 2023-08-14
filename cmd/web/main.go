package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
	"github.com/alexedwards/scs/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		ReadTimeout int64 `yaml:"readTimeout"`
		WriteTimeout int64 `yaml:"writeTimeout"`
		IdleTimeout int64 `yaml:"idleTimeout"`
		MaxUpload int64 `yaml:"maxUpload"`
	}
}

type application struct {
	templateCache map[string]*template.Template
	sessionManager *scs.SessionManager
	config config
}

var appConfig config

func main() {
	readConfig(&appConfig)
	
	templateCache, err := newTemplateCache()
	if err != nil {
        log.Fatal(err)
    }

	sessionManager := scs.New()
	sessionManager.Cookie.Persist = false

	app := &application{
		templateCache: templateCache,
		config: appConfig,
		sessionManager: sessionManager,
	}
	
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port),
		Handler: app.routes(),
		IdleTimeout: time.Duration(app.config.Server.IdleTimeout) * time.Second,
		ReadTimeout: time.Duration(app.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(app.config.Server.WriteTimeout) * time.Second,
	}

	log.Infof(`
	"starting server with following parameters:
	address - %s
	idleTimeout - %d
	readTimeout - %d
	writeTimeout - %d
	maxUpload - %d"`, srv.Addr, srv.IdleTimeout,
	srv.ReadTimeout, srv.WriteTimeout, app.config.Server.MaxUpload)
	
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