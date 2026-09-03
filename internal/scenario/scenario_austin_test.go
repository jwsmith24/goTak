package scenario

import "testing"

// TestLoad_AustinCapitolScenario guards the checked-in scenario file
// against becoming invalid JSON or losing a required field.
func TestLoad_AustinCapitolScenario(t *testing.T) {
	sc, err := Load("../../scenarios/austin-capitol.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sc.Tracks) != 2 {
		t.Fatalf("len(Tracks) = %d, want 2", len(sc.Tracks))
	}

	for _, tr := range sc.Tracks {
		if tr.UID == "" || tr.Callsign == "" {
			t.Errorf("track %+v missing uid or callsign", tr)
		}
		if tr.Sensor == nil {
			t.Errorf("track %q: Sensor = nil, want a populated sensor", tr.UID)
		}
	}
}

// TestLoad_AustinCapitolHelicoptersScenario guards the checked-in
// composite scenario (helicopters, UAS, and ground units) against
// becoming invalid JSON or losing a required field.
func TestLoad_AustinCapitolHelicoptersScenario(t *testing.T) {
	sc, err := Load("../../scenarios/austin-capitol-helicopters.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const wantTracks = 7
	if len(sc.Tracks) != wantTracks {
		t.Fatalf("len(Tracks) = %d, want %d", len(sc.Tracks), wantTracks)
	}

	byUID := make(map[string]TrackConfig, len(sc.Tracks))
	for _, tr := range sc.Tracks {
		if tr.UID == "" || tr.Callsign == "" {
			t.Errorf("track %+v missing uid or callsign", tr)
		}
		byUID[tr.UID] = tr
	}

	helo1 := byUID["gotak-austin-helo01"]
	helo2 := byUID["gotak-austin-helo02"]
	if helo1.Orbit == nil || helo2.Orbit == nil {
		t.Fatal("expected both helicopters to have a populated orbit")
	}
	if helo1.Orbit.InitialBearingDeg == helo2.Orbit.InitialBearingDeg {
		t.Error("both helicopters start at the same bearing; expected them offset from each other")
	}

	hoveringUAS := byUID["gotak-austin-uas-hover"]
	if hoveringUAS.SpeedMPS != 0 || hoveringUAS.Orbit != nil || hoveringUAS.RaceTrack != nil {
		t.Errorf("hovering UAS = %+v, want a stationary (speed 0) straight-line track", hoveringUAS)
	}

	raceTrackUAS := byUID["gotak-austin-uas-racetrack"]
	if raceTrackUAS.RaceTrack == nil {
		t.Error("expected the second UAS to have a populated raceTrack pattern")
	}

	for _, uid := range []string{"gotak-austin-arty1", "gotak-austin-armor1", "gotak-austin-inf1"} {
		tr, ok := byUID[uid]
		if !ok {
			t.Errorf("missing expected ground unit track %q", uid)
			continue
		}
		if tr.Orbit != nil || tr.RaceTrack != nil {
			t.Errorf("ground unit %q should fly a plain straight-line (or stationary) track", uid)
		}
	}

	for _, uid := range []string{"gotak-austin-helo01", "gotak-austin-helo02", "gotak-austin-uas-hover", "gotak-austin-uas-racetrack"} {
		if byUID[uid].Sensor == nil {
			t.Errorf("track %q: Sensor = nil, want a populated sensor", uid)
		}
	}
}
