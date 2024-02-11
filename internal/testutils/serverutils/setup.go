// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package serverutils

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

func SetupBackendServer(
	useHTTPS bool,
	expected string,
	certFile string,
	keyFile string,
	caFile string,
) (*httptest.Server, string) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(expected))
		if err != nil {
			log.Fatalf("[SetupBackendServer]: Failed to write response: %v", err)
		}
	})

	var server *httptest.Server
	if useHTTPS {

		caCert, err := os.ReadFile(caFile)
		if err != nil {
			log.Fatalf("[SetupBackendServer]: Failed to read CA certificate: %v", err)
		}

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			log.Fatal("[SetupBackendServer]: Failed to append CA certificate")
		}

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("[SetupBackendServer]: Failed to load key pair: %v", err)
		}
		server = httptest.NewUnstartedServer(handler)
		server.TLS = &tls.Config{
			RootCAs:      caCertPool,
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{cert},
		}
		server.StartTLS()
	} else {
		server = httptest.NewServer(handler)
	}

	return server, server.URL
}
