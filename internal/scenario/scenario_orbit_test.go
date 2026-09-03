package scenario

import (
	"strings"
	"testing"
)

const orbitTrackJSON = `{
	"tracks": [
		{
			"uid": "helo-1",
			"callsign": "HELO01",
			"type": "a-f-A-M-H",
			"hae": 300,
			"orbit": {
				"centerLat": 30.2747,
				"centerLon": -97.7404,
				"radiusMeters": 800,
				"speedMps": 30,
				"clockwise": true,
				"initialBearingDeg": 90
			}
		}
	]
}`

func TestParse_OrbitTrack(t *testing.T) {
	sc, err := Parse([]byte(orbitTrackJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tr := sc.Tracks[0]
	if tr.Orbit == nil {
		t.Fatal("Orbit = nil, want a populated orbit config")
	}
	if tr.Orbit.CenterLat != 30.2747 || tr.Orbit.CenterLon != -97.7404 {
		t.Errorf("orbit center = (%v, %v), unexpected", tr.Orbit.CenterLat, tr.Orbit.CenterLon)
	}
	if tr.Orbit.RadiusMeters != 800 || tr.Orbit.SpeedMPS != 30 {
		t.Errorf("orbit radius/speed = (%v, %v), unexpected", tr.Orbit.RadiusMeters, tr.Orbit.SpeedMPS)
	}
	if !tr.Orbit.Clockwise {
		t.Error("Clockwise = false, want true")
	}
	if tr.Orbit.InitialBearingDeg != 90 {
		t.Errorf("InitialBearingDeg = %v, want 90", tr.Orbit.InitialBearingDeg)
	}
}

func TestParse_TrackWithoutOrbitHasNilOrbit(t *testing.T) {
	sc, err := Parse([]byte(twoTrackJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Tracks[0].Orbit != nil {
		t.Errorf("Orbit = %+v, want nil for a straight-line track", sc.Tracks[0].Orbit)
	}
}

func TestParse_OrbitRequiresPositiveRadius(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{
		"uid": "helo-1", "callsign": "HELO01",
		"orbit": {"centerLat": 1, "centerLon": 2, "radiusMeters": 0, "speedMps": 30}
	}]}`))
	if err == nil {
		t.Fatal("expected error for zero radiusMeters, got nil")
	}
	if !strings.Contains(err.Error(), "radiusMeters") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "radiusMeters")
	}
}

func TestParse_OrbitRequiresPositiveSpeed(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{
		"uid": "helo-1", "callsign": "HELO01",
		"orbit": {"centerLat": 1, "centerLon": 2, "radiusMeters": 800, "speedMps": 0}
	}]}`))
	if err == nil {
		t.Fatal("expected error for zero speedMps, got nil")
	}
	if !strings.Contains(err.Error(), "speedMps") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "speedMps")
	}
}
