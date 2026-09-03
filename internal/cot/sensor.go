package cot

import "math"

// SensorFOV describes a steerable sensor's field of view, per the CoT
// sensor detail schema. Its azimuth is derived each time an event is
// built from the track's current course, so the FOV stays aligned with
// the direction of travel as the track maneuvers.
type SensorFOV struct {
	FOVDeg           float64 // horizontal field of view, degrees
	RangeMeters      float64 // sensor range, meters
	AzimuthOffsetDeg float64 // offset from the track's course; 0 = sensor points straight ahead
}

type sensorXML struct {
	Azimuth float64 `xml:"azimuth,attr"`
	FOV     float64 `xml:"fov,attr"`
	Range   float64 `xml:"range,attr"`
}

func normalizeDegrees(d float64) float64 {
	d = math.Mod(d, 360)
	if d < 0 {
		d += 360
	}
	return d
}
