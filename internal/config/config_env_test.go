package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFlagsWithEnvFile_FallsBackToEnvFileValues(t *testing.T) {
	envPath := writeTempEnvFile(t, ""+
		"GOTAK_SERVER=192.168.1.50\n"+
		"GOTAK_USERNAME=dev\n"+
		"GOTAK_PASSWORD=devpass\n")

	cfg, err := parseFlags(nil, envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.ServerAddress != "192.168.1.50" {
		t.Errorf("ServerAddress = %q, want %q", cfg.ServerAddress, "192.168.1.50")
	}
	if cfg.Username != "dev" {
		t.Errorf("Username = %q, want %q", cfg.Username, "dev")
	}
	if cfg.Password != "devpass" {
		t.Errorf("Password = %q, want %q", cfg.Password, "devpass")
	}
}

func TestParseFlagsWithEnvFile_FlagsOverrideEnvFileValues(t *testing.T) {
	envPath := writeTempEnvFile(t, ""+
		"GOTAK_SERVER=192.168.1.50\n"+
		"GOTAK_USERNAME=fileuser\n"+
		"GOTAK_PASSWORD=filepass\n")

	cfg, err := parseFlags([]string{"-username", "cliuser"}, envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Username != "cliuser" {
		t.Errorf("Username = %q, want flag value %q to win over env file", cfg.Username, "cliuser")
	}
	if cfg.ServerAddress != "192.168.1.50" {
		t.Errorf("ServerAddress = %q, want env file fallback %q", cfg.ServerAddress, "192.168.1.50")
	}
}

func TestParseFlagsWithEnvFile_MissingEnvFileIsNotAnError(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "does-not-exist.env")

	cfg, err := parseFlags([]string{"-server", "10.0.0.1", "-username", "dev", "-password", "devpass"}, missingPath)
	if err != nil {
		t.Fatalf("unexpected error when env file is absent: %v", err)
	}
	if cfg.ServerAddress != "10.0.0.1" {
		t.Errorf("ServerAddress = %q, want %q", cfg.ServerAddress, "10.0.0.1")
	}
}

func TestParseFlagsWithEnvFile_InvalidEnvFilePropagatesError(t *testing.T) {
	envPath := writeTempEnvFile(t, "not-a-valid-line\n")

	_, err := parseFlags([]string{"-server", "10.0.0.1", "-username", "dev", "-password", "devpass"}, envPath)
	if err == nil {
		t.Fatal("expected error for a malformed env file, got nil")
	}
}

func TestParseFlags_UsesDotEnvInCurrentDirectory(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("GOTAK_SERVER=10.0.0.9\nGOTAK_USERNAME=dev\nGOTAK_PASSWORD=devpass\n"), 0o600); err != nil {
		t.Fatalf("writing .env: %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("changing to temp dir: %v", err)
	}
	defer func() { _ = os.Chdir(originalDir) }()

	cfg, err := ParseFlags(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ServerAddress != "10.0.0.9" {
		t.Errorf("ServerAddress = %q, want %q", cfg.ServerAddress, "10.0.0.9")
	}
}
