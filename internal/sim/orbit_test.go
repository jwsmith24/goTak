package sim

import (
	"math"
	"testing"
	"time"
)

// orbitTestState builds a track orbiting at 1 radian/second so tests can
// reason in whole radians of rotation without picking awkward numbers.
func orbitTestState(clockwise bool, initialBearingDeg float64) TrackState {
	const radius = 1000.0
	return NewOrbitTrackState(500, OrbitState{
		CenterLat:            30.2747,
		CenterLon:            -97.7404,
		RadiusMeters:         radius,
		SpeedMPS:             radius, // SpeedMPS/RadiusMeters == 1 rad/s
		Clockwise:            clockwise,
		BearingFromCenterDeg: initialBearingDeg,
	})
}

func TestNewOrbitTrackState_StartsAtBearingFromCenter(t *testing.T) {
	start := orbitTestState(true, 0)

	// Starting due north of the center, at the orbit radius.
	wantLat, wantLon := destinationPoint(30.2747, -97.7404, 0, 1000)
	if !almostEqual(start.Lat, wantLat, 1e-6) || !almostEqual(start.Lon, wantLon, 1e-6) {
		t.Errorf("start = (%v, %v), want (%v, %v)", start.Lat, start.Lon, wantLat, wantLon)
	}
	if start.HAE != 500 {
		t.Errorf("HAE = %v, want 500", start.HAE)
	}
	// Clockwise and due north of center: tangential course is due east.
	if !almostEqual(start.CourseDeg, 90, 1e-6) {
		t.Errorf("CourseDeg = %v, want 90", start.CourseDeg)
	}
	if start.SpeedMPS != 1000 {
		t.Errorf("SpeedMPS = %v, want 1000", start.SpeedMPS)
	}
}

func TestTrackState_Advance_OrbitClockwiseQuarterTurn(t *testing.T) {
	start := orbitTestState(true, 0)

	// 1 rad/s for pi/2 seconds is a quarter turn clockwise: 0 -> 90.
	quarterTurnSeconds := math.Pi / 2
	next := start.Advance(time.Duration(quarterTurnSeconds * float64(time.Second)))

	wantLat, wantLon := destinationPoint(30.2747, -97.7404, 90, 1000)
	if !almostEqual(next.Lat, wantLat, 1e-4) || !almostEqual(next.Lon, wantLon, 1e-4) {
		t.Errorf("next = (%v, %v), want (%v, %v)", next.Lat, next.Lon, wantLat, wantLon)
	}
	// Now due east of center, moving clockwise: tangential course is due south.
	if !almostEqual(next.CourseDeg, 180, 1e-2) {
		t.Errorf("CourseDeg = %v, want 180", next.CourseDeg)
	}
}

func TestTrackState_Advance_OrbitCounterclockwiseQuarterTurn(t *testing.T) {
	start := orbitTestState(false, 0)

	quarterTurnSeconds := math.Pi / 2
	next := start.Advance(time.Duration(quarterTurnSeconds * float64(time.Second)))

	// Counterclockwise from due north: 0 -> -90 (i.e. 270, due west).
	wantLat, wantLon := destinationPoint(30.2747, -97.7404, 270, 1000)
	if !almostEqual(next.Lat, wantLat, 1e-4) || !almostEqual(next.Lon, wantLon, 1e-4) {
		t.Errorf("next = (%v, %v), want (%v, %v)", next.Lat, next.Lon, wantLat, wantLon)
	}
}

func TestTrackState_Advance_OrbitFullTurnReturnsToStart(t *testing.T) {
	start := orbitTestState(true, 45)

	fullTurnSeconds := 2 * math.Pi
	next := start.Advance(time.Duration(fullTurnSeconds * float64(time.Second)))

	if !almostEqual(next.Lat, start.Lat, 1e-4) || !almostEqual(next.Lon, start.Lon, 1e-4) {
		t.Errorf("after a full turn, next = (%v, %v), want back at start (%v, %v)", next.Lat, next.Lon, start.Lat, start.Lon)
	}
}

func TestTrackState_Advance_OrbitPreservesRadiusFromCenter(t *testing.T) {
	start := orbitTestState(true, 30)

	next := start.Advance(700 * time.Millisecond)

	gotDistance := greatCircleDistanceMeters(30.2747, -97.7404, next.Lat, next.Lon)
	if !almostEqual(gotDistance, 1000, 1.0) {
		t.Errorf("distance from center = %v, want ~1000", gotDistance)
	}
}

// greatCircleDistanceMeters is a test-only haversine helper independent of
// destinationPoint, so TestTrackState_Advance_OrbitPreservesRadiusFromCenter
// doesn't just re-derive the same formula it's checking.
func greatCircleDistanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	p1, p2 := degreesToRadians(lat1), degreesToRadians(lat2)
	dp := degreesToRadians(lat2 - lat1)
	dl := degreesToRadians(lon2 - lon1)
	a := math.Sin(dp/2)*math.Sin(dp/2) + math.Cos(p1)*math.Cos(p2)*math.Sin(dl/2)*math.Sin(dl/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusMeters * c
}
