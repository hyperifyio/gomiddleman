// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package tlsutils

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"testing"
)

func TestLoadTLSConfig(t *testing.T) {

	// Setup: Assume you've prepared these files in your test environment
	certFile := "../../../cert.pem"
	keyFile := "../../../key.pem"
	caFile := "../../../ca.pem"

	tlsConfig := LoadTLSConfig(certFile, keyFile, caFile)

	// Verify that tlsConfig is not nil
	if tlsConfig == nil {
		t.Fatal("Expected non-nil tls.Config")
	}

	// Verify that Certificates is correctly loaded
	if len(tlsConfig.Certificates) != 1 {
		t.Errorf("Expected 1 certificate, got %d", len(tlsConfig.Certificates))
	}

	// Verify that RootCAs and ClientCAs are correctly loaded
	if tlsConfig.RootCAs == nil || tlsConfig.ClientCAs == nil {
		t.Fatal("Expected non-nil RootCAs and ClientCAs")
	}

	// Load the CA file to compare
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		t.Fatalf("Failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		t.Fatal("Failed to append CA certificate to pool for comparison")
	}

	// This is a basic check, for a more thorough test you might compare the actual
	// certificates loaded into the pools by parsing them and comparing fields.
	if len(tlsConfig.RootCAs.Subjects()) != len(caCertPool.Subjects()) {
		t.Error("RootCAs loaded does not match the CA certificate provided")
	}

	// Verify ClientAuth is correctly set
	if tlsConfig.ClientAuth != tls.RequireAndVerifyClientCert {
		t.Errorf("Expected ClientAuth to be RequireAndVerifyClientCert, got %v", tlsConfig.ClientAuth)
	}
}
