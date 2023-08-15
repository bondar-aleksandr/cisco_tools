package main

import (
	"net/http"
	"testing"
	"github.com/bondar-aleksandr/ios-config-parsing/internal/assert"
)

func Test_Home_Pages(t *testing.T) {
	
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	data := []struct{
		name string
		pageURL string
		statusCode int
		body string
	}{
		{"Homepage", "/", http.StatusOK, "<h2>List of apps available</h2>"},
		{"ConfigParserHome", "/config-parser", http.StatusOK, "This app parses text configuration IOS/NXOS file"},
	}

	for _, val := range data {
		t.Run(val.name, func(t *testing.T){
			code, _, body := ts.get(t, val.pageURL)
			assert.Equal(t, code, val.statusCode)
			assert.StringContains(t, body, val.body)
		})
	}
}
