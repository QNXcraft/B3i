package device

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListApps(t *testing.T) {
	// Mock BB10 server
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/apps" {
			apps := []App{
				{ID: "com.example.app", Name: "Example App", Version: "1.0.0", Status: "running"},
			}
			json.NewEncoder(w).Encode(apps)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Insecure skip verify for test server
	client := &Client{
		BaseURL: server.URL,
		HTTPClient: server.Client(),
	}

	apps, err := client.ListApps()
	if err != nil {
		t.Fatalf("Failed to list apps: %v", err)
	}

	if len(apps) != 1 {
		t.Errorf("Expected 1 app, got %d", len(apps))
	}

	if apps[0].ID != "com.example.app" {
		t.Errorf("Expected app ID com.example.app, got %s", apps[0].ID)
	}
}

func TestLogin(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth" {
			w.Header().Set("X-Auth-Token", "test-token")
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Password:   "password",
		HTTPClient: server.Client(),
	}

	err := client.Login()
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if client.Token != "test-token" {
		t.Errorf("Expected token test-token, got %s", client.Token)
	}
}
