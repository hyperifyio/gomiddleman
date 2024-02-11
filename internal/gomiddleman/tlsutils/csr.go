// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package tlsutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
)

func CreateCertificateRequest(privateKey *rsa.PrivateKey) (*x509.CertificateRequest, error) {

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   "Client CN",
			Organization: []string{"Client Org"},
		},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return nil, err
	}

	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		return nil, err
	}

	return csr, nil
}
