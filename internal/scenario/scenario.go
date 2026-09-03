// Package scenario loads a JSON description of one or more simulated air
// tracks.
package scenario

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const defaultTickIntervalSeconds = 2.0

// TrackConfig describes one track's callsign, CoT type, and starting
// kinematic state. A track either flies a straight course (Lat/Lon/
// CourseDeg/SpeedMPS, Orbit nil) or loops around a fixed point (Orbit
// set, in which case Lat/Lon/CourseDeg/SpeedMPS are ignored — the
// starting position and heading are derived from the orbit instead).
type TrackConfig struct {
	UID       string       `json:"uid"`
	Callsign  string       `json:"callsign"`
	Type      string       `json:"type,omitempty"`
	Lat       float64      `json:"lat,omitempty"`
	Lon       float64      `json:"lon,omitempty"`
	HAE       float64      `json:"hae,omitempty"`
	CourseDeg float64      `json:"courseDeg,omitempty"`
	SpeedMPS  float64      `json:"speedMps,omitempty"`
	Orbit     *OrbitConfig `json:"orbit,omitempty"`
}

// OrbitConfig describes a track looping at a fixed radius and speed
// around a center point.
type OrbitConfig struct {
	CenterLat         float64 `json:"centerLat"`
	CenterLon         float64 `json:"centerLon"`
	RadiusMeters      float64 `json:"radiusMeters"`
	SpeedMPS          float64 `json:"speedMps"`
	Clockwise         bool    `json:"clockwise,omitempty"`
	InitialBearingDeg float64 `json:"initialBearingDeg,omitempty"`
}

// Scenario describes a full simulation run: how often to send position
// updates, and the tracks to simulate.
type Scenario struct {
	TickIntervalSeconds float64       `json:"tickIntervalSeconds,omitempty"`
	Tracks              []TrackConfig `json:"tracks"`
}

// TickInterval returns how often each track's position should be updated.
func (s Scenario) TickInterval() time.Duration {
	return time.Duration(s.TickIntervalSeconds * float64(time.Second))
}

// Load reads and parses a scenario from a JSON file at path.
func Load(path string) (Scenario, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Scenario{}, err
	}
	return Parse(data)
}

// Parse validates and parses scenario JSON, applying a default tick
// interval when one isn't specified.
func Parse(data []byte) (Scenario, error) {
	var sc Scenario
	if err := json.Unmarshal(data, &sc); err != nil {
		return Scenario{}, fmt.Errorf("scenario: parsing JSON: %w", err)
	}

	if len(sc.Tracks) == 0 {
		return Scenario{}, fmt.Errorf("scenario: must define at least one track")
	}

	seen := make(map[string]bool, len(sc.Tracks))
	for i, tr := range sc.Tracks {
		if tr.UID == "" {
			return Scenario{}, fmt.Errorf("scenario: track %d: uid is required", i)
		}
		if seen[tr.UID] {
			return Scenario{}, fmt.Errorf("scenario: duplicate track uid %q", tr.UID)
		}
		seen[tr.UID] = true

		if tr.Callsign == "" {
			return Scenario{}, fmt.Errorf("scenario: track %q: callsign is required", tr.UID)
		}

		if tr.Orbit != nil {
			if tr.Orbit.RadiusMeters <= 0 {
				return Scenario{}, fmt.Errorf("scenario: track %q: orbit radiusMeters must be positive", tr.UID)
			}
			if tr.Orbit.SpeedMPS <= 0 {
				return Scenario{}, fmt.Errorf("scenario: track %q: orbit speedMps must be positive", tr.UID)
			}
		}
	}

	if sc.TickIntervalSeconds <= 0 {
		sc.TickIntervalSeconds = defaultTickIntervalSeconds
	}

	return sc, nil
}
