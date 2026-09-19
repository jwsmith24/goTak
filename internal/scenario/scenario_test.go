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
			"offsetNorthMeters": 30.2747,
			"offsetEastMeters": -97.76,
			"hae": 1500,
			"courseDeg": 90,
			"speedMps": 120
		},
		{
			"uid": "track-2",
			"callsign": "EAGLE02",
			"offsetNorthMeters": 30.26,
			"offsetEastMeters": -97.7404,
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
	if first.OffsetNorthMeters != 30.2747 || first.OffsetEastMeters != -97.76 || first.HAE != 1500 {
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

func TestParse_ParsesDescription(t *testing.T) {
	sc, err := Parse([]byte(`{"description": "Two tracks near the Capitol", "tracks": [{"uid": "t1", "callsign": "C1", "offsetNorthMeters": 1, "offsetEastMeters": 2}]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Description != "Two tracks near the Capitol" {
		t.Errorf("Description = %q, want %q", sc.Description, "Two tracks near the Capitol")
	}
}

func TestParse_DescriptionIsOptional(t *testing.T) {
	sc, err := Parse([]byte(twoTrackJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Description != "" {
		t.Errorf("Description = %q, want empty when not set", sc.Description)
	}
}

func TestParse_DefaultsTickIntervalWhenOmitted(t *testing.T) {
	sc, err := Parse([]byte(`{"tracks": [{"uid": "t1", "callsign": "C1", "offsetNorthMeters": 1, "offsetEastMeters": 2}]}`))
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
	_, err := Parse([]byte(`{"tracks": [{"callsign": "C1", "offsetNorthMeters": 1, "offsetEastMeters": 2}]}`))
	if err == nil {
		t.Fatal("expected error for missing uid, got nil")
	}
	if !strings.Contains(err.Error(), "uid") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "uid")
	}
}

func TestParse_RequiresCallsign(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{"uid": "t1", "offsetNorthMeters": 1, "offsetEastMeters": 2}]}`))
	if err == nil {
		t.Fatal("expected error for missing callsign, got nil")
	}
	if !strings.Contains(err.Error(), "callsign") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "callsign")
	}
}

func TestParse_RejectsDuplicateUID(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [
		{"uid": "t1", "callsign": "C1", "offsetNorthMeters": 1, "offsetEastMeters": 2},
		{"uid": "t1", "callsign": "C2", "offsetNorthMeters": 3, "offsetEastMeters": 4}
	]}`))
	if err == nil {
		t.Fatal("expected error for duplicate uid, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "duplicate")
	}
}

func TestParse_RejectsMultipleMotionModels(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{
		"uid": "t1",
		"callsign": "C1",
		"orbit": {"radiusMeters": 100, "speedMps": 10},
		"raceTrack": {"legLengthMeters": 1000, "turnRadiusMeters": 100, "speedMps": 20}
	}]}`))
	if err == nil {
		t.Fatal("expected error when both orbit and raceTrack are set, got nil")
	}
	for _, want := range []string{"t1", "orbit", "raceTrack"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error = %q, want it to mention %q", err.Error(), want)
		}
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := Parse([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestParse_RejectsUnknownFields(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		wantField string
	}{
		{
			name:      "scenario",
			json:      `{"tickIntervalSecond": 2, "tracks": [{"uid": "t1", "callsign": "C1"}]}`,
			wantField: "tickIntervalSecond",
		},
		{
			name:      "track",
			json:      `{"tracks": [{"uid": "t1", "callsign": "C1", "speedMpss": 10}]}`,
			wantField: "speedMpss",
		},
		{
			name:      "orbit",
			json:      `{"tracks": [{"uid": "t1", "callsign": "C1", "orbit": {"radiusMeter": 100, "radiusMeters": 100, "speedMps": 10}}]}`,
			wantField: "radiusMeter",
		},
		{
			name:      "race track",
			json:      `{"tracks": [{"uid": "t1", "callsign": "C1", "raceTrack": {"legLengthMeter": 1000, "legLengthMeters": 1000, "turnRadiusMeters": 100, "speedMps": 20}}]}`,
			wantField: "legLengthMeter",
		},
		{
			name:      "sensor",
			json:      `{"tracks": [{"uid": "t1", "callsign": "C1", "sensor": {"fovDegrees": 30, "fovDeg": 30, "rangeMeters": 1000}}]}`,
			wantField: "fovDegrees",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.json))
			if err == nil || !strings.Contains(err.Error(), tt.wantField) {
				t.Fatalf("Parse returned %v, want unknown-field error mentioning %q", err, tt.wantField)
			}
		})
	}
}

func TestParse_RejectsTrailingJSON(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{"uid": "t1", "callsign": "C1"}]} {}`))
	if err == nil {
		t.Fatal("expected error for trailing JSON document, got nil")
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
