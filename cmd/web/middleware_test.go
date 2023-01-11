package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {

	// Recorder to hold the handler response
	rec := httptest.NewRecorder()

	// Dummy HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Dummy Handler to pass to the secureHeader handler
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Call the handler to be tested
	addSecureHeaders(handle).ServeHTTP(rec, req)

	// Get the result of the call
	res := rec.Result()

	// Check if the required headers are set
	frameOpt := res.Header.Get("X-Frame-Options")
	if frameOpt != "deny" {
		t.Errorf("want %s got %s", "deny", frameOpt)
	}

	xssProtect := res.Header.Get("X-XSS-Protection")
	if xssProtect != "1; mode:block" {
		t.Errorf("want %s got %s", "1; mode:block", xssProtect)
	}

	// Check if the Handler has been called successfully
	if res.StatusCode != http.StatusOK {
		t.Errorf("want %d got %d", http.StatusOK, res.StatusCode)
	}

	// Finally, also check the body and see if it is as expected
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
