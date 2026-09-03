package scenario

import (
	"strings"
	"testing"
)

const raceTrackJSON = `{
	"tracks": [
		{
			"uid": "uas-1",
			"callsign": "RQ01",
			"type": "a-f-A-M-F-Q",
			"hae": 1200,
			"raceTrack": {
				"centerLat": 30.28,
				"centerLon": -97.75,
				"headingDeg": 45,
				"legLengthMeters": 4000,
				"turnRadiusMeters": 800,
				"speedKts": 60,
				"clockwise": true
			}
		}
	]
}`

func TestParse_RaceTrackTrack(t *testing.T) {
	sc, err := Parse([]byte(raceTrackJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tr := sc.Tracks[0]
	if tr.RaceTrack == nil {
		t.Fatal("RaceTrack = nil, want a populated race-track config")
	}
	if tr.RaceTrack.CenterLat != 30.28 || tr.RaceTrack.CenterLon != -97.75 {
		t.Errorf("center = (%v, %v), unexpected", tr.RaceTrack.CenterLat, tr.RaceTrack.CenterLon)
	}
	if tr.RaceTrack.HeadingDeg != 45 {
		t.Errorf("HeadingDeg = %v, want 45", tr.RaceTrack.HeadingDeg)
	}
	if tr.RaceTrack.LegLengthMeters != 4000 || tr.RaceTrack.TurnRadiusMeters != 800 {
		t.Errorf("leg/radius = (%v, %v), unexpected", tr.RaceTrack.LegLengthMeters, tr.RaceTrack.TurnRadiusMeters)
	}
	if !tr.RaceTrack.Clockwise {
		t.Error("Clockwise = false, want true")
	}

	want := 60 * 1852.0 / 3600.0
	if !almostEqual(tr.RaceTrack.SpeedMPS, want, 1e-9) {
		t.Errorf("SpeedMPS = %v, want %v (60 kts converted)", tr.RaceTrack.SpeedMPS, want)
	}
}

func TestParse_TrackWithoutRaceTrackHasNilRaceTrack(t *testing.T) {
	sc, err := Parse([]byte(twoTrackJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Tracks[0].RaceTrack != nil {
		t.Errorf("RaceTrack = %+v, want nil when not configured", sc.Tracks[0].RaceTrack)
	}
}

func TestParse_RaceTrackRequiresPositiveLegLength(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{
		"uid": "uas-1", "callsign": "RQ01",
		"raceTrack": {"centerLat": 1, "centerLon": 2, "headingDeg": 0, "legLengthMeters": 0, "turnRadiusMeters": 800, "speedMps": 30}
	}]}`))
	if err == nil {
		t.Fatal("expected error for zero legLengthMeters, got nil")
	}
	if !strings.Contains(err.Error(), "legLengthMeters") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "legLengthMeters")
	}
}

func TestParse_RaceTrackRequiresPositiveTurnRadius(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{
		"uid": "uas-1", "callsign": "RQ01",
		"raceTrack": {"centerLat": 1, "centerLon": 2, "headingDeg": 0, "legLengthMeters": 4000, "turnRadiusMeters": 0, "speedMps": 30}
	}]}`))
	if err == nil {
		t.Fatal("expected error for zero turnRadiusMeters, got nil")
	}
	if !strings.Contains(err.Error(), "turnRadiusMeters") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "turnRadiusMeters")
	}
}

func TestParse_RaceTrackRequiresPositiveSpeed(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{
		"uid": "uas-1", "callsign": "RQ01",
		"raceTrack": {"centerLat": 1, "centerLon": 2, "headingDeg": 0, "legLengthMeters": 4000, "turnRadiusMeters": 800, "speedMps": 0}
	}]}`))
	if err == nil {
		t.Fatal("expected error for zero speedMps, got nil")
	}
	if !strings.Contains(err.Error(), "speedMps") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "speedMps")
	}
}

func TestParse_RaceTrackRejectsBothSpeedMpsAndSpeedKts(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{
		"uid": "uas-1", "callsign": "RQ01",
		"raceTrack": {"centerLat": 1, "centerLon": 2, "headingDeg": 0, "legLengthMeters": 4000, "turnRadiusMeters": 800, "speedMps": 30, "speedKts": 60}
	}]}`))
	if err == nil {
		t.Fatal("expected error when both raceTrack speedMps and speedKts are set, got nil")
	}
	if !strings.Contains(err.Error(), "speedMps") || !strings.Contains(err.Error(), "speedKts") {
		t.Errorf("error = %q, want it to mention both %q and %q", err.Error(), "speedMps", "speedKts")
	}
}
