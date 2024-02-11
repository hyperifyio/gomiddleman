// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import (
	"log"
	"net"
)

func handleConnection(clientConn net.Conn, handler ConnectionHandler, connector Connector) {

	defer clientConn.Close()

	if err := handler.Handle(clientConn); err != nil {
		log.Println(err)
		return
	}
	log.Printf("Accepted connection from %s ('%s', %s) to %s", clientConn.RemoteAddr(), handler.GetCommonName(), handler.GetSerialNumber(), connector.GetTarget())

	// Use the connector to connect to the target
	targetConn, err := connector.Connect()
	if err != nil {
		log.Printf("Failed to connect %s ('%s', %s) to target %s: %v", clientConn.RemoteAddr(), handler.GetCommonName(), handler.GetSerialNumber(), connector.GetTarget(), err)
		return
	}
	defer targetConn.Close()

	// Forward traffic between client and target
	ForwardTraffic(clientConn, targetConn)

}
