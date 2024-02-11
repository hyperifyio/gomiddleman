// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package connectors

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/connectionhandlers"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/tlsutils"
	"math/big"
	"net"
)

type TLSConnector struct {
	target       string
	config       *tls.Config
	caCert       *x509.Certificate
	caPrivateKey *rsa.PrivateKey
}

func NewTLSConnector(
	target string,
	config *tls.Config,
	caCert *x509.Certificate,
	caPrivateKey *rsa.PrivateKey,
) *TLSConnector {
	return &TLSConnector{target, config, caCert, caPrivateKey}
}

func (connector *TLSConnector) Connect(
	handler connectionhandlers.ConnectionHandler,
) (net.Conn, error) {

	serialNumberStr := handler.GetSerialNumber()
	serialNumber := new(big.Int)
	serialNumber.SetString(serialNumberStr, 16)

	// Generate a private key for the dynamic client certificate
	privateKey, err := tlsutils.GeneratePrivateKey(2048)
	if err != nil {
		return nil, fmt.Errorf("[TLSConnector.Connect]: failed to generate dynamic private client key: %v", err)
	}

	// Generate a CSR for the dynamic client certificate
	csr, err := tlsutils.CreateCertificateRequest(privateKey)
	if err != nil {
		return nil, fmt.Errorf("[TLSConnector.Connect]: failed to create CSR for dynamic mtls client key: %v", err)
	}

	// Sign the CSR with the CA to get the client certificate
	certBytes, err := tlsutils.SignCSR(connector.caCert, connector.caPrivateKey, csr, serialNumber)
	if err != nil {
		return nil, fmt.Errorf("[TLSConnector.Connect]: failed to sign dynamic mtls client certificate: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	// Load the client certificate
	cert, err := tls.X509KeyPair(certPEM, privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("[TLSConnector.Connect]: failed to load client certificate: %v", err)
	}

	// Update the TLS config with the dynamic client certificate
	tlsConfig := connector.config.Clone()
	tlsConfig.Certificates = append(tlsConfig.Certificates, cert)

	conn, err := tls.Dial("tcp", connector.target, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("[TLSConnector.Connect]: failed to connect to %s: %v", connector.target, err)
	}
	return conn, nil
}

func (connector *TLSConnector) GetTarget() string {
	return connector.target
}

func (connector *TLSConnector) GetType() string {
	return "tls"
}
