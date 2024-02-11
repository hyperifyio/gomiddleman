// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package proxyutils

import (
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/connectionhandlers"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/connectors"
	"log"
	"net"
)

func HandleConnection(
	clientConn net.Conn,
	handler connectionhandlers.ConnectionHandler,
	connector connectors.Connector,
) {

	defer func() {
		if err := clientConn.Close(); err != nil {
			log.Printf("[handleConnection]: Failed to close client connection %s ('%s', %s) to target %s: %v", clientConn.RemoteAddr(), handler.GetCommonName(), handler.GetSerialNumber(), connector.GetTarget(), err)
		}
	}()

	if err := handler.Handle(clientConn); err != nil {
		log.Println(err)
		return
	}
	log.Printf("[handleConnection]: Accepted connection from %s ('%s', %s) to %s", clientConn.RemoteAddr(), handler.GetCommonName(), handler.GetSerialNumber(), connector.GetTarget())

	// Use the connector to connect to the target
	targetConn, err := connector.Connect(handler)
	if err != nil {
		log.Printf("[handleConnection]: Failed to connect %s ('%s', %s) to target %s: %v", clientConn.RemoteAddr(), handler.GetCommonName(), handler.GetSerialNumber(), connector.GetTarget(), err)
		return
	}
	defer func() {
		if err := targetConn.Close(); err != nil {
			log.Printf("[handleConnection]: Failed to close target connection %s ('%s', %s) to target %s: %v", clientConn.RemoteAddr(), handler.GetCommonName(), handler.GetSerialNumber(), connector.GetTarget(), err)
		}
	}()

	// Forward traffic between proxyutils connection and target backend
	ForwardTraffic(clientConn, targetConn)

}
