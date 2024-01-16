package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/h3th-IV/H/internal/assert"
)

func TestPing(t *testing.T) {
	//init new httptest response recorder
	responseRecorder := httptest.NewRecorder()

	//nothing but a test request
	testRequest, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	//ping test instead of responsewriter use response Recorder and r(test request) as request
	ping(responseRecorder, testRequest)

	responseResult := responseRecorder.Result()

	//test if response code is "200 OK"
	assert.Equal(t, responseResult.StatusCode, http.StatusOK)

	defer responseResult.Body.Close()

	body, err := io.ReadAll(responseResult.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")

}
