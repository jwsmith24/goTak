package scenario

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const twoTrackJSON = `{
	"tickIntervalSeconds": 3,
	"tracks": [
		{
			"uid": "track-1",
			"callsign": "EAGLE01",
			"type": "a-f-A",
			"lat": 30.2747,
			"lon": -97.76,
			"hae": 1500,
			"courseDeg": 90,
			"speedMps": 120
		},
		{
			"uid": "track-2",
			"callsign": "EAGLE02",
			"lat": 30.26,
			"lon": -97.7404,
			"hae": 2000,
			"courseDeg": 0,
			"speedMps": 100
		}
	]
}`

func TestParse_ValidScenario(t *testing.T) {
	sc, err := Parse([]byte(twoTrackJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sc.Tracks) != 2 {
		t.Fatalf("len(Tracks) = %d, want 2", len(sc.Tracks))
	}

	first := sc.Tracks[0]
	if first.UID != "track-1" || first.Callsign != "EAGLE01" || first.Type != "a-f-A" {
		t.Errorf("first track = %+v, unexpected fields", first)
	}
	if first.Lat != 30.2747 || first.Lon != -97.76 || first.HAE != 1500 {
		t.Errorf("first track position = %+v, unexpected", first)
	}
	if first.CourseDeg != 90 || first.SpeedMPS != 120 {
		t.Errorf("first track kinematics = %+v, unexpected", first)
	}

	second := sc.Tracks[1]
	if second.UID != "track-2" || second.Callsign != "EAGLE02" {
		t.Errorf("second track = %+v, unexpected fields", second)
	}
	if second.Type != "" {
		t.Errorf("second track Type = %q, want empty (defaults applied later by cot package)", second.Type)
	}

	if sc.TickInterval() != 3*time.Second {
		t.Errorf("TickInterval() = %v, want 3s", sc.TickInterval())
	}
}

func TestParse_DefaultsTickIntervalWhenOmitted(t *testing.T) {
	sc, err := Parse([]byte(`{"tracks": [{"uid": "t1", "callsign": "C1", "lat": 1, "lon": 2}]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.TickInterval() != 2*time.Second {
		t.Errorf("TickInterval() = %v, want default 2s", sc.TickInterval())
	}
}

func TestParse_RequiresAtLeastOneTrack(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": []}`))
	if err == nil {
		t.Fatal("expected error for empty tracks list, got nil")
	}
}

func TestParse_RequiresUID(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{"callsign": "C1", "lat": 1, "lon": 2}]}`))
	if err == nil {
		t.Fatal("expected error for missing uid, got nil")
	}
	if !strings.Contains(err.Error(), "uid") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "uid")
	}
}

func TestParse_RequiresCallsign(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{"uid": "t1", "lat": 1, "lon": 2}]}`))
	if err == nil {
		t.Fatal("expected error for missing callsign, got nil")
	}
	if !strings.Contains(err.Error(), "callsign") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "callsign")
	}
}

func TestParse_RejectsDuplicateUID(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [
		{"uid": "t1", "callsign": "C1", "lat": 1, "lon": 2},
		{"uid": "t1", "callsign": "C2", "lat": 3, "lon": 4}
	]}`))
	if err == nil {
		t.Fatal("expected error for duplicate uid, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "duplicate")
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := Parse([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestLoad_ReadsFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "scenario.json")
	if err := os.WriteFile(path, []byte(twoTrackJSON), 0o600); err != nil {
		t.Fatalf("writing scenario file: %v", err)
	}

	sc, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sc.Tracks) != 2 {
		t.Fatalf("len(Tracks) = %d, want 2", len(sc.Tracks))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "does-not-exist.json"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
