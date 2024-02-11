// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package tlsutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

// GenerateClientCertificate generates a client certificate signed by the provided CA.
// The generated certificate includes the specified serial number.
func GenerateClientCertificate(
	caCert *x509.Certificate,
	caPrivateKey *rsa.PrivateKey,
	serialNumber *big.Int,
) (certPEM []byte, keyPEM []byte, err error) {

	// Step 1: Generate a new private key for the client
	clientPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Step 2: Create a certificate request (CSR) for the client
	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   "Client CN",
			Organization: []string{"Client Org"},
		},
	}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, clientPrivateKey)
	if err != nil {
		return nil, nil, err
	}
	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		return nil, nil, err
	}

	// Step 3: Sign the CSR with the CA's private key to generate the client certificate
	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      csr.Subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour), // 1 year validity
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, caCert, csr.PublicKey, caPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	// Step 4: Encode the client private key and certificate to PEM format
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientPrivateKey)})
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})

	return certPEM, keyPEM, nil
}
