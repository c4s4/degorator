package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingRoute(t *testing.T) {
	router := http.HandlerFunc(handler)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/hello?name=World", nil)
	router.ServeHTTP(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("%d != 200\n", recorder.Code)
	}
	if recorder.Body.String() != "Hello World!" {
		t.Errorf("%s != 'Hello World!'", recorder.Body.String())
	}
}
