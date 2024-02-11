// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package connectors

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/connectionhandlers"
	"net"
	"net/url"
)

// Connector is an interface for connecting to a target server.
type Connector interface {
	Connect(handler connectionhandlers.ConnectionHandler) (net.Conn, error)
	GetTarget() string
	GetType() string
}

func NewConnector(
	target string,
	tlsConfig *tls.Config,
	caCert *x509.Certificate,
	caPrivateKey *rsa.PrivateKey,
) (Connector, error) {
	var connector Connector
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("[NewConnector]: Invalid target URL: %v", err)
	}

	switch targetURL.Scheme {

	case "tcp", "http":
		connector = NewTCPConnector(
			targetURL.Host,
		)

	case "tls", "https":
		connector = NewTLSConnector(
			targetURL.Host,
			tlsConfig,
			caCert,
			caPrivateKey,
		)

	default:
		return nil, fmt.Errorf("[NewConnector]: Unsupported target scheme: %s", targetURL.Scheme)
	}
	return connector, nil
}
