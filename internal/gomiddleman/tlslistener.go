// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sync"
)

type TLSListener struct {
	addr     string
	config   *tls.Config
	listener net.Listener
	wg       sync.WaitGroup
}

func NewTLSListener(addr string, config *tls.Config) *TLSListener {
	return &TLSListener{addr: addr, config: config}
}

// Logic to accept a new TLS connection
func (l *TLSListener) Listen(handleConnection func(net.Conn)) error {

	var err error
	l.listener, err = tls.Listen("tcp", l.GetAddress(), l.config)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", l.GetAddress(), err)
	}
	log.Printf("Listening on %s", l.addr)

	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for {
			clientConn, err := l.listener.Accept()
			if err != nil {
				log.Println("Listener closed, stopping accept loop")
				return
			}

			// Handle the connection in a new goroutine.
			go handleConnection(clientConn)
		}
	}()

	return nil
}

// Logic to close the TLS listener
func (l *TLSListener) Close() error {
	log.Printf("Shutting down TLS listener at %s", l.GetAddress())
	if err := l.listener.Close(); err != nil {
		return fmt.Errorf("failed to close TLS listener on %s: %v", l.GetAddress(), err)
	}
	l.wg.Wait()
	return nil
}

func (l *TLSListener) NewConnectionHandler() ConnectionHandler {
	return NewTLSConnectionHandler()
}

func (l *TLSListener) GetAddress() string {
	return l.addr
}
