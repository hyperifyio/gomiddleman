// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package tlsutils

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
)

func LoadTLSConfig(
	certFile,
	keyFile,
	caFile string,
) *tls.Config {
	var certificates []tls.Certificate

	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("[LoadTLSConfig]: Failed to load serverutils certificate and key: %v", err)
		}
		certificates = append(certificates, cert)
	}

	caCertPool := x509.NewCertPool()
	if caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			log.Fatalf("[LoadTLSConfig]: Failed to load CA certificate: %v", err)
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}

	tlsConfig := &tls.Config{
		Certificates: certificates,
		RootCAs:      caCertPool,
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return tlsConfig
}
