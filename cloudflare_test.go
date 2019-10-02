package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
)

func TestIndexHandler(t *testing.T) {
	environment = "dev"
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

func TestJsonify(t *testing.T) {
	var recordType, name, ip string = "A", "sub.example.com", "127.0.0.1"
	var proxied bool = true
	jsonData := jsonify(recordType, name, ip, proxied)
	if len(jsonData) == 0 {
		t.Fail()
	}
}

func TestGetIP(t *testing.T) {
	ip := getIP()
	if len(ip) == 0 {
		t.Fail()
	}
}

func TestGetCredentials(t *testing.T) {
	os.Setenv("CF_EMAIL", "email@email.com")
	os.Setenv("CF_KEY", "myAPIKey")
	os.Setenv("CF_ZONE", "myDNSZone")
	email, gapik, zone := getCredentials()

	if email == "" || gapik == "" || zone == "" {
		t.Fail()
	}
}
