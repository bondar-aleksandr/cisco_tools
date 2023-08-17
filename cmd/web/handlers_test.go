package main

import (
	"bytes"
	"github.com/bondar-aleksandr/ios-config-parsing/internal/assert"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	// "net/url"
)

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

	_, _, body := ts.get(t, "/config-parser")
	validCSRFToken := extractCSRFToken(t, body)
	t.Logf("CSRF token is: %q", validCSRFToken)

	mbody := new(bytes.Buffer)

	var mw = *multipart.NewWriter(mbody)
	_ = mw.WriteField("osFamily", "ios")
	_ = mw.WriteField("outputFormat", "csv")
	_ = mw.WriteField("csrf_token", validCSRFToken)

	file, err := os.Open("./../../parser/test_data/ASR-P.txt")
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

	code, _, _ := ts.post(t, "/config-parser/upload", mw.FormDataContentType(), mbody)
	if code != http.StatusOK {
		t.Errorf("wrong status code")
	}
}
