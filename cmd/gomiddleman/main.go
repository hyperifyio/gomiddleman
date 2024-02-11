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
	listenerType = flag.String("type", getEnvOrDefault("GOMIDDLEMAN_TYPE", "tls"), "type of proxy (tcp or tls)")
	listenPort   = flag.String("port", getEnvOrDefault("GOMIDDLEMAN_PORT", "8080"), "port on which the proxy listens")
	target       = flag.String("target", getEnvOrDefault("GOMIDDLEMAN_TARGET", "http://localhost:3000"), "target where to proxy connections")
	certFile     = flag.String("cert", getEnvOrDefault("GOMIDDLEMAN_CERT_FILE", "cert.pem"), "proxy certificate as PEM file")
	keyFile      = flag.String("key", getEnvOrDefault("GOMIDDLEMAN_KEY_FILE", "key.pem"), "proxy key as PEM file")
	caFile       = flag.String("ca", getEnvOrDefault("GOMIDDLEMAN_CA_FILE", "ca.pem"), "proxy ca as PEM file")
	caKeyFile    = flag.String("ca-key", getEnvOrDefault("GOMIDDLEMAN_CA_KEY_FILE", "ca-key.pem"), "proxy ca as PEM file")
)

func main() {

	flag.Parse()

	var wg sync.WaitGroup

	listenAddr := fmt.Sprintf(":%s", *listenPort)

	listenerTlsConfig := tlsutils.LoadTLSConfig(*certFile, *keyFile, *caFile)

	listener, err := listeners.NewListener(*listenerType, listenAddr, listenerTlsConfig)
	if err != nil {
		log.Fatalf("[main]: Failed to initialize proxy: %v", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Fatalf("[main]: Failed to close listener: %v", err)
		}
	}()

	connectorTlsConfig := tlsutils.LoadTLSConfig("", "", *caFile)

	caCert := tlsutils.ReadCACertificateFile(*caFile)
	caPrivateKey := tlsutils.ReadCAKeyFile(*caKeyFile)

	connector, err := connectors.NewConnector(
		*target,
		connectorTlsConfig,
		caCert,
		caPrivateKey,
	)
	if err != nil {
		log.Fatalf("[main]: Failed to initialize target connector: %v", err)
	}

	if err := gomiddleman.StartProxy(listener, connector); err != nil {
		log.Fatalf("[main]: Error when starting the proxy: %v", err)
	}

	// Setup signal handling for graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	log.Println("Shutting down proxy...")
	if err := listener.Close(); err != nil {
		log.Printf("[main]: Failed to close listener: %v", err)
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
