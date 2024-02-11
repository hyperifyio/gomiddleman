// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

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

func (c *TLSConnector) Connect() (net.Conn, error) {
	conn, err := tls.Dial("tcp", c.target, c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to TLS target %s: %v", c.target, err)
	}
	return conn, nil
}

func (c *TLSConnector) GetTarget() string {
	return c.target
}
