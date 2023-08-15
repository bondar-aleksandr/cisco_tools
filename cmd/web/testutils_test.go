package main

import (
    "bytes"
    "io"
    // "log"
    "net/http"
    "net/http/httptest"
    "testing"
	"net/http/cookiejar"
	"github.com/alexedwards/scs/v2"
)

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
    pathToTemplates = "./../../ui/html/"
	tc, _ := newTemplateCache()

	sessionManager := scs.New()
	sessionManager.Cookie.Persist = false

	app := application{
		templateCache: tc,
		sessionManager: sessionManager,
	}
	return &app
}

// Define a custom testServer type which embeds a httptest.Server instance.
type testServer struct {
    *httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
    ts := httptest.NewServer(h)

	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
    // Add the cookie jar to the test server client. Any response cookies will
    // now be stored and sent with subsequent requests when using this client.
    ts.Client().Jar = jar

    return &testServer{ts}
}

// Implement a get() method on our custom testServer type. This makes a GET
// request to a given url path using the test server client, and returns the 
// response status code, headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
    rs, err := ts.Client().Get(ts.URL + urlPath)
    if err != nil {
        t.Fatal(err)
    }

    defer rs.Body.Close()
    body, err := io.ReadAll(rs.Body)
    if err != nil {
        t.Fatal(err)
    }
    bytes.TrimSpace(body)

    return rs.StatusCode, rs.Header, string(body)
}