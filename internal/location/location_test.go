package location

import "testing"

func TestAll_HasThreeNamedLocations(t *testing.T) {
	if len(All) != 3 {
		t.Fatalf("len(All) = %d, want 3", len(All))
	}

	want := map[string]struct{ lat, lon float64 }{
		"Austin, TX":                {30.2747, -97.7404},
		"Fort Campbell, KY":         {36.6722, -87.4925},
		"Wheeler Army Airfield, HI": {21.4816, -158.0371},
	}

	seen := make(map[string]bool, len(All))
	for _, loc := range All {
		wantCoords, ok := want[loc.Name]
		if !ok {
			t.Errorf("unexpected location name %q", loc.Name)
			continue
		}
		seen[loc.Name] = true

		if !almostEqual(loc.Lat, wantCoords.lat, 0.01) || !almostEqual(loc.Lon, wantCoords.lon, 0.01) {
			t.Errorf("%s: (%v, %v), want approximately (%v, %v)", loc.Name, loc.Lat, loc.Lon, wantCoords.lat, wantCoords.lon)
		}
	}
	for name := range want {
		if !seen[name] {
			t.Errorf("missing expected location %q", name)
		}
	}
}

func TestAll_NamesAreUnique(t *testing.T) {
	seen := make(map[string]bool, len(All))
	for _, loc := range All {
		if seen[loc.Name] {
			t.Errorf("duplicate location name %q", loc.Name)
		}
		seen[loc.Name] = true
	}
}

func almostEqual(a, b, tolerance float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= tolerance
}
