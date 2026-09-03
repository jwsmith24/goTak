package enroll

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestNewCSR_ProducesValidRequestWithSubject(t *testing.T) {
	kp, err := NewCSR("alice", "TAK", "TAK")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	keyBlock, _ := pem.Decode(kp.PrivateKeyPEM)
	if keyBlock == nil {
		t.Fatal("PrivateKeyPEM did not contain a PEM block")
	}

	csrBlock, rest := pem.Decode(kp.CSRPEM)
	if csrBlock == nil {
		t.Fatal("CSRPEM did not contain a PEM block")
	}
	if len(rest) != 0 {
		t.Errorf("unexpected trailing data after CSR PEM block: %q", rest)
	}
	if csrBlock.Type != "CERTIFICATE REQUEST" {
		t.Errorf("CSR PEM block type = %q, want %q", csrBlock.Type, "CERTIFICATE REQUEST")
	}

	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to parse CSR: %v", err)
	}

	if err := csr.CheckSignature(); err != nil {
		t.Errorf("CSR signature is invalid: %v", err)
	}

	if csr.Subject.CommonName != "alice" {
		t.Errorf("CommonName = %q, want %q", csr.Subject.CommonName, "alice")
	}
	if len(csr.Subject.Organization) != 1 || csr.Subject.Organization[0] != "TAK" {
		t.Errorf("Organization = %v, want [TAK]", csr.Subject.Organization)
	}
	if len(csr.Subject.OrganizationalUnit) != 1 || csr.Subject.OrganizationalUnit[0] != "TAK" {
		t.Errorf("OrganizationalUnit = %v, want [TAK]", csr.Subject.OrganizationalUnit)
	}
}

func TestNewCSR_RequiresCommonName(t *testing.T) {
	_, err := NewCSR("", "TAK", "TAK")
	if err == nil {
		t.Fatal("expected error for empty common name, got nil")
	}
}
