// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package connectionhandlers

import "net"

type TCPConnectionHandler struct{}

func NewTCPConnectionHandler() *TCPConnectionHandler {
	return &TCPConnectionHandler{}
}

// Handle For now, TCP connections do not require additional setup.
//
//	Future checks or setup can be implemented here.
func (handle *TCPConnectionHandler) Handle(_ net.Conn) error {
	return nil
}

// GetCommonName TCP connections might not have a common name.
func (handle *TCPConnectionHandler) GetCommonName() string {
	return "n/a"
}

// GetSerialNumber TCP connections might not have a serial number.
func (handle *TCPConnectionHandler) GetSerialNumber() string {
	return "n/a"
}
