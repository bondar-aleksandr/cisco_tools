package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"github.com/bondar-aleksandr/ios-config-parsing/parser"
	log "github.com/sirupsen/logrus"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := &templateData{}
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) configParser(w http.ResponseWriter, r *http.Request) {
	data := &templateData{}
	app.render(w, http.StatusOK, "config-parser.tmpl", data)
}

func (app *application) configUpload(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	file, fileHeader, err := r.FormFile("configFile")
    if err != nil {
        log.Errorf("Error Retrieving the File: %s", err)
		app.clientError(w, http.StatusBadRequest)
        return
    }
    defer file.Close()
	log.Infof("Uploaded File: %+v\n", fileHeader.Filename)
	log.Infof("File Size: %+v\n", fileHeader.Size)
	log.Infof("MIME Header: %+v\n", fileHeader.Header.Values("Content-Type"))

    fileBytes, err := io.ReadAll(file)
    if err != nil {
        app.serverError(w, err)
		return
    }
    
	//processing part
	osFamily := r.FormValue("osFamily")
	outputFormat := r.FormValue("outputFormat")
	log.Infof("Selected the following: OS family - %s, Outputformat - %s", osFamily, outputFormat)

	buf := new(bytes.Buffer)
	buf.Write(fileBytes)

	interface_map := parser.Parsing(buf, osFamily)

	tempFile, err := os.CreateTemp("./temp", fmt.Sprintf("output-*.%s", outputFormat))
	if err != nil {
		log.Error(err)
		app.serverError(w, err)
		return
	}

	if outputFormat == "csv" {
		interface_map.ToCSV(tempFile)
	} else if outputFormat == "json" {
		interface_map.ToJSON(tempFile)
	}
	log.Infof("Saved as %s", tempFile.Name())
	tempFile.Close()

	// clean up
	defer func() {
		err = os.Remove(tempFile.Name())
		if err != nil {
			log.Error(err)
        }
    }()

	//response to client
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(tempFile.Name()))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, tempFile.Name())

	// data := &templateData{}
	// app.render(w, http.StatusAccepted, "config-upload.tmpl", data)
}