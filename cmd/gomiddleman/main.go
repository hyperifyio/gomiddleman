// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	listenerType = flag.String("type", getEnvOrDefault("GOMIDDLEMAN_TYPE", "tls"), "type of proxy (tcp or tls)")
	listenPort   = flag.String("port", getEnvOrDefault("GOMIDDLEMAN_PORT", "8080"), "port on which the proxy listens")
	target       = flag.String("target", getEnvOrDefault("GOMIDDLEMAN_TARGET", "http://localhost:3000"), "target where to proxy connections")
	certFile     = flag.String("cert", getEnvOrDefault("GOMIDDLEMAN_CERT_FILE", "cert.pem"), "proxy certificate as PEM file")
	keyFile      = flag.String("key", getEnvOrDefault("GOMIDDLEMAN_KEY_FILE", "key.pem"), "proxy key as PEM file")
	caFile       = flag.String("ca", getEnvOrDefault("GOMIDDLEMAN_CA_FILE", "ca.pem"), "proxy ca as PEM file")
)

func main() {

	flag.Parse()

	var wg sync.WaitGroup

	targetURL, err := url.Parse(*target)
	if err != nil {
		log.Fatalf("Invalid target URL: %v", err)
	}

	listenAddr := fmt.Sprintf(":%s", *listenPort)

	var listener gomiddleman.Listener

	// Choose the listener and connector based on the listener type
	if *listenerType == "tls" {
		tlsConfig := gomiddleman.LoadTLSConfig(*certFile, *keyFile, *caFile)
		listener = gomiddleman.NewTLSListener(listenAddr, tlsConfig)
	} else if *listenerType == "tcp" {
		listener = gomiddleman.NewTCPListener(listenAddr)
	} else {
		log.Fatalf("Unsupported listener type: %s", *listenerType)
	}
	defer listener.Close()

	var connector gomiddleman.Connector
	switch targetURL.Scheme {

	case "tcp", "http":
		connector = gomiddleman.NewTCPConnector(targetURL.Host)

	case "tls", "https":
		// TLSConnector needs a tls.Config to establish TLS connections
		if *listenerType != "tls" {
			log.Fatalf("TLS connector requires tls proxy type")
		}
		tlsConfig := gomiddleman.LoadTLSConfig(*certFile, *keyFile, *caFile)
		connector = gomiddleman.NewTLSConnector(targetURL.Host, tlsConfig)

	default:
		log.Fatalf("Unsupported target scheme: %s", targetURL.Scheme)
	}

	if err := gomiddleman.StartProxy(listener, connector); err != nil {
		log.Fatalf("Error when starting the proxy: %v", err)
	}

	// Setup signal handling for graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	log.Println("Shutting down proxy...")
	if err := listener.Close(); err != nil {
		log.Printf("Failed to close listener: %v", err)
	}
	wg.Wait()

}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
