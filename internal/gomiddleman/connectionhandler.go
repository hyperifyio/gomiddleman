// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import (
	"net"
)

// ConnectionHandler defines the interface for handling connections.
type ConnectionHandler interface {
	Handle(conn net.Conn) error
	GetCommonName() string
	GetSerialNumber() string
}
