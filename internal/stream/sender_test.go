package stream

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"testing"
	"time"
)

type testPKI struct {
	caCertPEM     []byte
	serverCertPEM []byte
	serverKeyPEM  []byte
	clientCertPEM []byte
	clientKeyPEM  []byte
}

func newTestPKI(t *testing.T) testPKI {
	t.Helper()

	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating CA key: %v", err)
	}
	caTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "gotak-test-ca"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("creating CA certificate: %v", err)
	}
	caCert, err := x509.ParseCertificate(caDER)
	if err != nil {
		t.Fatalf("parsing CA certificate: %v", err)
	}

	issue := func(commonName string, ips []net.IP, extKeyUsage []x509.ExtKeyUsage) ([]byte, []byte) {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatalf("generating key for %s: %v", commonName, err)
		}
		template := x509.Certificate{
			SerialNumber: big.NewInt(time.Now().UnixNano()),
			Subject:      pkix.Name{CommonName: commonName},
			NotBefore:    time.Now().Add(-time.Hour),
			NotAfter:     time.Now().Add(time.Hour),
			KeyUsage:     x509.KeyUsageDigitalSignature,
			ExtKeyUsage:  extKeyUsage,
			IPAddresses:  ips,
		}
		der, err := x509.CreateCertificate(rand.Reader, &template, caCert, &key.PublicKey, caKey)
		if err != nil {
			t.Fatalf("creating certificate for %s: %v", commonName, err)
		}
		keyDER, err := x509.MarshalECPrivateKey(key)
		if err != nil {
			t.Fatalf("marshaling key for %s: %v", commonName, err)
		}
		certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
		keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
		return certPEM, keyPEM
	}

	serverCertPEM, serverKeyPEM := issue("localhost", []net.IP{net.ParseIP("127.0.0.1")}, []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth})
	clientCertPEM, clientKeyPEM := issue("dev", nil, []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth})

	return testPKI{
		caCertPEM:     pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER}),
		serverCertPEM: serverCertPEM,
		serverKeyPEM:  serverKeyPEM,
		clientCertPEM: clientCertPEM,
		clientKeyPEM:  clientKeyPEM,
	}
}

func startMTLSListener(t *testing.T, pki testPKI) (addr string, received chan []byte) {
	t.Helper()

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(pki.caCertPEM) {
		t.Fatal("failed to add CA cert to pool")
	}

	serverCert, err := tls.X509KeyPair(pki.serverCertPEM, pki.serverKeyPEM)
	if err != nil {
		t.Fatalf("loading server key pair: %v", err)
	}

	listener, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
	})
	if err != nil {
		t.Fatalf("starting TLS listener: %v", err)
	}
	t.Cleanup(func() { _ = listener.Close() })

	received = make(chan []byte, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		received <- buf[:n]
	}()

	return listener.Addr().String(), received
}

func TestSend_DeliversEventOverMTLSConnection(t *testing.T) {
	pki := newTestPKI(t)
	addr, received := startMTLSListener(t, pki)

	sender, err := Dial(context.Background(), addr, pki.clientCertPEM, pki.clientKeyPEM, [][]byte{pki.caCertPEM})
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer sender.Close()

	payload := []byte(`<event uid="gotak-sim-1"/>`)
	if err := sender.Send(payload); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	select {
	case got := <-received:
		if string(got) != string(payload) {
			t.Errorf("server received %q, want %q", got, payload)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for server to receive event")
	}
}

func TestDial_FailsWithoutValidClientCert(t *testing.T) {
	pki := newTestPKI(t)
	addr, _ := startMTLSListener(t, pki)

	otherPKI := newTestPKI(t)

	sender, err := Dial(context.Background(), addr, otherPKI.clientCertPEM, otherPKI.clientKeyPEM, [][]byte{pki.caCertPEM})
	if err != nil {
		// The server rejected the client certificate during the
		// handshake itself.
		return
	}
	defer sender.Close()

	// Under TLS 1.3 the client's handshake can complete before the
	// server has finished validating the client certificate: the
	// rejection then arrives as an alert on the next read/write instead.
	deadline := time.Now().Add(2 * time.Second)
	_ = sender.conn.SetDeadline(deadline)
	if writeErr := sender.Send([]byte(`<event uid="should-be-rejected"/>`)); writeErr != nil {
		return
	}
	buf := make([]byte, 1)
	if _, readErr := sender.conn.Read(buf); readErr == nil {
		t.Fatal("expected the server to reject a client cert from an untrusted CA")
	}
}
