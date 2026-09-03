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
