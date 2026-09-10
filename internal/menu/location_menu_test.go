package menu

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/jwsmith24/goTak/internal/location"
)

// nonTerminalStdin returns an *os.File that is never a terminal (a
// pipe), so choose() deterministically falls back to the plain
// numbered-list prompt regardless of how tests are invoked.
func nonTerminalStdin(t *testing.T) *os.File {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("creating pipe: %v", err)
	}
	t.Cleanup(func() {
		r.Close()
		w.Close()
	})
	return r
}

func TestRunLocationMenu_SelectsAnExistingEntry(t *testing.T) {
	var out bytes.Buffer
	customPath := filepath.Join(t.TempDir(), "custom_locations.json")
	locations := []location.Location{
		{Name: "Austin, TX", Lat: 30.2747, Lon: -97.7404},
		{Name: "Fort Campbell, KY", Lat: 36.6722, Lon: -87.4925},
	}

	got, err := runLocationMenu(nonTerminalStdin(t), bufReader("2\n"), &out, locations, customPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != locations[1] {
		t.Errorf("runLocationMenu() = %+v, want %+v", got, locations[1])
	}
}

func TestRunLocationMenu_ListsACustomLocationEntry(t *testing.T) {
	var out bytes.Buffer
	customPath := filepath.Join(t.TempDir(), "custom_locations.json")
	locations := []location.Location{{Name: "Austin, TX", Lat: 30.2747, Lon: -97.7404}}

	if _, err := runLocationMenu(nonTerminalStdin(t), bufReader("1\n"), &out, locations, customPath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(out.Bytes(), []byte(customLocationLabel)) {
		t.Errorf("printed menu = %q, want it to list %q", out.String(), customLocationLabel)
	}
}

func TestRunLocationMenu_CustomEntryPromptsAndSaves(t *testing.T) {
	var out bytes.Buffer
	customPath := filepath.Join(t.TempDir(), "custom_locations.json")
	locations := []location.Location{{Name: "Austin, TX", Lat: 30.2747, Lon: -97.7404}}

	// "2" selects the custom-location entry (index 1, the only entry
	// after the single named location), then name/lat/lon.
	got, err := runLocationMenu(nonTerminalStdin(t), bufReader("2\nHome Base\n30.0\n-97.0\n"), &out, locations, customPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := location.Location{Name: "Home Base", Lat: 30.0, Lon: -97.0}
	if got != want {
		t.Errorf("runLocationMenu() = %+v, want %+v", got, want)
	}

	saved, err := location.LoadCustom(customPath)
	if err != nil {
		t.Fatalf("LoadCustom: unexpected error: %v", err)
	}
	if len(saved) != 1 || saved[0] != want {
		t.Errorf("saved custom locations = %+v, want [%+v]", saved, want)
	}
}

func TestRunLocationMenu_InvalidCustomCoordinatesReturnsErrorWithoutSaving(t *testing.T) {
	var out bytes.Buffer
	customPath := filepath.Join(t.TempDir(), "custom_locations.json")
	locations := []location.Location{{Name: "Austin, TX", Lat: 30.2747, Lon: -97.7404}}

	_, err := runLocationMenu(nonTerminalStdin(t), bufReader("2\nBad\n120\n-97.0\n"), &out, locations, customPath)
	if err == nil {
		t.Fatal("expected error for out-of-range latitude, got nil")
	}

	if _, statErr := os.Stat(customPath); !os.IsNotExist(statErr) {
		t.Errorf("expected no custom locations file to be created, stat err = %v", statErr)
	}
}
