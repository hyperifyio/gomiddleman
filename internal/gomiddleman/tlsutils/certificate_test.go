// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package tlsutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"
)

func TestSignCSR(t *testing.T) {
	// Step 1: Generate CA private key and certificate
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate CA private key: %v", err)
	}
	caCertTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "Test CA",
			Organization: []string{"Test Org"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caCertBytes, err := x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		t.Fatalf("Failed to generate CA certificate: %v", err)
	}
	caCert, err := x509.ParseCertificate(caCertBytes)
	if err != nil {
		t.Fatalf("Failed to parse CA certificate: %v", err)
	}

	// Step 2: Generate a private key and CSR
	clientPrivateKey, err := GeneratePrivateKey(2048)
	if err != nil {
		t.Fatalf("Failed to generate private key for CSR: %v", err)
	}
	csr, err := CreateCertificateRequest(clientPrivateKey)
	if err != nil {
		t.Fatalf("Failed to create CSR: %v", err)
	}

	// Step 3: Sign the CSR
	serialNumber := big.NewInt(2) // Example serial number
	signedCertBytes, err := SignCSR(caCert, caPrivateKey, csr, serialNumber)
	if err != nil {
		t.Fatalf("SignCSR failed: %v", err)
	}

	// Step 4: Parse and verify the signed certificate
	signedCert, err := x509.ParseCertificate(signedCertBytes)
	if err != nil {
		t.Fatalf("Failed to parse signed certificate: %v", err)
	}

	// Verify attributes of the signed certificate
	if signedCert.SerialNumber.Cmp(serialNumber) != 0 {
		t.Errorf("Expected serial number %v, got %v", serialNumber, signedCert.SerialNumber)
	}

	if signedCert.Subject.CommonName != "Client CN" {
		t.Errorf("Expected subject CommonName 'Client CN', got '%s'", signedCert.Subject.CommonName)
	}

	if !signedCert.NotBefore.Before(time.Now()) || !signedCert.NotAfter.After(time.Now()) {
		t.Errorf("Signed certificate validity period is incorrect")
	}

	if signedCert.KeyUsage&x509.KeyUsageDigitalSignature == 0 || signedCert.ExtKeyUsage[0] != x509.ExtKeyUsageClientAuth {
		t.Errorf("Signed certificate key usage or extended key usage does not include expected values")
	}
}
