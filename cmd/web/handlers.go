package main

import (
	"bytes"
	"fmt"
	"github.com/bondar-aleksandr/ios-config-parsing/parser"
	"github.com/gabriel-vasile/mimetype"
	"github.com/justinas/nosurf"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
)

var tempDir = "./temp"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := &templateData{}
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) configParserHome(w http.ResponseWriter, r *http.Request) {
	data := &templateData{
		MaxUploadSize:   appConfig.Server.MaxUpload,
		CSRFToken:       nosurf.Token(r),
		UploadMIMETypes: appConfig.Server.UploadMIMETypes,
	}
	app.render(w, http.StatusOK, "configParserHome.tmpl", data)
}

// configUpload func parses http request, gets values from hmtl form, and calls "parser.Parsing" func.
// After parsing is done, it puts the result into ./temp/ directory and shows result download page to user.
// Result filename transered to "configDownload" handler as session key
func (app *application) configUpload(w http.ResponseWriter, r *http.Request) {
	//parsing part
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

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		app.serverError(w, err)
		return
	}
	mtype := mimetype.Detect(fileBytes)
	allowed := appConfig.Server.UploadMIMETypes
	if !mimetype.EqualsAny(mtype.String(), allowed...) {
		log.Errorf("Wrong filetype uploaded, got: %s, expect: %s",
			mtype, appConfig.Server.UploadMIMETypes)
		data := &templateData{
			Message: "Only text/plain file types upload is allowed",
		}
		app.render(w, http.StatusUnprocessableEntity, "configParserResult.tmpl", data)
		return
	}
	log.Infof("Detected MIME type: %s", mtype)

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
		app.render(w, http.StatusUnprocessableEntity, "configParserResult.tmpl", data)
		return
	}

	tempFile, err := os.CreateTemp(tempDir, fmt.Sprintf("output-*.%s", outputFormat))
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

	//response to client
	app.sessionManager.Put(r.Context(), "downloadFile", tempFile.Name())
	data := &templateData{
		Message:       "Parsed successfully. Download link below",
		ParsingStatus: true,
	}
	app.render(w, http.StatusOK, "configParserResult.tmpl", data)
}

// configDownload func gets result filename from session, and serve this file to user as attachment.
// Afterwards, result file is deleted from "./temp" directory
func (app *application) configDownload(w http.ResponseWriter, r *http.Request) {

	downloadFile := app.sessionManager.PopString(r.Context(), "downloadFile")
	// for cases when path is accessed directly
	if downloadFile == "" {
		data := &templateData{
			Message: "You need to upload configFile first!",
		}
		app.render(w, http.StatusBadRequest, "configParserResult.tmpl", data)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(downloadFile))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, downloadFile)

	// clean up
	defer func() {
		err := os.Remove(downloadFile)
		if err != nil {
			log.Error(err)
		}
		log.Infof("%s file deleted", downloadFile)
	}()
}
