// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import "net"

// Connector is an interface for connecting to a target server.
type Connector interface {
	Connect() (net.Conn, error)
	GetTarget() string
}
