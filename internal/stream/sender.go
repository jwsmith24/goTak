// Package stream sends Cursor-on-Target events to a TAK server over an
// mTLS connection, using the client certificate obtained during
// enrollment.
package stream

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
)

// Sender is a persistent mTLS connection to a TAK server's CoT streaming
// port, over which CoT XML events are written.
type Sender struct {
	conn net.Conn
}

// Dial opens an mTLS connection to addr (host:port, typically the
// server's CoT streaming port), authenticating with the client
// certificate/key pair issued during enrollment and verifying the server
// against the CA chain returned by that same enrollment.
func Dial(ctx context.Context, addr string, clientCertPEM, clientKeyPEM []byte, caCertsPEM [][]byte) (*Sender, error) {
	// The server may only trust a root CA while the client certificate
	// was issued by an intermediate; without the intermediate in the
	// handshake, the server can fail to build a trust chain even though
	// Dial itself appears to succeed (TLS 1.3 completes the client's
	// side of the handshake before the server finishes verifying it).
	// Presenting the full chain (leaf followed by every CA certificate
	// enrollment returned) covers both single-CA and root+intermediate
	// deployments.
	certChainPEM := bytes.Join(append([][]byte{clientCertPEM}, caCertsPEM...), []byte("\n"))
	cert, err := tls.X509KeyPair(certChainPEM, clientKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("stream: loading client certificate chain: %w", err)
	}

	rootCAs := x509.NewCertPool()
	for _, caCertPEM := range caCertsPEM {
		if !rootCAs.AppendCertsFromPEM(caCertPEM) {
			return nil, fmt.Errorf("stream: failed to parse CA certificate")
		}
	}

	dialer := &tls.Dialer{
		Config: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      rootCAs,
		},
	}

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("stream: dialing %s: %w", addr, err)
	}

	return &Sender{conn: conn}, nil
}

// Send writes a single CoT XML event to the connection.
func (s *Sender) Send(eventXML []byte) error {
	_, err := s.conn.Write(eventXML)
	return err
}

// Close closes the underlying connection.
func (s *Sender) Close() error {
	return s.conn.Close()
}
