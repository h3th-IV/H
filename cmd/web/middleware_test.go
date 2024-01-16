package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/h3th-IV/H/internal/assert"
)

func TestMW(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	testRequest, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHaedersMW(next).ServeHTTP(responseRecorder, testRequest)

	reponseResult := responseRecorder.Result()

	// Check if MW has correctly set the Content-Security-Policy
	// header on the response.
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, reponseResult.Header.Get("Content-Security-Policy"), expectedValue)

	//check if MW has set Referer policy
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, reponseResult.Header.Get("Referrer-Policy"), expectedValue)

	//check if mW has set X-Content-Type-Optionsof the header
	expectedValue = "nosniff"
	assert.Equal(t, reponseResult.Header.Get("X-Content-Type-Options"), expectedValue)

	//check if MW set header X-Frame-Options
	expectedValue = "deny"
	assert.Equal(t, reponseResult.Header.Get("X-Frame-Options"), expectedValue)

	//check if mW correctly set X-site scripting protection in header
	expectedValue = "0"
	assert.Equal(t, reponseResult.Header.Get("X-XSS-Protection"), expectedValue)
}
