package main

import (
	"html/template"
	"path"
	"path/filepath"
)

type templateData struct {
	CSRFToken       string
	MaxUploadSize   int64
	Message         string
	ParsingStatus   bool
	UploadMIMETypes []string
}

var pathToTemplates = "./ui/html/"

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	//find all files with .tmpl extension
	pages, err := filepath.Glob(path.Join(pathToTemplates, "/pages/*.tmpl"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		//parse base template
		ts, err := template.New(name).ParseFiles(path.Join(pathToTemplates, "/base.tmpl"))
		if err != nil {
			return nil, err
		}
		//parse partials folder
		ts, err = ts.ParseGlob(path.Join(pathToTemplates, "/partials/*.tmpl"))
		if err != nil {
			return nil, err
		}
		//parse the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
