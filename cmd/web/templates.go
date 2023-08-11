package main

import (
	"html/template"
	"path/filepath"
)

type templateData struct {
	// CurrentYear int
	// Snippet *models.Snippet
	// Snippets []*models.Snippet
	// Form any
	// Flash string
	// IsAuthenticated bool
	CSRFToken string
	MaxUploadSize int64
	Message string
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	//find all files with .tmpl extension
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
    if err != nil {
        return nil, err
    }
	for _, page := range pages {
		name := filepath.Base(page)
		//parse base template
		ts, err := template.New(name).ParseFiles("./ui/html/base.tmpl")
        if err != nil {
            return nil, err
        }
		//parse partials folder
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
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