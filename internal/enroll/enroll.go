package enroll

import (
	"context"
	"crypto/tls"
	"net/http"
)

// EnrollmentResult holds everything needed to open an mTLS connection to
// the TAK server after a successful enrollment.
type EnrollmentResult struct {
	PrivateKeyPEM []byte
	ClientCertPEM []byte
	CACertsPEM    [][]byte
}

// DefaultBaseURL builds the enrollment base URL for a TAK server address,
// using the server's standard certificate-enrollment port.
func DefaultBaseURL(serverAddress string) string {
	return "https://" + serverAddress + ":8446"
}

// InsecureHTTPClient returns an HTTP client that skips TLS verification.
// This is required for enrollment: the client has no preconfigured trust
// store, so it authenticates with a username and password instead of
// validating the server's certificate. Trust is established afterwards,
// once the server hands back its CA chain.
func InsecureHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // intentional: no preconfigured trust store exists yet
		},
	}
}

// Enroll runs the full client-certificate enrollment flow against
// baseURL: it fetches the CSR subject fields, generates a key pair and
// CSR for username, and submits it for signing using username/password
// authentication.
func Enroll(ctx context.Context, httpClient *http.Client, baseURL, username, password string) (EnrollmentResult, error) {
	tlsConfig, err := FetchTLSConfig(ctx, httpClient, baseURL, username, password)
	if err != nil {
		return EnrollmentResult{}, err
	}

	kp, err := NewCSR(username, tlsConfig.Organization, tlsConfig.OrganizationalUnit)
	if err != nil {
		return EnrollmentResult{}, err
	}

	signed, err := SignCSR(ctx, httpClient, baseURL, username, password, kp.CSRPEM, "gotak-"+username)
	if err != nil {
		return EnrollmentResult{}, err
	}

	return EnrollmentResult{
		PrivateKeyPEM: kp.PrivateKeyPEM,
		ClientCertPEM: signed.ClientCertPEM,
		CACertsPEM:    signed.CACertsPEM,
	}, nil
}
