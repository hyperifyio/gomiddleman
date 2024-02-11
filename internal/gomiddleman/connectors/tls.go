// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package connectors

import (
	"crypto/tls"
	"fmt"
	"net"
)

type TLSConnector struct {
	target string
	config *tls.Config
}

func NewTLSConnector(target string, config *tls.Config) *TLSConnector {
	return &TLSConnector{target, config}
}

func (connector *TLSConnector) Connect() (net.Conn, error) {
	conn, err := tls.Dial("tcp", connector.target, connector.config)
	if err != nil {
		return nil, fmt.Errorf("[TLSConnector.Connect]: failed connect to %s: %v", connector.target, err)
	}
	return conn, nil
}

func (connector *TLSConnector) GetTarget() string {
	return connector.target
}

func (connector *TLSConnector) GetType() string {
	return "tls"
}
