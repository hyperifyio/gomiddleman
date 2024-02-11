// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

// StartProxy starts the proxy server and returns a function to stop it.
func StartProxy(listenPort string, target string, targetConfig *tls.Config) func() {

	var wg sync.WaitGroup

	listenAddr := fmt.Sprintf(":%s", listenPort)

	listener, err := tls.Listen("tcp", listenAddr, targetConfig)
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

func LoadTLSConfig(certFile, keyFile string, caFile string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load server certificate and key: %v", err)
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
}

func handleConnection(clientConn net.Conn, target string) {

	defer clientConn.Close()

	tlsConn, ok := clientConn.(*tls.Conn)
	if !ok {
		log.Println("Failed to cast to TLS connection")
		return
	}

	// It's a good practice to perform a handshake explicitly to check for any TLS errors early
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("TLS handshake failed: %v", err)
		return
	}

	// After performing the handshake:
	var commonName string
	var serialNumber string
	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) > 0 {
		clientCert := state.PeerCertificates[0]
		commonName = clientCert.Subject.CommonName
		serialNumber = clientCert.SerialNumber.Text(16)
	} else {
		commonName = "n/a"
		serialNumber = "n/a"
	}

	// Connect to the backend server.
	targetConn, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("Failed to connect %s ('%s', %s) to backend server at %s: %v", tlsConn.RemoteAddr(), commonName, serialNumber, target, err)
		return
	}
	defer targetConn.Close()
	log.Printf("Accepted connection from %s ('%s', %s) to %s", tlsConn.RemoteAddr(), commonName, serialNumber, target)

	// Forward traffic between client and backend server concurrently.
	errChan := make(chan error, 2)

	go func() {
		_, err := io.Copy(targetConn, tlsConn)
		errChan <- err
	}()

	go func() {
		_, err := io.Copy(tlsConn, targetConn)
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
