// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package tlsutils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

// ReadCACertificateFile reads and decodes CA certificate
func ReadCACertificateFile(caFile string) *x509.Certificate {

	caCertPEM, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("[ReadCACertificateFile]: Failed to read CA certificate: %v", err)
	}

	// Decode the CA certificate
	caCertBlock, _ := pem.Decode(caCertPEM)
	if caCertBlock == nil {
		log.Fatalf("[ReadCACertificateFile]: Failed to decode CA certificate PEM")
	}
	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		log.Fatalf("[ReadCACertificateFile]: Failed to parse CA certificate: %v", err)
	}

	return caCert

}

// ReadCAKeyFile reads and decodes PEM CA key
func ReadCAKeyFile(caKeyFile string) *rsa.PrivateKey {

	caKeyPEM, err := os.ReadFile(caKeyFile)
	if err != nil {
		log.Fatalf("[ReadCAKeyFile]: Failed to read CA key: %v", err)
	}

	// Decode the CA private key
	caKeyBlock, _ := pem.Decode(caKeyPEM)
	if caKeyBlock == nil {
		log.Fatalf("[ReadCAKeyFile]: Failed to decode CA key PEM")
	}
	keyInterface, err := x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes) // Adjust parsing function based on your key format
	if err != nil {
		log.Fatalf("[ReadCAKeyFile]: Failed to parse CA private key: %v", err)
	}
	caPrivateKey, ok := keyInterface.(*rsa.PrivateKey) // Adjust type assertion based on your key type
	if !ok {
		log.Fatalf("[ReadCAKeyFile]: Decoded key is not RSA private key")
	}

	return caPrivateKey

}
