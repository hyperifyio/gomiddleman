// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/connectors"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/listeners"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/tlsutils"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	listenerType = flag.String("type", getEnvOrDefault("GOMIDDLEMAN_TYPE", "tls"), "type of proxyutils (tcp or tls)")
	listenPort   = flag.String("port", getEnvOrDefault("GOMIDDLEMAN_PORT", "8080"), "port on which the proxyutils listens")
	target       = flag.String("target", getEnvOrDefault("GOMIDDLEMAN_TARGET", "http://localhost:3000"), "target where to proxyutils connections")
	certFile     = flag.String("cert", getEnvOrDefault("GOMIDDLEMAN_CERT_FILE", "cert.pem"), "proxyutils certificate as PEM file")
	keyFile      = flag.String("key", getEnvOrDefault("GOMIDDLEMAN_KEY_FILE", "key.pem"), "proxyutils key as PEM file")
	caFile       = flag.String("ca", getEnvOrDefault("GOMIDDLEMAN_CA_FILE", "ca.pem"), "proxyutils ca as PEM file")
)

func main() {

	flag.Parse()

	var wg sync.WaitGroup

	listenAddr := fmt.Sprintf(":%s", *listenPort)

	tlsConfig := tlsutils.LoadTLSConfig(*certFile, *keyFile, *caFile)

	listener, err := listeners.NewListener(*listenerType, listenAddr, tlsConfig)
	if err != nil {
		log.Fatalf("Failed to initialize proxyutils: %v", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Fatalf("Failed to close listener: %v", err)
		}
	}()

	connector, err := connectors.NewConnector(*target, tlsConfig)
	if err != nil {
		log.Fatalf("Failed to initialize target connector: %v", err)
	}

	if err := gomiddleman.StartProxy(listener, connector); err != nil {
		log.Fatalf("Error when starting the proxyutils: %v", err)
	}

	// Setup signal handling for graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	log.Println("Shutting down proxyutils...")
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
