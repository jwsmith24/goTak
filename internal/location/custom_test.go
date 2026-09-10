package location

import (
	"path/filepath"
	"testing"
)

func TestNew_ValidCoordinates(t *testing.T) {
	loc, err := New("Home", 30.0, -97.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.Name != "Home" || loc.Lat != 30.0 || loc.Lon != -97.0 {
		t.Errorf("New() = %+v, unexpected", loc)
	}
}

func TestNew_RejectsLatitudeOutOfRange(t *testing.T) {
	for _, lat := range []float64{90.1, -90.1} {
		if _, err := New("Bad", lat, 0); err == nil {
			t.Errorf("New() with lat %v: expected error, got nil", lat)
		}
	}
}

func TestNew_RejectsLongitudeOutOfRange(t *testing.T) {
	for _, lon := range []float64{180.1, -180.1} {
		if _, err := New("Bad", 0, lon); err == nil {
			t.Errorf("New() with lon %v: expected error, got nil", lon)
		}
	}
}

func TestNew_AcceptsBoundaryCoordinates(t *testing.T) {
	for _, coords := range [][2]float64{{90, 180}, {-90, -180}} {
		if _, err := New("Edge", coords[0], coords[1]); err != nil {
			t.Errorf("New(%v, %v): unexpected error: %v", coords[0], coords[1], err)
		}
	}
}

func TestLoadCustom_MissingFileReturnsEmptySlice(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.json")

	locs, err := LoadCustom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(locs) != 0 {
		t.Errorf("LoadCustom() = %+v, want empty", locs)
	}
}

func TestSaveCustom_ThenLoadCustom_RoundTrips(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "custom_locations.json")

	loc, err := New("Home Base", 30.0, -97.0)
	if err != nil {
		t.Fatalf("unexpected error building location: %v", err)
	}
	if err := SaveCustom(path, loc); err != nil {
		t.Fatalf("SaveCustom: unexpected error: %v", err)
	}

	locs, err := LoadCustom(path)
	if err != nil {
		t.Fatalf("LoadCustom: unexpected error: %v", err)
	}
	if len(locs) != 1 || locs[0] != loc {
		t.Errorf("LoadCustom() = %+v, want [%+v]", locs, loc)
	}
}

func TestSaveCustom_AppendsToExisting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom_locations.json")

	first, _ := New("First", 1, 1)
	second, _ := New("Second", 2, 2)

	if err := SaveCustom(path, first); err != nil {
		t.Fatalf("SaveCustom (first): unexpected error: %v", err)
	}
	if err := SaveCustom(path, second); err != nil {
		t.Fatalf("SaveCustom (second): unexpected error: %v", err)
	}

	locs, err := LoadCustom(path)
	if err != nil {
		t.Fatalf("LoadCustom: unexpected error: %v", err)
	}
	if len(locs) != 2 || locs[0] != first || locs[1] != second {
		t.Errorf("LoadCustom() = %+v, want [%+v, %+v]", locs, first, second)
	}
}

func TestCustomLocationsPath_ReturnsNonEmptyPath(t *testing.T) {
	path, err := CustomLocationsPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path == "" {
		t.Error("CustomLocationsPath() = \"\", want a non-empty path")
	}
}
