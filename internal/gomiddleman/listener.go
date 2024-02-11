// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import "net"

// Listener is an interface for listening for incoming connections.
type Listener interface {
	Listen(handleConnection func(net.Conn)) error
	Close() error
	NewConnectionHandler() ConnectionHandler
	GetAddress() string
}
