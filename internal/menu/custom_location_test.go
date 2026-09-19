package menu

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jwsmith24/goTak/internal/location"
)

func nonTerminalStdin(t *testing.T) *os.File {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("creating pipe: %v", err)
	}
	t.Cleanup(func() {
		_ = r.Close()
		_ = w.Close()
	})
	return r
}

func TestRunLocationMenu_CustomEntryPromptsAndSaves(t *testing.T) {
	var out bytes.Buffer
	path := filepath.Join(t.TempDir(), "custom_locations.json")
	locations := []location.Location{{Name: "Austin", Lat: 30, Lon: -97}}

	got, err := runLocationMenu(nonTerminalStdin(t), bufReader("2\nHome\n31\n-98\n"), &out, locations, path)
	if err != nil {
		t.Fatalf("runLocationMenu: %v", err)
	}
	want := location.Location{Name: "Home", Lat: 31, Lon: -98}
	if got != want {
		t.Fatalf("location = %+v, want %+v", got, want)
	}
	saved, err := location.LoadCustom(path)
	if err != nil || len(saved) != 1 || saved[0] != want {
		t.Fatalf("saved = %+v, %v; want [%+v]", saved, err, want)
	}
	for _, text := range []string{customLocationLabel, "Name", "Latitude", "Longitude"} {
		if !strings.Contains(out.String(), text) {
			t.Errorf("output %q does not contain %q", out.String(), text)
		}
	}
}

func TestRunLocationMenu_CustomEntryRejectsInvalidCoordinates(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom_locations.json")
	_, err := runLocationMenu(nonTerminalStdin(t), bufReader("2\nBad\n91\n0\n"), &bytes.Buffer{}, []location.Location{{Name: "Austin"}}, path)
	if err == nil {
		t.Fatal("expected invalid coordinate error")
	}
}

func TestRunLocationMenu_InteractiveCustomEntryPromptsAndSaves(t *testing.T) {
	var out bytes.Buffer
	path := filepath.Join(t.TempDir(), "custom_locations.json")
	input := bufReader("\x1b[B\nHome\n31\n-98\n")
	interactive := func(_ *os.File, r *bufio.Reader, w io.Writer, title string, entries []string) (int, error) {
		return chooseInteractive(w, r, title, entries)
	}

	got, err := runLocationMenuWithChooser(nonTerminalStdin(t), input, &out, []location.Location{{Name: "Austin"}}, path, interactive)
	if err != nil {
		t.Fatalf("runLocationMenuWithChooser: %v", err)
	}
	if got.Name != "Home" || got.Lat != 31 || got.Lon != -98 {
		t.Fatalf("location = %+v", got)
	}
}
