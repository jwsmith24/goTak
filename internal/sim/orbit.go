package sim

import "time"

// OrbitState describes a track looping at a fixed radius and speed
// around a center point, such as a helicopter orbiting a landmark.
type OrbitState struct {
	CenterLat, CenterLon float64 // degrees
	RadiusMeters         float64
	SpeedMPS             float64 // tangential ground speed
	Clockwise            bool
	// BearingFromCenterDeg is the current angle, clockwise from north,
	// from the center to the track's position.
	BearingFromCenterDeg float64
}

// NewOrbitTrackState returns the initial state of a track orbiting per
// orbit, at height hae, positioned at orbit.BearingFromCenterDeg.
func NewOrbitTrackState(hae float64, orbit OrbitState) TrackState {
	lat, lon := destinationPoint(orbit.CenterLat, orbit.CenterLon, orbit.BearingFromCenterDeg, orbit.RadiusMeters)
	return TrackState{
		Lat:       lat,
		Lon:       lon,
		HAE:       hae,
		CourseDeg: tangentialCourseDeg(orbit),
		SpeedMPS:  orbit.SpeedMPS,
		Orbit:     &orbit,
	}
}

// advanceOrbit moves s.Orbit.BearingFromCenterDeg forward by the angle
// covered in elapsed time at the orbit's tangential speed, then
// recomputes position and heading from the new bearing.
func (s TrackState) advanceOrbit(elapsed time.Duration) TrackState {
	orbit := *s.Orbit

	angularSpeedRadPerSec := orbit.SpeedMPS / orbit.RadiusMeters
	deltaDeg := radiansToDegrees(angularSpeedRadPerSec * elapsed.Seconds())
	if !orbit.Clockwise {
		deltaDeg = -deltaDeg
	}
	orbit.BearingFromCenterDeg = normalizeDegrees(orbit.BearingFromCenterDeg + deltaDeg)

	s.Lat, s.Lon = destinationPoint(orbit.CenterLat, orbit.CenterLon, orbit.BearingFromCenterDeg, orbit.RadiusMeters)
	s.CourseDeg = tangentialCourseDeg(orbit)
	s.SpeedMPS = orbit.SpeedMPS
	s.Orbit = &orbit
	return s
}

// tangentialCourseDeg returns the heading of travel at the orbit's
// current position: perpendicular to the radius, in the direction of
// rotation.
func tangentialCourseDeg(orbit OrbitState) float64 {
	if orbit.Clockwise {
		return normalizeDegrees(orbit.BearingFromCenterDeg + 90)
	}
	return normalizeDegrees(orbit.BearingFromCenterDeg - 90)
}
