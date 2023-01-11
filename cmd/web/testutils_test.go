package main

import (
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/narinderv/snipText/pkg/models/stubs"
)

var csrfTokenRX = regexp.MustCompile("<input type='hidden' name='csrf_token' value='(.+)'>")

func getCSRFToken(t *testing.T, body []byte) string {

	match := csrfTokenRX.FindSubmatch(body)
	if len(match) < 2 {
		t.Fatal("No csrf token found in the body")
	}

	return html.UnescapeString(string(match[1]))
}

// Return a test configuration structure
func newTestConfig(t *testing.T) *configuration {

	tempCache, err := templateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	session := sessions.New([]byte("s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge"))
	session.Lifetime = 12 * time.Hour

	// Create a dummy configuration entry
	config := configuration{
		errorLog:       log.New(ioutil.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		sessionManager: session,
		snips:          &stubs.SnipModel{},
		users:          &stubs.UserModel{},
		templateCache:  tempCache,
	}

	return &config
}

// Return a test server
type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, handler http.Handler) *testServer {

	server := httptest.NewTLSServer(handler)

	// Enable storing of cookies, if any returned by the tested handler
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	server.Client().Jar = jar

	// Disable any redirects, if returned on calling the handler and use the last returned result
	server.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{server}
}

// Get method for the test Server
func (server *testServer) get(t *testing.T, url string) (int, http.Header, []byte) {

	// Now hit the URL
	res, err := server.Client().Get(server.URL + url)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return res.StatusCode, res.Header, body
}

// Post method for the test Server
func (server *testServer) postForm(t *testing.T, url string, formData url.Values) (int, http.Header, []byte) {

	// Now hit the URL
	res, err := server.Client().PostForm(server.URL+url, formData)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return res.StatusCode, res.Header, body
}
