package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvFile(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("writing temp env file: %v", err)
	}
	return path
}

func TestLoadEnvFile_ParsesKeyValuePairs(t *testing.T) {
	path := writeTempEnvFile(t, ""+
		"GOTAK_SERVER=192.168.1.50\n"+
		"# a comment\n"+
		"\n"+
		"GOTAK_USERNAME=dev\n"+
		`GOTAK_PASSWORD="s3cret value"`+"\n")

	values, err := LoadEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[string]string{
		"GOTAK_SERVER":   "192.168.1.50",
		"GOTAK_USERNAME": "dev",
		"GOTAK_PASSWORD": "s3cret value",
	}
	for key, wantVal := range want {
		if got := values[key]; got != wantVal {
			t.Errorf("values[%q] = %q, want %q", key, got, wantVal)
		}
	}
}

func TestLoadEnvFile_MissingFileReturnsError(t *testing.T) {
	_, err := LoadEnvFile(filepath.Join(t.TempDir(), "does-not-exist.env"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
	if !os.IsNotExist(err) {
		t.Errorf("expected a not-exist error, got %v", err)
	}
}

func TestLoadEnvFile_InvalidLineReturnsError(t *testing.T) {
	path := writeTempEnvFile(t, "not-a-valid-line\n")

	_, err := LoadEnvFile(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}
