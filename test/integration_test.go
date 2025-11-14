//go:build integration
// +build integration

package challenge_test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const (
	apiBaseURL = "https://localhost:10000/v1"
	maxRetries = 5
	retryDelay = 2 * time.Second
)

var httpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func waitForAPI(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := httpClient.Get(url)
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("API not available after %v", timeout)
}

func TestServiceHealth(t *testing.T) {
	// Wait for the service to be available
	url := apiBaseURL + "/users"
	if err := waitForAPI(url, time.Duration(maxRetries)*retryDelay); err != nil {
		t.Fatalf("Service did not start: %v", err)
	}

	// Test the service endpoint
	resp, err := httpClient.Get(url)
	if err != nil {
		t.Fatalf("Failed to connect to service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestListUsers(t *testing.T) {
	resp, err := httpClient.Get(apiBaseURL + "/users")
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestAddUser(t *testing.T) {
	user := map[string]string{
		"firstname": "Test",
		"lastname":  "User",
		"login":     "testuser",
		"password":  "testpass",
	}
	body, _ := json.Marshal(user)

	resp, err := httpClient.Post(apiBaseURL+"/users", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to add user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
}

func TestGetUser(t *testing.T) {
	resp, err := httpClient.Get(apiBaseURL + "/users/1")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestDelUser(t *testing.T) {
	req, _ := http.NewRequest("DELETE", apiBaseURL+"/users/2", nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 204 or 404, got %d", resp.StatusCode)
	}
}
