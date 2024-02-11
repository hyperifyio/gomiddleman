package connectionhandlers

import (
	"crypto/tls"
	"fmt"
	"net"
)

type TLSConnectionHandler struct {
	commonName   string
	serialNumber string
}

func NewTLSConnectionHandler() *TLSConnectionHandler {
	return &TLSConnectionHandler{}
}

func (handle *TLSConnectionHandler) Handle(conn net.Conn) error {
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return fmt.Errorf("[TLSConnectionHandler.Handle]: failed to cast to TLS connection")
	}

	// It's a good practice to perform a handshake explicitly to check for any TLS errors early
	if err := tlsConn.Handshake(); err != nil {
		return fmt.Errorf("[TLSConnectionHandler.Handle]: TLS handshake failed: %v", err)
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) > 0 {
		clientCert := state.PeerCertificates[0]
		handle.commonName = clientCert.Subject.CommonName
		handle.serialNumber = clientCert.SerialNumber.Text(16)
	}

	return nil
}

func (handle *TLSConnectionHandler) GetCommonName() string {
	return handle.commonName
}

func (handle *TLSConnectionHandler) GetSerialNumber() string {
	return handle.serialNumber
}
