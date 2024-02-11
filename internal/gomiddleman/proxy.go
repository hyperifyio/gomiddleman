// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

// StartProxy starts the proxy server and returns a function to stop it.
func StartProxy(listenPort, target string) func() {

	var wg sync.WaitGroup

	listenAddr := fmt.Sprintf(":%s", listenPort)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", listenPort, err)
	}
	log.Printf("Proxy listening on port %s and forwarding to %s", listenPort, target)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {

			// Accept new connections.
			clientConn, err := listener.Accept()
			if err != nil {
				log.Println("Proxy listener closed, stopping accept loop")
				return
			}

			// Handle the connection in a new goroutine.
			go handleConnection(clientConn, target)
		}
	}()

	// Return a function that stops the proxy, cleans up resources, etc.
	return func() {
		log.Println("Shutting down proxy...")
		err := listener.Close()
		if err != nil {
			log.Printf("Failed to close listener on port %s: %v", listenPort, err)
		}
		wg.Wait()
	}
}

func handleConnection(clientConn net.Conn, target string) {

	defer clientConn.Close()

	// Connect to the backend server.
	targetConn, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("Failed to connect %s to backend server at %s: %v", clientConn.RemoteAddr(), target, err)
		return
	}
	defer targetConn.Close()
	log.Printf("Accepted connection from %s to %s", clientConn.RemoteAddr(), target)

	// Forward traffic between client and backend server concurrently.
	errChan := make(chan error, 2)

	go func() {
		_, err := io.Copy(targetConn, clientConn)
		errChan <- err
	}()

	go func() {
		_, err := io.Copy(clientConn, targetConn)
		errChan <- err
	}()

	for i := 0; i < 2; i++ {
		err := <-errChan
		if err != nil {
			log.Printf("Error forwarding traffic: %v", err)
			// Handle the error, e.g., by closing both connections.
			break
		}
	}

}
