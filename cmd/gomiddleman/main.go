// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package main

import (
	"flag"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman"
	"os"
	"os/signal"
	"syscall"
)

var (
	listenPort = flag.String("port", getEnvOrDefault("GOMIDDLEMAN_PORT", "8080"), "port on which the proxy listens")
	target     = flag.String("target", getEnvOrDefault("GOMIDDLEMAN_TARGET", "127.0.0.1:3000"), "target where to proxy connections")
	certFile   = flag.String("cert", getEnvOrDefault("GOMIDDLEMAN_CERT_FILE", "cert.pem"), "proxy certificate as PEM file")
	keyFile    = flag.String("key", getEnvOrDefault("GOMIDDLEMAN_KEY_FILE", "key.pem"), "proxy key as PEM file")
	caFile     = flag.String("ca", getEnvOrDefault("GOMIDDLEMAN_CA_FILE", "ca.pem"), "proxy ca as PEM file")
)

func main() {

	flag.Parse()

	tlsConfig := gomiddleman.LoadTLSConfig(*certFile, *keyFile, *caFile)

	stopProxy := gomiddleman.StartProxy(*listenPort, *target, tlsConfig)
	defer stopProxy()

	// Setup signal handling for graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	stopProxy()

}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}