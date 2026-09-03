package enroll

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

// SignResult holds the PEM-encoded certificates returned by the server
// once a CSR has been signed.
type SignResult struct {
	ClientCertPEM []byte
	CACertsPEM    [][]byte
}

var caFieldPattern = regexp.MustCompile(`^ca\d+$`)

// SignCSR submits csrPEM to baseURL's /Marti/api/tls/signClient/v2
// endpoint, authenticating with username and password rather than a
// preconfigured trust store. httpClient is expected to skip TLS
// verification for this call, since the client has no way to validate the
// server's certificate until enrollment completes.
func SignCSR(ctx context.Context, httpClient *http.Client, baseURL, username, password string, csrPEM []byte, clientUID string) (SignResult, error) {
	endpoint := baseURL + "/Marti/api/tls/signClient/v2?" + url.Values{"clientUid": {clientUID}}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(stripPEMMarkers(csrPEM)))
	if err != nil {
		return SignResult{}, err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return SignResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return SignResult{}, fmt.Errorf("enroll: signing CSR: authentication failed for user %q", username)
	}
	if resp.StatusCode != http.StatusOK {
		return SignResult{}, fmt.Errorf("enroll: signing CSR: unexpected status %d", resp.StatusCode)
	}

	var fields map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&fields); err != nil {
		return SignResult{}, fmt.Errorf("enroll: parsing sign response: %w", err)
	}

	signedCert, ok := fields["signedCert"]
	if !ok {
		return SignResult{}, fmt.Errorf("enroll: sign response missing %q field", "signedCert")
	}

	var caKeys []string
	for key := range fields {
		if caFieldPattern.MatchString(key) {
			caKeys = append(caKeys, key)
		}
	}
	sort.Strings(caKeys)

	seen := make(map[string]bool, len(caKeys))
	var caCerts [][]byte
	for _, key := range caKeys {
		value := fields[key]
		if seen[value] {
			continue
		}
		seen[value] = true
		caCerts = append(caCerts, base64ToPEM(value))
	}

	return SignResult{
		ClientCertPEM: base64ToPEM(signedCert),
		CACertsPEM:    caCerts,
	}, nil
}

// stripPEMMarkers returns the base64 payload of a PEM block with its
// header/footer lines and newlines removed, matching what the TAK Server
// v2 enrollment endpoint expects in the request body.
func stripPEMMarkers(pemBytes []byte) []byte {
	var b strings.Builder
	for _, line := range strings.Split(string(pemBytes), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "-----") {
			continue
		}
		b.WriteString(line)
	}
	return []byte(b.String())
}

// base64ToPEM wraps a bare base64 certificate payload, as returned by the
// server, in standard PEM certificate markers.
func base64ToPEM(base64Body string) []byte {
	return []byte("-----BEGIN CERTIFICATE-----\n" + base64Body + "\n-----END CERTIFICATE-----\n")
}
