// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package clientutils

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"testing"
	"time"
)

func ConfigureClient(
	t *testing.T,
	useClientCert bool,
	clientCertFile string,
	clientKeyFile string,
	caFile string,
) *http.Client {

	certPool := x509.NewCertPool()

	// Use CA to trust the server cert
	if useClientCert {
		certData, err := os.ReadFile(caFile)
		if err != nil {
			t.Fatalf("[ConfigureClient]: failed to read CA certificate file: %v", err)
		}
		if !certPool.AppendCertsFromPEM(certData) {
			t.Fatalf("[ConfigureClient]: failed to append CA certificate to pool")
		}
	}

	tlsConfig := &tls.Config{
		RootCAs:   certPool,
		ClientCAs: certPool,
	}

	if useClientCert {
		clientCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			t.Fatalf("[ConfigureClient]: failed to load certificate: %v", err)
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 5 * time.Second,
	}

	return client
}
