// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import (
	"fmt"
	"net"
)

type TCPConnector struct {
	target string
}

func NewTCPConnector(target string) *TCPConnector {
	return &TCPConnector{target}
}

// Logic to connect to a TCP target
func (c *TCPConnector) Connect() (net.Conn, error) {
	conn, err := net.Dial("tcp", c.target)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to TCP target %s: %v", c.target, err)
	}
	return conn, nil
}

func (c *TCPConnector) GetTarget() string {
	return c.target
}
