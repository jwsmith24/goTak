package sim

import (
	"math"
	"testing"
	"time"
)

func TestRaceTrackState_Perimeter(t *testing.T) {
	r := RaceTrackState{LegLengthMeters: 3000, TurnRadiusMeters: 500}
	want := 2*3000.0 + 2*math.Pi*500.0
	if !almostEqual(r.Perimeter(), want, 1e-6) {
		t.Errorf("Perimeter() = %v, want %v", r.Perimeter(), want)
	}
}

func TestNewRaceTrackTrackState_StartsHeadingAlongFirstLeg(t *testing.T) {
	state := NewRaceTrackTrackState(500, RaceTrackState{
		CenterLat: 30.2747, CenterLon: -97.7404,
		HeadingDeg: 45, LegLengthMeters: 3000, TurnRadiusMeters: 500,
		SpeedMPS: 30, Clockwise: true,
	})

	if !almostEqual(state.CourseDeg, 45, 1e-6) {
		t.Errorf("CourseDeg = %v, want 45", state.CourseDeg)
	}
	if state.HAE != 500 {
		t.Errorf("HAE = %v, want 500", state.HAE)
	}
	if state.SpeedMPS != 30 {
		t.Errorf("SpeedMPS = %v, want 30", state.SpeedMPS)
	}
}

func TestRaceTrackState_PathIsContinuousAndMatchesSpeed(t *testing.T) {
	for _, clockwise := range []bool{true, false} {
		pattern := RaceTrackState{
			CenterLat: 30.2747, CenterLon: -97.7404,
			HeadingDeg: 30, LegLengthMeters: 4000, TurnRadiusMeters: 800,
			SpeedMPS: 25, Clockwise: clockwise,
		}
		perimeter := pattern.Perimeter()

		const steps = 400
		step := perimeter / steps

		prev := NewRaceTrackTrackState(300, pattern)
		var totalDistance float64
		var maxStepDistance float64

		cur := pattern
		for i := 1; i <= steps; i++ {
			cur.DistanceMeters = step * float64(i)
			lat, lon, course := cur.positionAndCourse()
			d := greatCircleDistanceMeters(prev.Lat, prev.Lon, lat, lon)
			totalDistance += d
			if d > maxStepDistance {
				maxStepDistance = d
			}
			// Each step should move roughly `step` meters; a large jump
			// would indicate a discontinuity (a sign error at a segment
			// boundary).
			if d > step*1.5 {
				t.Fatalf("clockwise=%v: step %d jumped %v meters, want ~%v (course=%v)", clockwise, i, d, step, course)
			}
			prev.Lat, prev.Lon = lat, lon
		}

		if !almostEqual(totalDistance, perimeter, perimeter*0.02) {
			t.Errorf("clockwise=%v: total sampled distance = %v, want ~%v (perimeter)", clockwise, totalDistance, perimeter)
		}
	}
}

func TestRaceTrackState_FullCircuitReturnsToStart(t *testing.T) {
	pattern := RaceTrackState{
		CenterLat: 30.2747, CenterLon: -97.7404,
		HeadingDeg: 60, LegLengthMeters: 2000, TurnRadiusMeters: 400,
		SpeedMPS: 20, Clockwise: true,
	}
	start := NewRaceTrackTrackState(500, pattern)

	circuitSeconds := pattern.Perimeter() / pattern.SpeedMPS
	next := start.Advance(time.Duration(circuitSeconds * float64(time.Second)))

	if !almostEqual(next.Lat, start.Lat, 1e-4) || !almostEqual(next.Lon, start.Lon, 1e-4) {
		t.Errorf("after one full circuit, next = (%v, %v), want back at start (%v, %v)", next.Lat, next.Lon, start.Lat, start.Lon)
	}
	if !almostEqual(next.CourseDeg, start.CourseDeg, 1e-2) {
		t.Errorf("CourseDeg = %v, want back to start %v", next.CourseDeg, start.CourseDeg)
	}
}

func TestRaceTrackState_PreservesAltitudeAndSpeed(t *testing.T) {
	pattern := RaceTrackState{
		CenterLat: 30.2747, CenterLon: -97.7404,
		HeadingDeg: 10, LegLengthMeters: 1500, TurnRadiusMeters: 300,
		SpeedMPS: 18, Clockwise: false,
	}
	state := NewRaceTrackTrackState(700, pattern)

	next := state.Advance(37 * time.Second)

	if next.HAE != 700 {
		t.Errorf("HAE = %v, want unchanged 700", next.HAE)
	}
	if next.SpeedMPS != 18 {
		t.Errorf("SpeedMPS = %v, want unchanged 18", next.SpeedMPS)
	}
}
