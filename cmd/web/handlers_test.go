package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// Individual Handler test
func TestPingHandler(t *testing.T) {

	// Create a dummy configuration entry
	config := newTestConfig(t)

	// Recorder to hold the handler response
	rec := httptest.NewRecorder()

	// Dummy HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Call the handler to be tested
	config.pingResponse(rec, req)

	// Get the result of the call
	res := rec.Result()

	// Check if status code is as expected (200)
	if res.StatusCode != http.StatusOK {
		t.Errorf("want %d got %d", http.StatusOK, res.StatusCode)
	}

	// Also check the body and see if it is as expected
	// Defer the close of the response body
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body %s got %s", "OK", string(body))
	}
}

// Test Ping Flow
func TestPingFlow(t *testing.T) {

	// Create a dummy configuration entry
	config := newTestConfig(t)

	// Start a test server
	testServer := newTestServer(t, config.routes())
	defer testServer.Close()

	// Now hit the URL
	statusCode, _, body := testServer.get(t, "/ping")

	// Check if status code is as expected (200)
	if statusCode != http.StatusOK {
		t.Errorf("want %d got %d", http.StatusOK, statusCode)
	}

	// Also check the body and see if it is as expected
	if string(body) != "OK" {
		t.Errorf("want body %s got %s", "OK", string(body))
	}
}

func TestShowSnip(t *testing.T) {

	// Create a dummy configuration entry
	config := newTestConfig(t)

	// Start a test server
	testServer := newTestServer(t, config.routes())
	defer testServer.Close()

	// Set up test cases
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/sniptext/1", http.StatusOK, []byte("Snip For Testing")},
		{"Non-existent ID", "/sniptext/2", http.StatusNotFound, nil},
		{"Negative ID", "/sniptext/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/sniptext/1.23", http.StatusNotFound, nil},
		{"String ID", "/sniptext/foo", http.StatusNotFound, nil},
		{"Empty ID", "/sniptext/", http.StatusNotFound, nil},
		{"Trailing slash", "/sniptext/1/", http.StatusNotFound, nil},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			// Now hit the URL
			statusCode, _, body := testServer.get(t, v.urlPath)

			// Check if status code is as expected
			if statusCode != v.wantCode {
				t.Errorf("want %d got %d", v.wantCode, statusCode)
			}

			// Also check the body and see if it is as expected
			if !bytes.Contains(body, v.wantBody) {
				t.Errorf("want body %s", string(v.wantBody))
			}
		})
	}
}

func TestSignup(t *testing.T) {

	// Create a dummy configuration entry
	config := newTestConfig(t)

	// Start a test server
	testServer := newTestServer(t, config.routes())
	defer testServer.Close()

	// First get the CSRF token for the request
	_, _, body := testServer.get(t, "/user/signup")
	csrfToken := getCSRFToken(t, body)

	// Test cases
	tests := []struct {
		name      string
		userName  string
		email     string
		password  string
		csrfToken string
		wantCode  int
	}{
		{"ValidCase", "User1", "user@example.com", "pa$$worD", csrfToken, http.StatusOK},
		{"Empty name", "", "user@example.com", "validPa$$word", csrfToken, http.StatusOK},
		{"Empty email", "user", "", "validPa$$word", csrfToken, http.StatusOK},
		{"Empty password", "user", "user@example.com", "", csrfToken, http.StatusOK},
		{"Invalid email (incomplete domain)", "user", "user@example.", "validPa$$word", csrfToken, http.StatusOK},
		{"Invalid email (missing @)", "user", "userexample.com", "validPa$$word", csrfToken, http.StatusOK},
		{"Invalid email (missing local part)", "user", "@example.com", "validPa$$", csrfToken, http.StatusOK},
		{"Short password", "user", "user@example.com", "pa$$word", csrfToken, http.StatusOK},
		{"Duplicate email", "user", "narinderv@gmail.com", "validPa$$word", csrfToken, http.StatusOK},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {

			form := url.Values{}
			form.Add("name", v.userName)
			form.Add("email", v.email)
			form.Add("password", v.password)
			form.Add("csrf_token", v.csrfToken)

			// Now hit the URL
			statusCode, _, _ := testServer.postForm(t, "/user/signup", form)

			// Check if status code is as expected
			if statusCode != v.wantCode {
				t.Errorf("want %d got %d", v.wantCode, statusCode)
			}
		})

	}

}
