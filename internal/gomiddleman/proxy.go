// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import (
	"fmt"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/connectors"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/listeners"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/proxyutils"
	"log"
	"net"
)

// StartProxy starts the proxy server and returns a function to stop it.
func StartProxy(listener listeners.Listener, connector connectors.Connector) error {

	handleConnection := func(clientConn net.Conn) {
		handler := listener.NewConnectionHandler()
		proxyutils.HandleConnection(clientConn, handler, connector)
	}

	// Start listening for connections
	if err := listener.Listen(handleConnection); err != nil {
		return fmt.Errorf("[StartProxy]: Error in listen: %v", err)
	}

	log.Printf(
		"[StartProxy]: Proxy listening on %s (%s) and forwarding to %s (%s)",
		listener.GetAddress(),
		listener.GetType(),
		connector.GetTarget(),
		connector.GetType(),
	)

	return nil
}
