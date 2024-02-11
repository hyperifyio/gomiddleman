// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// TestProxyForwarding verifies the proxy forwards HTTP requests and responses correctly.
func TestProxyForwarding(t *testing.T) {

	// Start a mock HTTP backend server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a simple message
		w.Write([]byte("Hello from backend"))
	}))
	defer mockServer.Close()

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 5 * time.Second, // Set an appropriate timeout duration
	}

	mockServerURL, err := url.Parse(mockServer.URL)
	if err != nil {
		t.Fatalf("Failed to parse mock server URL: %v", err)
	}

	stopProxy := StartProxy("8080", mockServerURL.Host)
	defer stopProxy()

	// Make an HTTP request through the proxy to the mock server
	resp, err := client.Get("http://localhost:8080") // Assuming the proxy forwards to mockServer
	if err != nil {
		t.Fatalf("Failed to make request through proxy: %v", err)
	}
	defer resp.Body.Close()

	// Read and verify the response from the backend server via the proxy
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expected := "Hello from backend"
	if string(body) != expected {
		t.Errorf("Expected response body %q, got %q", expected, body)
	}

}
