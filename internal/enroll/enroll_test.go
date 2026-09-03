package enroll

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newFakeTAKServer(t *testing.T, username, password string) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/Marti/api/tls/config", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<ns2:certificateConfig xmlns="http://bbn.com/marti/xml/config" xmlns:ns2="com.bbn.marti.config">` +
			`<nameEntries><nameEntry name="O" value="TAK"/><nameEntry name="OU" value="TAK"/></nameEntries>` +
			`</ns2:certificateConfig>`))
	})
	mux.HandleFunc("/Marti/api/tls/signClient/v2", func(w http.ResponseWriter, r *http.Request) {
		gotUser, gotPass, ok := r.BasicAuth()
		if !ok || gotUser != username || gotPass != password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		resp := map[string]string{
			"signedCert": fixtureCertBase64(t, "dev"),
			"ca0":        fixtureCertBase64(t, "TAK CA"),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	return httptest.NewTLSServer(mux)
}

func TestEnroll_FullFlowProducesUsableCredentials(t *testing.T) {
	srv := newFakeTAKServer(t, "dev", "devpass")
	defer srv.Close()

	result, err := Enroll(context.Background(), srv.Client(), srv.URL, "dev", "devpass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	keyBlock, _ := pem.Decode(result.PrivateKeyPEM)
	if keyBlock == nil {
		t.Fatal("PrivateKeyPEM did not contain a PEM block")
	}
	if _, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes); err != nil {
		t.Errorf("failed to parse private key: %v", err)
	}

	certBlock, _ := pem.Decode(result.ClientCertPEM)
	if certBlock == nil {
		t.Fatal("ClientCertPEM did not contain a PEM block")
	}
	if _, err := x509.ParseCertificate(certBlock.Bytes); err != nil {
		t.Errorf("failed to parse client certificate: %v", err)
	}

	if len(result.CACertsPEM) != 1 {
		t.Fatalf("len(CACertsPEM) = %d, want 1", len(result.CACertsPEM))
	}
}

func TestEnroll_WrongPasswordFails(t *testing.T) {
	srv := newFakeTAKServer(t, "dev", "devpass")
	defer srv.Close()

	_, err := Enroll(context.Background(), srv.Client(), srv.URL, "dev", "wrongpass")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestDefaultBaseURL(t *testing.T) {
	got := DefaultBaseURL("192.168.1.50")
	want := "https://192.168.1.50:8446"
	if got != want {
		t.Errorf("DefaultBaseURL(%q) = %q, want %q", "192.168.1.50", got, want)
	}
}

func TestInsecureHTTPClient_SkipsTLSVerification(t *testing.T) {
	client := InsecureHTTPClient()
	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("Transport is %T, want *http.Transport", client.Transport)
	}
	if transport.TLSClientConfig == nil || !transport.TLSClientConfig.InsecureSkipVerify {
		t.Error("expected InsecureSkipVerify to be true, since enrollment has no preconfigured trust store")
	}
}
