// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package listeners

import (
	"crypto/tls"
	"fmt"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/connectionhandlers"
	"net"
)

// Listener is an interface for listening for incoming connections.
type Listener interface {
	Listen(handleConnection func(net.Conn)) error
	Close() error
	NewConnectionHandler() connectionhandlers.ConnectionHandler
	GetAddress() string
	GetType() string
}

// NewListener Choose the listener and connector based on the listener type
func NewListener(listenerType string, listenAddr string, tlsConfig *tls.Config) (Listener, error) {
	var listener Listener
	switch listenerType {
	case "tcp", "http":
		listener = NewTCPListener(listenAddr)

	case "tls", "https":
		listener = NewTLSListener(listenAddr, tlsConfig)

	default:
		return nil, fmt.Errorf("[NewListener]: Unsupported listener type: %s", listenerType)

	}
	return listener, nil
}
