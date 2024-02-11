// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package listeners

import (
	"crypto/tls"
	"fmt"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/connectionhandlers"
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

// Listen Logic to accept a new TLS connection
func (listener *TLSListener) Listen(handleConnection func(net.Conn)) error {

	var err error
	listener.listener, err = tls.Listen("tcp", listener.GetAddress(), listener.config)
	if err != nil {
		return fmt.Errorf("[TLSListener.Close]: failed to listen on %s: %v", listener.GetAddress(), err)
	}
	log.Printf("[TLSListener.Close]: Listening on %s", listener.addr)

	listener.wg.Add(1)
	go func() {
		defer listener.wg.Done()
		for {
			clientConn, err := listener.listener.Accept()
			if err != nil {
				log.Println("[TLSListener.Close]: Listener closed, stopping accept loop")
				return
			}

			// Handle the connection in a new goroutine.
			go handleConnection(clientConn)
		}
	}()

	return nil
}

// Close Logic to close the TLS listener
func (listener *TLSListener) Close() error {
	log.Printf("[TLSListener.Close]: Shutting down TLS listener at %s", listener.GetAddress())
	if err := listener.listener.Close(); err != nil {
		return fmt.Errorf("[TLSListener.Close]: failed to close TLS listener on %s: %v", listener.GetAddress(), err)
	}
	listener.wg.Wait()
	return nil
}

func (listener *TLSListener) NewConnectionHandler() connectionhandlers.ConnectionHandler {
	return connectionhandlers.NewTLSConnectionHandler()
}

func (listener *TLSListener) GetAddress() string {
	return listener.addr
}

func (listener *TLSListener) GetType() string {
	return "tls"
}
