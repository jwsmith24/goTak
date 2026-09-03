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
type TrackState struct {
	Lat, Lon  float64 // degrees
	HAE       float64 // height above ellipsoid, meters
	CourseDeg float64 // true course, degrees clockwise from north
	SpeedMPS  float64 // ground speed, meters/second
}

// Advance returns the track's state after flying its current course and
// speed for the given elapsed time, using a spherical-earth
// destination-point formula. Altitude and course are left unchanged.
func (s TrackState) Advance(elapsed time.Duration) TrackState {
	distance := s.SpeedMPS * elapsed.Seconds()
	if distance == 0 {
		return s
	}

	angularDistance := distance / earthRadiusMeters
	bearing := degreesToRadians(s.CourseDeg)
	lat1 := degreesToRadians(s.Lat)
	lon1 := degreesToRadians(s.Lon)

	lat2 := math.Asin(math.Sin(lat1)*math.Cos(angularDistance) +
		math.Cos(lat1)*math.Sin(angularDistance)*math.Cos(bearing))
	lon2 := lon1 + math.Atan2(
		math.Sin(bearing)*math.Sin(angularDistance)*math.Cos(lat1),
		math.Cos(angularDistance)-math.Sin(lat1)*math.Sin(lat2),
	)

	s.Lat = radiansToDegrees(lat2)
	s.Lon = radiansToDegrees(lon2)
	return s
}

func degreesToRadians(d float64) float64 { return d * math.Pi / 180 }
func radiansToDegrees(r float64) float64 { return r * 180 / math.Pi }
