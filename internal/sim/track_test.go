package sim

import (
	"math"
	"testing"
	"time"
)

func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestTrackState_Advance_DueNorth(t *testing.T) {
	state := TrackState{
		Lat: 0, Lon: 0, HAE: 1000,
		CourseDeg: 0, // north
		SpeedMPS:  earthRadiusMeters * math.Pi / 180 / 3600,
	}

	next := state.Advance(time.Hour)

	// Traveling due north for the distance covering exactly one degree
	// of latitude should land ~1 degree north of the equator.
	if !almostEqual(next.Lat, 1.0, 0.01) {
		t.Errorf("Lat = %v, want ~1.0", next.Lat)
	}
	if !almostEqual(next.Lon, 0.0, 0.01) {
		t.Errorf("Lon = %v, want ~0.0", next.Lon)
	}
	if next.HAE != 1000 {
		t.Errorf("HAE = %v, want unchanged 1000", next.HAE)
	}
	if next.CourseDeg != 0 {
		t.Errorf("CourseDeg = %v, want unchanged 0", next.CourseDeg)
	}
}

func TestTrackState_Advance_DueEastAlongEquator(t *testing.T) {
	state := TrackState{
		Lat: 0, Lon: 0,
		CourseDeg: 90, // east
		SpeedMPS:  earthRadiusMeters * math.Pi / 180 / 3600,
	}

	next := state.Advance(time.Hour)

	if !almostEqual(next.Lat, 0.0, 0.01) {
		t.Errorf("Lat = %v, want ~0.0", next.Lat)
	}
	if !almostEqual(next.Lon, 1.0, 0.01) {
		t.Errorf("Lon = %v, want ~1.0", next.Lon)
	}
}

func TestTrackState_Advance_Stationary(t *testing.T) {
	state := TrackState{Lat: 12.3, Lon: 45.6, SpeedMPS: 0, CourseDeg: 270}

	next := state.Advance(time.Minute)

	if next.Lat != 12.3 || next.Lon != 45.6 {
		t.Errorf("stationary track moved: got (%v, %v)", next.Lat, next.Lon)
	}
}
