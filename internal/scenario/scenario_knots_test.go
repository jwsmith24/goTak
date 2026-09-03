package scenario

import (
	"math"
	"strings"
	"testing"
)

func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestParse_TrackSpeedKtsConvertedToMPS(t *testing.T) {
	sc, err := Parse([]byte(`{"tracks": [{
		"uid": "t1", "callsign": "C1", "lat": 1, "lon": 2,
		"speedKts": 100
	}]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := 100 * 1852.0 / 3600.0
	if !almostEqual(sc.Tracks[0].SpeedMPS, want, 1e-9) {
		t.Errorf("SpeedMPS = %v, want %v (100 kts converted)", sc.Tracks[0].SpeedMPS, want)
	}
}

func TestParse_TrackRejectsBothSpeedMpsAndSpeedKts(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{
		"uid": "t1", "callsign": "C1", "lat": 1, "lon": 2,
		"speedMps": 50, "speedKts": 100
	}]}`))
	if err == nil {
		t.Fatal("expected error when both speedMps and speedKts are set, got nil")
	}
	if !strings.Contains(err.Error(), "speedMps") || !strings.Contains(err.Error(), "speedKts") {
		t.Errorf("error = %q, want it to mention both %q and %q", err.Error(), "speedMps", "speedKts")
	}
}

func TestParse_TrackWithoutSpeedKtsLeavesSpeedMpsUnchanged(t *testing.T) {
	sc, err := Parse([]byte(twoTrackJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Tracks[0].SpeedMPS != 120 {
		t.Errorf("SpeedMPS = %v, want unchanged 120", sc.Tracks[0].SpeedMPS)
	}
}

func TestParse_OrbitSpeedKtsConvertedToMPS(t *testing.T) {
	sc, err := Parse([]byte(`{"tracks": [{
		"uid": "helo-1", "callsign": "HELO01",
		"orbit": {"centerLat": 1, "centerLon": 2, "radiusMeters": 800, "speedKts": 70}
	}]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := 70 * 1852.0 / 3600.0
	if !almostEqual(sc.Tracks[0].Orbit.SpeedMPS, want, 1e-9) {
		t.Errorf("Orbit.SpeedMPS = %v, want %v (70 kts converted)", sc.Tracks[0].Orbit.SpeedMPS, want)
	}
}

func TestParse_OrbitRejectsBothSpeedMpsAndSpeedKts(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{
		"uid": "helo-1", "callsign": "HELO01",
		"orbit": {"centerLat": 1, "centerLon": 2, "radiusMeters": 800, "speedMps": 30, "speedKts": 70}
	}]}`))
	if err == nil {
		t.Fatal("expected error when both orbit speedMps and speedKts are set, got nil")
	}
	if !strings.Contains(err.Error(), "speedMps") || !strings.Contains(err.Error(), "speedKts") {
		t.Errorf("error = %q, want it to mention both %q and %q", err.Error(), "speedMps", "speedKts")
	}
}
