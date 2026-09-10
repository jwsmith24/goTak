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

func TestParseFlags_ScenarioPathIsOptional(t *testing.T) {
	args := []string{"-server", "192.168.1.50", "-username", "alice", "-password", "s3cret"}

	cfg, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScenarioPath != "" {
		t.Errorf("ScenarioPath = %q, want empty when not provided", cfg.ScenarioPath)
	}
}

func TestParseFlags_ScenarioPathFromFlag(t *testing.T) {
	args := []string{
		"-server", "192.168.1.50",
		"-username", "alice",
		"-password", "s3cret",
		"-scenario", "scenarios/austin-capitol.json",
	}

	cfg, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScenarioPath != "scenarios/austin-capitol.json" {
		t.Errorf("ScenarioPath = %q, want %q", cfg.ScenarioPath, "scenarios/austin-capitol.json")
	}
	if !cfg.ScenarioFromFlag {
		t.Error("ScenarioFromFlag = false, want true when -scenario is passed on the command line")
	}
}

func TestParseFlags_ScenarioFromFlagFalseWhenNotProvided(t *testing.T) {
	args := []string{"-server", "192.168.1.50", "-username", "alice", "-password", "s3cret"}

	cfg, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScenarioFromFlag {
		t.Error("ScenarioFromFlag = true, want false when -scenario was never passed")
	}
}

func TestParseFlags_LocationNameIsOptional(t *testing.T) {
	args := []string{"-server", "192.168.1.50", "-username", "alice", "-password", "s3cret"}

	cfg, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LocationName != "" {
		t.Errorf("LocationName = %q, want empty when not provided", cfg.LocationName)
	}
}

func TestParseFlags_LocationNameFromFlag(t *testing.T) {
	args := []string{
		"-server", "192.168.1.50",
		"-username", "alice",
		"-password", "s3cret",
		"-location", "Austin, TX",
	}

	cfg, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LocationName != "Austin, TX" {
		t.Errorf("LocationName = %q, want %q", cfg.LocationName, "Austin, TX")
	}
	if !cfg.LocationFromFlag {
		t.Error("LocationFromFlag = false, want true when -location is passed on the command line")
	}
}

func TestParseFlags_LocationFromFlagFalseWhenNotProvided(t *testing.T) {
	args := []string{"-server", "192.168.1.50", "-username", "alice", "-password", "s3cret"}

	cfg, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LocationFromFlag {
		t.Error("LocationFromFlag = true, want false when -location was never passed")
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
