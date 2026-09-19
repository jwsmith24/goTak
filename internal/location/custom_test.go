package location

import (
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestNew_ValidatesCoordinates(t *testing.T) {
	tests := []struct {
		name    string
		lat     float64
		lon     float64
		wantErr bool
	}{
		{name: "valid", lat: 30, lon: -97},
		{name: "boundaries", lat: 90, lon: -180},
		{name: "latitude high", lat: 90.1, wantErr: true},
		{name: "longitude low", lon: -180.1, wantErr: true},
		{name: "latitude NaN", lat: math.NaN(), wantErr: true},
		{name: "longitude infinity", lon: math.Inf(1), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New("Custom", tt.lat, tt.lon)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New returned %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveAndLoadCustom(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "custom_locations.json")
	first, _ := New("First", 1, 2)
	second, _ := New("Second", 3, 4)
	for _, loc := range []Location{first, second} {
		if err := SaveCustom(path, loc); err != nil {
			t.Fatalf("SaveCustom: %v", err)
		}
	}

	got, err := LoadCustom(path)
	if err != nil {
		t.Fatalf("LoadCustom: %v", err)
	}
	if len(got) != 2 || got[0] != first || got[1] != second {
		t.Fatalf("LoadCustom = %+v, want [%+v %+v]", got, first, second)
	}
}

func TestLoadCustom_MissingFileIsEmpty(t *testing.T) {
	got, err := LoadCustom(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil || len(got) != 0 {
		t.Fatalf("LoadCustom = %+v, %v; want empty, nil", got, err)
	}
}

func TestSaveCustom_RejectsDuplicateNames(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom_locations.json")
	loc, _ := New("Home", 1, 2)
	if err := SaveCustom(path, loc); err != nil {
		t.Fatal(err)
	}
	if err := SaveCustom(path, Location{Name: "Home", Lat: 3, Lon: 4}); err == nil {
		t.Fatal("expected duplicate saved name error")
	}
	if err := SaveCustom(path, Location{Name: All[0].Name, Lat: 3, Lon: 4}); err == nil {
		t.Fatal("expected built-in name collision error")
	}
}

func TestSaveCustom_ReplacesRelaxedFileWithPrivatePermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom_locations.json")
	if err := os.WriteFile(path, []byte("[]"), 0o644); err != nil {
		t.Fatal(err)
	}
	loc, _ := New("Home", 1, 2)
	if err := SaveCustom(path, loc); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("permissions = %o, want 600", got)
	}
}

func TestLoadCustom_RejectsDuplicateAndBuiltInNames(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{name: "duplicate custom", json: `[{"name":"Home","lat":1,"lon":2},{"name":"Home","lat":3,"lon":4}]`},
		{name: "built-in collision", json: `[{"name":"Austin, TX","lat":1,"lon":2}]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "custom_locations.json")
			if err := os.WriteFile(path, []byte(tt.json), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := LoadCustom(path); err == nil {
				t.Fatal("expected duplicate name error")
			}
		})
	}
}
