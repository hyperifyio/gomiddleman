// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package tlsutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func TestGenerateClientCertificate(t *testing.T) {

	// Generate a CA private key
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate CA private key: %v", err)
	}

	// Create a CA certificate
	caCertTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "Test CA",
			Organization: []string{"Test Org"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caCertBytes, err := x509.CreateCertificate(rand.Reader, &caCertTemplate, &caCertTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		t.Fatalf("Failed to create CA certificate: %v", err)
	}

	caCert, err := x509.ParseCertificate(caCertBytes)
	if err != nil {
		t.Fatalf("Failed to parse CA certificate: %v", err)
	}

	// Now test GenerateClientCertificate function
	serialNumber := big.NewInt(12345)
	certPEM, keyPEM, err := GenerateClientCertificate(caCert, caPrivateKey, serialNumber)
	if err != nil {
		t.Fatalf("GenerateClientCertificate failed: %v", err)
	}

	// Load the generated certificate
	block, _ := pem.Decode(certPEM)
	if block == nil {
		t.Fatal("Failed to decode PEM block containing the certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse generated certificate: %v", err)
	}

	// Verify the certificate attributes
	if cert.SerialNumber.Cmp(serialNumber) != 0 {
		t.Errorf("Expected serial number %v, got %v", serialNumber, cert.SerialNumber)
	}

	if cert.Subject.CommonName != "Client CN" {
		t.Errorf("Expected CommonName 'Client CN', got '%s'", cert.Subject.CommonName)
	}

	if len(cert.Subject.Organization) == 0 || cert.Subject.Organization[0] != "Client Org" {
		t.Errorf("Expected Organization 'Client Org', got '%v'", cert.Subject.Organization)
	}

	if !cert.NotBefore.Before(time.Now()) || !cert.NotAfter.After(time.Now()) {
		t.Errorf("Certificate validity is outside expected range")
	}

	// Load the generated key
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		t.Fatal("Failed to decode PEM block containing the private key")
	}

	if keyBlock.Type != "RSA PRIVATE KEY" {
		t.Errorf("Expected RSA PRIVATE KEY, got %s", keyBlock.Type)
	}

}
