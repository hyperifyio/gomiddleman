// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package connectors

import (
	"fmt"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/connectionhandlers"
	"net"
)

type TCPConnector struct {
	target string
}

func NewTCPConnector(target string) *TCPConnector {
	return &TCPConnector{target}
}

// Connect Logic to connect to a TCP target
func (connector *TCPConnector) Connect(_ connectionhandlers.ConnectionHandler) (net.Conn, error) {
	conn, err := net.Dial("tcp", connector.target)
	if err != nil {
		return nil, fmt.Errorf("[TCPConnector.Connect]: failed to connect to TCP target %s: %v", connector.target, err)
	}
	return conn, nil
}

func (connector *TCPConnector) GetTarget() string {
	return connector.target
}

func (connector *TCPConnector) GetType() string {
	return "tcp"
}
