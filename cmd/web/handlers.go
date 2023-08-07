package main

import (
	// "fmt"
	"io/ioutil"
	"net/http"
	"os"
	// "strings"

	"github.com/bondar-aleksandr/ios-config-parsing/parser"
	log "github.com/sirupsen/logrus"
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
	log.Infof("MIME Header: %+v\n", fileHeader.Header)

	tempFile, err := os.CreateTemp("./temp", "config-*.txt")
	if err != nil {
		log.Error(err)
		app.serverError(w, err)
		return
	}
	
	// clean up
	defer func() {
        err := os.Remove(tempFile.Name())
        if err != nil {
			log.Error(err)
            app.serverError(w, err)
        }
    }()

    // read all of the contents of our uploaded file into a byte array
    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        app.serverError(w, err)
		return
    }
    // write this byte array to our temporary file
    tempFile.Write(fileBytes)

	//TODO: processing part

	osFamily := r.FormValue("osFamily")
	outputFormat := r.FormValue("outputFormat")
	log.Infof("Selected the following: OS family - %s, Outputformat - %s", osFamily, outputFormat)

	interface_map := parser.Parsing(tempFile, osFamily)
	if outputFormat == "csv" {
		parser.ToCSV(interface_map, "./csv-output.csv")
	} else if outputFormat == "json" {
		interface_map.ToJSON("./json-output.json")
	}

	//end of processing
	tempFile.Close()

    // return that we have successfully uploaded our file!
	data := &templateData{}
	app.render(w, http.StatusAccepted, "config-upload.tmpl", data)
}