// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import (
	"fmt"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/connectors"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/listeners"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/tlsutils"
	"github.com/hyperifyio/gomiddleman/internal/testutils/clientutils"
	"github.com/hyperifyio/gomiddleman/internal/testutils/httputils"
	"github.com/hyperifyio/gomiddleman/internal/testutils/serverutils"
	"log"
	"testing"
)

// Test case for HTTPS to middleman to HTTP backend
func TestProxyForwardingHTTPSToHTTP(t *testing.T) {
	
	listenPort := "8080"
	caFile := "../../ca.pem"
	testContent := "Hello from backend"
	proxyClose := setupProxyTest(
		t,
		"tcp",
		"tls",
		listenPort,
		"../../cert.pem",
		"../../key.pem",
		caFile,
		"../../client-cert.pem",
		"../../client-key.pem",
		caFile,
		testContent,
	)
	defer proxyClose()
	client := clientutils.ConfigureClient(
		t,
		true,
		"../../client-cert.pem",
		"../../client-key.pem",
		caFile,
	)
	targetAddr := fmt.Sprintf("https://localhost:%s", listenPort)
	httputils.MakeRequestAndVerifyResponse(t, client, targetAddr, testContent)

}

// Test case for HTTPS to middleman to HTTPS backend
func TestProxyForwardingHTTPSToHTTPS(t *testing.T) {

	listenPort := "8443"
	caFile := "../../ca.pem"
	testContent := "Hello from backend"
	proxyClose := setupProxyTest(
		t,
		"tls",
		"tls",
		listenPort,
		"../../cert.pem",
		"../../key.pem",
		caFile,
		"../../client-cert.pem",
		"../../client-key.pem",
		caFile,
		testContent,
	)
	defer proxyClose()
	client := clientutils.ConfigureClient(
		t,
		true,
		"../../client-cert.pem",
		"../../client-key.pem",
		caFile,
	)
	targetAddr := fmt.Sprintf("https://localhost:%s", listenPort)

	log.Printf("targetAddr = %s", targetAddr)

	httputils.MakeRequestAndVerifyResponse(t, client, targetAddr, testContent)

}

func setupProxyTest(
	t *testing.T,
	connectorType string,
	listenerType string,
	listenPort string,
	certFile string,
	keyFile string,
	caFile string,
	connectorCertFile string,
	connectorKeyFile string,
	connectorCaFile string,
	expected string,
) func() {

	listenAddr := fmt.Sprintf(":%s", listenPort)

	server, serverURL := serverutils.SetupBackendServer(
		connectorType == "tls",
		expected,
		certFile,
		keyFile,
		caFile,
	)

	tlsConfig := tlsutils.LoadTLSConfig(certFile, keyFile, caFile)
	connectorTlsConfig := tlsutils.LoadTLSConfig(connectorCertFile, connectorKeyFile, connectorCaFile)

	listener, err := listeners.NewListener(listenerType, listenAddr, tlsConfig)
	if err != nil {
		t.Fatalf("[setupProxyTest]: Failed to initialize listeners: %v", err)
	}

	connector, err := connectors.NewConnector(serverURL, connectorTlsConfig)
	if err != nil {
		t.Fatalf("[setupProxyTest]: Failed to initialize connector: %v", err)
	}

	if err := StartProxy(listener, connector); err != nil {
		t.Fatalf("[setupProxyTest]: Error when starting the proxyutils: %v", err)
	}

	return func() {
		server.Close()
		if err := listener.Close(); err != nil {
			t.Errorf("[setupProxyTest]: Failed to close listener: %v", err)
		}
	}

}
