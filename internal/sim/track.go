// Package sim advances a simulated air track's position over time and
// streams it to a TAK server as CoT events.
package sim

import (
	"math"
	"time"
)

// earthRadiusMeters is the mean radius used for the spherical-earth
// destination-point calculation in Advance.
const earthRadiusMeters = 6371000.0

// TrackState is the current kinematic state of a simulated air track.
// A track flies a straight course by default; setting Orbit or
// RaceTrack (mutually exclusive) instead loops it around a fixed point
// or around a stadium-shaped pattern. Either way, CourseDeg and SpeedMPS
// report the track's current heading and speed rather than necessarily
// driving its motion directly.
type TrackState struct {
	Lat, Lon  float64 // degrees
	HAE       float64 // height above ellipsoid, meters
	CourseDeg float64 // true course, degrees clockwise from north
	SpeedMPS  float64 // ground speed, meters/second
	Orbit     *OrbitState
	RaceTrack *RaceTrackState
}

// Advance returns the track's state after elapsed time has passed:
// flying its current course and speed in a straight line, continuing
// around its orbit, or continuing around its race-track pattern,
// depending on which of those is set. Altitude is left unchanged
// either way.
func (s TrackState) Advance(elapsed time.Duration) TrackState {
	switch {
	case s.Orbit != nil:
		return s.advanceOrbit(elapsed)
	case s.RaceTrack != nil:
		return s.advanceRaceTrack(elapsed)
	default:
		distance := s.SpeedMPS * elapsed.Seconds()
		s.Lat, s.Lon = destinationPoint(s.Lat, s.Lon, s.CourseDeg, distance)
		return s
	}
}

// destinationPoint returns the point reached by traveling distanceMeters
// from (lat, lon) along bearingDeg (degrees clockwise from north), using
// a spherical-earth destination-point formula.
func destinationPoint(lat, lon, bearingDeg, distanceMeters float64) (float64, float64) {
	if distanceMeters == 0 {
		return lat, lon
	}

	angularDistance := distanceMeters / earthRadiusMeters
	bearing := degreesToRadians(bearingDeg)
	lat1 := degreesToRadians(lat)
	lon1 := degreesToRadians(lon)

	lat2 := math.Asin(math.Sin(lat1)*math.Cos(angularDistance) +
		math.Cos(lat1)*math.Sin(angularDistance)*math.Cos(bearing))
	lon2 := lon1 + math.Atan2(
		math.Sin(bearing)*math.Sin(angularDistance)*math.Cos(lat1),
		math.Cos(angularDistance)-math.Sin(lat1)*math.Sin(lat2),
	)

	return radiansToDegrees(lat2), radiansToDegrees(lon2)
}

// offsetLatLon returns the point eastMeters/northMeters from (lat0, lon0),
// using a local equirectangular (flat-earth) approximation. This is
// accurate enough for the km-scale local geometry of a race-track
// pattern, and much simpler than exact spherical geodesics for that
// case.
func offsetLatLon(lat0, lon0, eastMeters, northMeters float64) (float64, float64) {
	lat := lat0 + radiansToDegrees(northMeters/earthRadiusMeters)
	lon := lon0 + radiansToDegrees(eastMeters/(earthRadiusMeters*math.Cos(degreesToRadians(lat0))))
	return lat, lon
}

func degreesToRadians(d float64) float64 { return d * math.Pi / 180 }
func radiansToDegrees(r float64) float64 { return r * 180 / math.Pi }

func normalizeDegrees(d float64) float64 {
	d = math.Mod(d, 360)
	if d < 0 {
		d += 360
	}
	return d
}
