package main

import (
	"bytes"
	"fmt"
	"github.com/bondar-aleksandr/ios-config-parsing/internal/assert"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

var testDataDir = "./../../test_data"

func Test_Home_Pages(t *testing.T) {

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	data := []struct {
		name       string
		pageURL    string
		statusCode int
		body       string
	}{
		{"Homepage", "/", http.StatusOK, "<h2>List of apps available</h2>"},
		{"ConfigParserHome", "/config-parser", http.StatusOK, "This app parses text configuration IOS/NXOS file"},
	}

	for _, val := range data {
		t.Run(val.name, func(t *testing.T) {
			code, _, body := ts.get(t, val.pageURL)
			assert.Equal(t, code, val.statusCode)
			assert.StringContains(t, body, val.body)
		})
	}
}

func Test_config_Upload(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	data := []struct {
		name       string
		filename   string
		osFamily   string
		statusCode int
		body       string
	}{
		{
			"upload correct MIME type with correct data",
			filepath.Join(testDataDir, "INET-R01.txt"),
			"ios",
			http.StatusOK,
			"Parsed successfully. Download link below",
		},
		{
			"upload wrong MIME type",
			filepath.Join(testDataDir, "INET-R01.json"),
			"ios",
			http.StatusUnprocessableEntity,
			fmt.Sprintf("Only %s file types upload is allowed", app.config.Server.UploadMIMETypes),
		},
		{
			"upload correct MIME type with wrong data",
			filepath.Join(testDataDir, "not_config_textfile.txt"),
			"ios",
			http.StatusUnprocessableEntity,
			"It's not config file",
		},
		{
			"upload bigger than allowed file",
			filepath.Join(testDataDir, "big_pdf.pdf"),
			"ios",
			http.StatusBadRequest,
			"Bad Request",
		},
	}

	for _, val := range data {
		t.Run(val.name, func(t *testing.T) {
			// get CSRF token
			_, _, body := ts.get(t, "/config-parser")
			validCSRFToken := extractCSRFToken(t, body)

			// construct multipart form
			mbody := new(bytes.Buffer)
			mw := *multipart.NewWriter(mbody)
			_ = mw.WriteField("osFamily", "ios")
			_ = mw.WriteField("outputFormat", "csv")
			_ = mw.WriteField("csrf_token", validCSRFToken)

			// attach file to form
			file, err := os.Open(val.filename)
			if err != nil {
				t.Fatal(err)
			}
			w, err := mw.CreateFormFile("configFile", file.Name())
			if err != nil {
				t.Fatal(err)
			}
			if _, err := io.Copy(w, file); err != nil {
				t.Fatal(err)
			}
			mw.Close()

			// make request
			code, _, body := ts.post(t, "/config-parser/upload", mw.FormDataContentType(), mbody)

			// compare request results with expected values
			assert.Equal(t, code, val.statusCode)
			assert.StringContains(t, body, val.body)
		})
	}

}
