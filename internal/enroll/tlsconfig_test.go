package enroll

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchTLSConfig_ParsesOrganizationFields(t *testing.T) {
	const username = "dev"
	const password = "devpass"

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Marti/api/tls/config" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		gotUser, gotPass, ok := r.BasicAuth()
		if !ok || gotUser != username || gotPass != password {
			t.Errorf("BasicAuth = (%q, %q, %v), want (%q, %q, true)", gotUser, gotPass, ok, username, password)
		}
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		_, _ = w.Write([]byte(`<ns2:certificateConfig xmlns="http://bbn.com/marti/xml/config" xmlns:ns2="com.bbn.marti.config">` +
			`<nameEntries>` +
			`<nameEntry name="O" value="TAK"/>` +
			`<nameEntry name="OU" value="TAK"/>` +
			`</nameEntries>` +
			`</ns2:certificateConfig>`))
	}))
	defer srv.Close()

	cfg, err := FetchTLSConfig(context.Background(), srv.Client(), srv.URL, username, password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Organization != "TAK" {
		t.Errorf("Organization = %q, want %q", cfg.Organization, "TAK")
	}
	if cfg.OrganizationalUnit != "TAK" {
		t.Errorf("OrganizationalUnit = %q, want %q", cfg.OrganizationalUnit, "TAK")
	}
}

func TestFetchTLSConfig_ServerError(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := FetchTLSConfig(context.Background(), srv.Client(), srv.URL, "dev", "devpass")
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestFetchTLSConfig_Unauthorized(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	_, err := FetchTLSConfig(context.Background(), srv.Client(), srv.URL, "dev", "wrongpass")
	if err == nil {
		t.Fatal("expected error for 401 response, got nil")
	}
}
