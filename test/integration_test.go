//go:build integration
// +build integration

package challenge_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

const (
	apiBaseURL = "http://localhost:10000/v1"
	maxRetries = 5
	retryDelay = 2 * time.Second
)

func waitForAPI(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
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
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to connect to service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
