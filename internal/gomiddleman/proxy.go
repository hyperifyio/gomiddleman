// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import (
	"fmt"
	"log"
	"net"
)

// StartProxy starts the proxy server and returns a function to stop it.
func StartProxy(listener Listener, connector Connector) error {

	handleConnection := func(clientConn net.Conn) {
		handler := listener.NewConnectionHandler()
		handleConnection(clientConn, handler, connector)
	}

	// Start listening for connections
	if err := listener.Listen(handleConnection); err != nil {
		return fmt.Errorf("Error in listen: %v", err)
	}

	log.Printf("Proxy listening on %s and forwarding to %s", listener.GetAddress(), connector.GetTarget())

	return nil
}
