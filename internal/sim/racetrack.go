package sim

import (
	"math"
	"time"
)

// RaceTrackState describes a track flying a stadium-shaped ("race
// track") pattern: two straight legs of LegLengthMeters, joined by two
// 180-degree turns of TurnRadiusMeters, centered on (CenterLat,
// CenterLon) and oriented so the legs run along HeadingDeg.
type RaceTrackState struct {
	CenterLat, CenterLon float64 // degrees
	HeadingDeg           float64 // heading of the straight legs, degrees clockwise from north
	LegLengthMeters      float64
	TurnRadiusMeters     float64
	SpeedMPS             float64
	Clockwise            bool
	// DistanceMeters is the current distance traveled along the
	// pattern's perimeter, in [0, Perimeter()). 0 is the start of the
	// first straight leg.
	DistanceMeters float64
}

// Perimeter returns the length of one full circuit of the pattern: two
// straight legs plus one full circle's worth of turns (two semicircles).
func (r RaceTrackState) Perimeter() float64 {
	return 2*r.LegLengthMeters + 2*math.Pi*r.TurnRadiusMeters
}

// NewRaceTrackTrackState returns the initial state of a track flying
// pattern, at height hae, starting at pattern.DistanceMeters along the
// perimeter.
func NewRaceTrackTrackState(hae float64, pattern RaceTrackState) TrackState {
	lat, lon, course := pattern.positionAndCourse()
	return TrackState{
		Lat:       lat,
		Lon:       lon,
		HAE:       hae,
		CourseDeg: course,
		SpeedMPS:  pattern.SpeedMPS,
		RaceTrack: &pattern,
	}
}

// advanceRaceTrack moves s.RaceTrack.DistanceMeters forward by the
// distance covered in elapsed time, wrapping around the perimeter, then
// recomputes position and heading from the new distance.
func (s TrackState) advanceRaceTrack(elapsed time.Duration) TrackState {
	pattern := *s.RaceTrack

	perimeter := pattern.Perimeter()
	pattern.DistanceMeters = math.Mod(pattern.DistanceMeters+pattern.SpeedMPS*elapsed.Seconds(), perimeter)
	if pattern.DistanceMeters < 0 {
		pattern.DistanceMeters += perimeter
	}

	lat, lon, course := pattern.positionAndCourse()
	s.Lat = lat
	s.Lon = lon
	s.CourseDeg = course
	s.SpeedMPS = pattern.SpeedMPS
	s.RaceTrack = &pattern
	return s
}

// positionAndCourse computes the lat/lon and course for the pattern's
// current DistanceMeters, using a local tangent-plane approximation
// around the pattern center (accurate for the km-scale patterns this
// models). The pattern is built from two turn centers, C1 and C2, set
// LegLengthMeters apart along the heading, connected by straight legs
// offset TurnRadiusMeters to either side and semicircular turns at each
// end. Clockwise mirrors the whole shape across its long axis, which
// reverses the direction the turns bulge without changing the overall
// path shape.
func (r RaceTrackState) positionAndCourse() (lat, lon, courseDeg float64) {
	L := r.LegLengthMeters
	R := r.TurnRadiusMeters
	halfCircle := math.Pi * R
	perimeter := 2*L + 2*halfCircle

	d := math.Mod(r.DistanceMeters, perimeter)
	if d < 0 {
		d += perimeter
	}

	heading := degreesToRadians(r.HeadingDeg)
	along := [2]float64{math.Sin(heading), math.Cos(heading)} // unit vector (east, north) along the legs
	perp := [2]float64{math.Cos(heading), -math.Sin(heading)} // 90 degrees clockwise from along
	if !r.Clockwise {
		perp[0], perp[1] = -perp[0], -perp[1]
	}

	half := L / 2
	c1 := [2]float64{-half * along[0], -half * along[1]}
	c2 := [2]float64{half * along[0], half * along[1]}

	var posEast, posNorth, dirEast, dirNorth float64

	switch {
	case d < L:
		// Leg 1: from C1+R*perp to C2+R*perp, heading along.
		posEast = c1[0] + R*perp[0] + d*along[0]
		posNorth = c1[1] + R*perp[1] + d*along[1]
		dirEast, dirNorth = along[0], along[1]

	case d < L+halfCircle:
		// Turn around C2: from +R*perp to -R*perp, bulging outward (+along).
		theta := (d - L) / R
		cos, sin := math.Cos(theta), math.Sin(theta)
		posEast = c2[0] + R*(cos*perp[0]+sin*along[0])
		posNorth = c2[1] + R*(cos*perp[1]+sin*along[1])
		dirEast = -sin*perp[0] + cos*along[0]
		dirNorth = -sin*perp[1] + cos*along[1]

	case d < 2*L+halfCircle:
		// Leg 2: from C2-R*perp to C1-R*perp, heading -along.
		u := d - (L + halfCircle)
		posEast = c2[0] - R*perp[0] - u*along[0]
		posNorth = c2[1] - R*perp[1] - u*along[1]
		dirEast, dirNorth = -along[0], -along[1]

	default:
		// Turn around C1: from -R*perp back to +R*perp, bulging outward (-along).
		theta := (d - (2*L + halfCircle)) / R
		cos, sin := math.Cos(theta), math.Sin(theta)
		posEast = c1[0] - R*(cos*perp[0]+sin*along[0])
		posNorth = c1[1] - R*(cos*perp[1]+sin*along[1])
		dirEast = sin*perp[0] - cos*along[0]
		dirNorth = sin*perp[1] - cos*along[1]
	}

	lat, lon = offsetLatLon(r.CenterLat, r.CenterLon, posEast, posNorth)
	courseDeg = normalizeDegrees(radiansToDegrees(math.Atan2(dirEast, dirNorth)))
	return lat, lon, courseDeg
}
