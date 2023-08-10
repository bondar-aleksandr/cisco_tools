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

func (app *application) configParserHome(w http.ResponseWriter, r *http.Request) {
	data := &templateData{
		MaxUploadSize: appConfig.Server.MaxUpload,
	}
	app.render(w, http.StatusOK, "configParserHome.tmpl", data)
}

func (app *application) configUpload(w http.ResponseWriter, r *http.Request) {
	// limit upload file size
	r.Body = http.MaxBytesReader(w, r.Body, app.config.Server.MaxUpload)
	err := r.ParseMultipartForm(app.config.Server.MaxUpload)
	
	if err != nil {
		log.Errorf("Error parsing multipart form: %s", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("configFile")
    if err != nil {
        log.Errorf("Error retrieving the File: %s", err)
		app.clientError(w, http.StatusBadRequest)
        return
    }
    defer file.Close()
	log.Infof("Uploaded File: %+v", fileHeader.Filename)
	log.Infof("File Size: %+v", fileHeader.Size)
	log.Infof("MIME Header: %+v", fileHeader.Header.Values("Content-Type"))

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

	interface_map, err := parser.Parsing(buf, osFamily)
	if err != nil {
		data := &templateData{
			Message: "It's not config file, or there is no interfaces in it. Parsing failed.",
		}
		app.render(w, http.StatusUnprocessableEntity, "configParserAction.tmpl", data)
		return
	}

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
}