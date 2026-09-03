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
	}
}

// TestLoad_AustinCapitolHelicoptersScenario guards the checked-in orbit
// scenario against becoming invalid JSON or losing a required field.
func TestLoad_AustinCapitolHelicoptersScenario(t *testing.T) {
	sc, err := Load("../../scenarios/austin-capitol-helicopters.json")
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
		if tr.Orbit == nil {
			t.Errorf("track %q: Orbit = nil, want a populated orbit", tr.UID)
		}
	}

	if sc.Tracks[0].Orbit.InitialBearingDeg == sc.Tracks[1].Orbit.InitialBearingDeg {
		t.Error("both tracks start at the same bearing; expected them offset from each other")
	}
}
