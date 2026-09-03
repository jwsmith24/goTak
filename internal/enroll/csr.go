// Package enroll implements the TAK Server client-certificate enrollment
// flow: a client authenticates with a username and password (no
// preconfigured trust store) to obtain a signed client certificate and the
// server's CA chain, which are then used for subsequent mTLS connections.
package enroll

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
)

// KeyPair holds a freshly generated private key and the CSR built from it.
type KeyPair struct {
	PrivateKeyPEM []byte
	CSRPEM        []byte
}

// NewCSR generates an RSA key pair and a PKCS#10 certificate signing
// request for commonName, using organization and organizationalUnit as
// reported by the server's /Marti/api/tls/config endpoint.
func NewCSR(commonName, organization, organizationalUnit string) (KeyPair, error) {
	if commonName == "" {
		return KeyPair{}, errors.New("enroll: common name is required")
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return KeyPair{}, err
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:         commonName,
			Organization:       []string{organization},
			OrganizationalUnit: []string{organizationalUnit},
		},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, key)
	if err != nil {
		return KeyPair{}, err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDER,
	})

	return KeyPair{PrivateKeyPEM: keyPEM, CSRPEM: csrPEM}, nil
}
