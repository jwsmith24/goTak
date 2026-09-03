package enroll

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fixtureCertBase64 returns a realistic, PEM-marker-stripped, base64 DER
// certificate, mimicking what the server embeds in its JSON response.
func fixtureCertBase64(t *testing.T, commonName string) string {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating fixture key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: commonName},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("creating fixture certificate: %v", err)
	}

	return base64.StdEncoding.EncodeToString(der)
}

func TestSignCSR_ReturnsSignedCertificateAndCAChain(t *testing.T) {
	const username = "dev"
	const password = "devpass"
	const clientUID = "gotak-sim-1"

	signedCertB64 := fixtureCertBase64(t, "dev")
	caCertB64 := fixtureCertBase64(t, "TAK CA")

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Marti/api/tls/signClient/v2" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if got := r.URL.Query().Get("clientUid"); got != clientUID {
			t.Errorf("clientUid = %q, want %q", got, clientUID)
		}

		gotUser, gotPass, ok := r.BasicAuth()
		if !ok || gotUser != username || gotPass != password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading request body: %v", err)
		}
		if strings.Contains(string(body), "BEGIN CERTIFICATE REQUEST") {
			t.Error("expected PEM markers to be stripped from CSR body")
		}

		resp := map[string]string{
			"signedCert": signedCertB64,
			"ca0":        caCertB64,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	kp, err := NewCSR("dev", "TAK", "TAK")
	if err != nil {
		t.Fatalf("generating CSR: %v", err)
	}

	result, err := SignCSR(context.Background(), srv.Client(), srv.URL, username, password, kp.CSRPEM, clientUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	certBlock, _ := pem.Decode(result.ClientCertPEM)
	if certBlock == nil {
		t.Fatal("ClientCertPEM did not contain a PEM block")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to parse returned client certificate: %v", err)
	}
	if cert.Subject.CommonName != "dev" {
		t.Errorf("client cert CommonName = %q, want %q", cert.Subject.CommonName, "dev")
	}

	if len(result.CACertsPEM) != 1 {
		t.Fatalf("len(CACertsPEM) = %d, want 1", len(result.CACertsPEM))
	}
	caBlock, _ := pem.Decode(result.CACertsPEM[0])
	if caBlock == nil {
		t.Fatal("CACertsPEM[0] did not contain a PEM block")
	}
	caCert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to parse returned CA certificate: %v", err)
	}
	if caCert.Subject.CommonName != "TAK CA" {
		t.Errorf("CA cert CommonName = %q, want %q", caCert.Subject.CommonName, "TAK CA")
	}
}

func TestSignCSR_DeduplicatesRepeatedCACerts(t *testing.T) {
	caCertB64 := fixtureCertBase64(t, "TAK CA")

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"signedCert": fixtureCertBase64(t, "dev"),
			"ca0":        caCertB64,
			"ca1":        caCertB64,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	kp, err := NewCSR("dev", "TAK", "TAK")
	if err != nil {
		t.Fatalf("generating CSR: %v", err)
	}

	result, err := SignCSR(context.Background(), srv.Client(), srv.URL, "dev", "devpass", kp.CSRPEM, "uid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.CACertsPEM) != 1 {
		t.Errorf("len(CACertsPEM) = %d, want 1 (duplicates deduplicated)", len(result.CACertsPEM))
	}
}

func TestSignCSR_Unauthorized(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	kp, err := NewCSR("dev", "TAK", "TAK")
	if err != nil {
		t.Fatalf("generating CSR: %v", err)
	}

	_, err = SignCSR(context.Background(), srv.Client(), srv.URL, "dev", "wrongpass", kp.CSRPEM, "uid")
	if err == nil {
		t.Fatal("expected error for 401 response, got nil")
	}
}
