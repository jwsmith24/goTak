package config

import (
	"strings"
	"testing"
)

func TestParseFlags_AllFieldsProvided(t *testing.T) {
	args := []string{
		"-server", "192.168.1.50",
		"-username", "alice",
		"-password", "s3cret",
	}

	cfg, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.ServerAddress != "192.168.1.50" {
		t.Errorf("ServerAddress = %q, want %q", cfg.ServerAddress, "192.168.1.50")
	}
	if cfg.Username != "alice" {
		t.Errorf("Username = %q, want %q", cfg.Username, "alice")
	}
	if cfg.Password != "s3cret" {
		t.Errorf("Password = %q, want %q", cfg.Password, "s3cret")
	}
}

func TestParseFlags_MissingServerAddress(t *testing.T) {
	args := []string{"-username", "alice", "-password", "s3cret"}

	_, err := ParseFlags(args)
	if err == nil {
		t.Fatal("expected error for missing server address, got nil")
	}
	if !strings.Contains(err.Error(), "server") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "server")
	}
}

func TestParseFlags_MissingUsername(t *testing.T) {
	args := []string{"-server", "192.168.1.50", "-password", "s3cret"}

	_, err := ParseFlags(args)
	if err == nil {
		t.Fatal("expected error for missing username, got nil")
	}
	if !strings.Contains(err.Error(), "username") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "username")
	}
}

func TestParseFlags_MissingPassword(t *testing.T) {
	args := []string{"-server", "192.168.1.50", "-username", "alice"}

	_, err := ParseFlags(args)
	if err == nil {
		t.Fatal("expected error for missing password, got nil")
	}
	if !strings.Contains(err.Error(), "password") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "password")
	}
}

func TestParseFlags_MissingAllRequiredFields(t *testing.T) {
	_, err := ParseFlags([]string{})
	if err == nil {
		t.Fatal("expected error for missing fields, got nil")
	}
	for _, field := range []string{"server", "username", "password"} {
		if !strings.Contains(err.Error(), field) {
			t.Errorf("error = %q, want it to mention %q", err.Error(), field)
		}
	}
}
