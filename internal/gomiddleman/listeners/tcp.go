// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package listeners

import (
	"fmt"
	"github.com/hyperifyio/gomiddleman/internal/gomiddleman/connectionhandlers"
	"log"
	"net"
	"sync"
)

type TCPListener struct {
	addr     string
	listener net.Listener
	wg       sync.WaitGroup
}

func NewTCPListener(addr string) *TCPListener {
	return &TCPListener{addr: addr}
}

// Listen Logic to accept a new connection
func (listener *TCPListener) Listen(handleConnection func(net.Conn)) error {

	var err error
	listener.listener, err = net.Listen("tcp", listener.GetAddress())
	if err != nil {
		return fmt.Errorf("[TCPListener.Listen]: failed to listen on %s: %v", listener.GetAddress(), err)
	}
	log.Printf("[TCPListener.Listen]: Listening on %s", listener.GetAddress())

	listener.wg.Add(1)
	go func() {
		defer listener.wg.Done()
		for {
			clientConn, err := listener.listener.Accept()
			if err != nil {
				log.Println("[TCPListener.Listen]: Listener closed, stopping accept loop")
				return
			}

			// Handle the connection in a new goroutine.
			go handleConnection(clientConn)
		}
	}()

	return nil
}

// Close Logic to close the listener
func (listener *TCPListener) Close() error {
	log.Printf("[TCPListener.Close]: Shutting down TCP listener at %s", listener.GetAddress())
	if err := listener.listener.Close(); err != nil {
		return fmt.Errorf("[TCPListener.Close]: failed to close TCP listener on %s: %v", listener.GetAddress(), err)
	}
	listener.wg.Wait()
	return nil
}

func (listener *TCPListener) NewConnectionHandler() connectionhandlers.ConnectionHandler {
	return connectionhandlers.NewTCPConnectionHandler()
}

func (listener *TCPListener) GetAddress() string {
	return listener.addr
}

func (listener *TCPListener) GetType() string {
	return "tcp"
}
