// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package proxyutils

import (
	"io"
	"log"
	"net"
)

// ForwardTraffic concurrently forwards traffic between two net.Conn connections.
func ForwardTraffic(conn1, conn2 net.Conn) {
	errChan := make(chan error, 2)

	// Start forwarding from conn1 to conn2
	go func() {
		_, err := io.Copy(conn2, conn1)
		errChan <- err
	}()

	// Start forwarding from conn2 to conn1
	go func() {
		_, err := io.Copy(conn1, conn2)
		errChan <- err
	}()

	// Wait for forwarding to complete or an error to occur
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			log.Printf("[ForwardTraffic]: Error forwarding traffic: %v", err)
			break
		}
	}
}
