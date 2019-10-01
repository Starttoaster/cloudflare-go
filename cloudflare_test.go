package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	respRecord := httptest.NewRecorder()
	handler := setupHandlers()
	handler.ServeHTTP(respRecord, req)

	resp := respRecord.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
	}
}
