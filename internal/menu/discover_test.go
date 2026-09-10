package menu

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, dir, name, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o600); err != nil {
		t.Fatalf("writing %s: %v", name, err)
	}
}

const validOneTrack = `{
	"description": "A single test track",
	"tracks": [{"uid": "t1", "callsign": "C1", "lat": 1, "lon": 2}]
}`

const validTwoTracks = `{
	"tracks": [
		{"uid": "t1", "callsign": "C1", "lat": 1, "lon": 2},
		{"uid": "t2", "callsign": "C2", "lat": 3, "lon": 4}
	]
}`

func TestDiscoverScenarios_ListsJSONFilesSortedWithLabels(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "b-two.json", validTwoTracks)
	writeFile(t, dir, "a-one.json", validOneTrack)

	options := DiscoverScenarios(dir)

	if len(options) != 2 {
		t.Fatalf("len(options) = %d, want 2", len(options))
	}

	if options[0].Path != filepath.Join(dir, "a-one.json") {
		t.Errorf("options[0].Path = %q, want the a-one.json path", options[0].Path)
	}
	if options[1].Path != filepath.Join(dir, "b-two.json") {
		t.Errorf("options[1].Path = %q, want the b-two.json path", options[1].Path)
	}

	if options[0].Label == "" {
		t.Error("options[0].Label is empty")
	}
	// The description, when present, should show up in the label.
	if got := options[0].Label; !strings.Contains(got, "A single test track") {
		t.Errorf("options[0].Label = %q, want it to contain the description", got)
	}
}

func TestDiscoverScenarios_SkipsInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "good.json", validOneTrack)
	writeFile(t, dir, "bad.json", `not valid json`)

	options := DiscoverScenarios(dir)

	if len(options) != 1 {
		t.Fatalf("len(options) = %d, want 1 (invalid file skipped)", len(options))
	}
	if options[0].Path != filepath.Join(dir, "good.json") {
		t.Errorf("options[0].Path = %q, want good.json", options[0].Path)
	}
}

func TestDiscoverScenarios_IgnoresNonJSONFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "good.json", validOneTrack)
	writeFile(t, dir, "README.md", "not a scenario")

	options := DiscoverScenarios(dir)

	if len(options) != 1 {
		t.Fatalf("len(options) = %d, want 1", len(options))
	}
}

func TestDiscoverScenarios_MissingDirReturnsEmpty(t *testing.T) {
	options := DiscoverScenarios(filepath.Join(t.TempDir(), "does-not-exist"))
	if len(options) != 0 {
		t.Errorf("len(options) = %d, want 0 for a missing directory", len(options))
	}
}
