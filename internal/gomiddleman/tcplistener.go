// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package gomiddleman

import (
	"fmt"
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

// Logic to accept a new connection
func (l *TCPListener) Listen(handleConnection func(net.Conn)) error {

	var err error
	l.listener, err = net.Listen("tcp", l.GetAddress())
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", l.GetAddress(), err)
	}
	log.Printf("Listening on %s", l.GetAddress())

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

// Logic to close the listener
func (l *TCPListener) Close() error {
	log.Printf("Shutting down TCP listener at %s", l.GetAddress())
	if err := l.listener.Close(); err != nil {
		return fmt.Errorf("failed to close TCP listener on %s: %v", l.GetAddress(), err)
	}
	l.wg.Wait()
	return nil
}

func (l *TCPListener) NewConnectionHandler() ConnectionHandler {
	return NewTCPConnectionHandler()
}

func (l *TCPListener) GetAddress() string {
	return l.addr
}
