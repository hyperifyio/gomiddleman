package gomiddleman

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

func (h *TLSConnectionHandler) Handle(conn net.Conn) error {
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return fmt.Errorf("failed to cast to TLS connection")
	}

	// It's a good practice to perform a handshake explicitly to check for any TLS errors early
	if err := tlsConn.Handshake(); err != nil {
		return fmt.Errorf("TLS handshake failed: %v", err)
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) > 0 {
		clientCert := state.PeerCertificates[0]
		h.commonName = clientCert.Subject.CommonName
		h.serialNumber = clientCert.SerialNumber.Text(16)
	}

	return nil
}

func (h *TLSConnectionHandler) GetCommonName() string {
	return h.commonName
}

func (h *TLSConnectionHandler) GetSerialNumber() string {
	return h.serialNumber
}
