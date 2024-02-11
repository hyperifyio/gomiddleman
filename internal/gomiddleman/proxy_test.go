// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// TestProxyForwarding verifies the proxy forwards HTTP requests and responses correctly.
func TestProxyForwarding(t *testing.T) {

	// Load the certificate
	certData, err := ioutil.ReadFile("../../cert.pem")
	if err != nil {
		t.Fatalf("Failed to read certificate file: %v", err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(certData) {
		t.Fatalf("Failed to append certificate to pool")
	}

	// Start a mock HTTP backend server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a simple message
		w.Write([]byte("Hello from backend"))
	}))
	defer mockServer.Close()

	tlsConfig := LoadTLSConfig("../../cert.pem", "../../key.pem", "../../ca.pem")

	// Load client certificate and key
	clientCert, err := tls.LoadX509KeyPair("../../client-cert.pem", "../../client-key.pem")
	if err != nil {
		t.Fatalf("Failed to load client certificate: %v", err)
	}

	// Create an HTTP client with a timeout
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      certPool,
				Certificates: []tls.Certificate{clientCert},
			},
		},
		Timeout: 5 * time.Second, // Set an appropriate timeout duration
	}

	mockServerURL, err := url.Parse(mockServer.URL)
	if err != nil {
		t.Fatalf("Failed to parse mock server URL: %v", err)
	}

	stopProxy := StartProxy("8080", mockServerURL.Host, tlsConfig)
	defer stopProxy()

	// Make an HTTP request through the proxy to the mock server
	resp, err := client.Get("https://localhost:8080")
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
