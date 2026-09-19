// Package scenario loads a JSON description of one or more simulated air
// tracks.
package scenario

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

const defaultTickIntervalSeconds = 2.0

// knotsToMPS converts knots to meters/second using the international
// definition of the nautical mile (1852 meters).
func knotsToMPS(knots float64) float64 {
	return knots * 1852.0 / 3600.0
}

// TrackConfig describes one track's callsign, CoT type, and starting
// kinematic state. Positions are given as an offset in meters from
// whichever central location the scenario runs from (see the location
// package), not as absolute lat/lon, so the same scenario file can run
// from any of them. A track either flies a straight course
// (OffsetNorthMeters/OffsetEastMeters/CourseDeg/SpeedMPS, Orbit and
// RaceTrack nil) or loops around a fixed point (Orbit or RaceTrack set,
// in which case those straight-course fields are ignored — the starting
// position and heading are derived from the orbit/race-track instead).
type TrackConfig struct {
	UID               string           `json:"uid"`
	Callsign          string           `json:"callsign"`
	Type              string           `json:"type,omitempty"`
	OffsetNorthMeters float64          `json:"offsetNorthMeters,omitempty"`
	OffsetEastMeters  float64          `json:"offsetEastMeters,omitempty"`
	HAE               float64          `json:"hae,omitempty"`
	CourseDeg         float64          `json:"courseDeg,omitempty"`
	SpeedMPS          float64          `json:"speedMps,omitempty"`
	SpeedKts          float64          `json:"speedKts,omitempty"` // alternative to speedMps; converted into SpeedMPS during Parse
	Orbit             *OrbitConfig     `json:"orbit,omitempty"`
	RaceTrack         *RaceTrackConfig `json:"raceTrack,omitempty"`
	Sensor            *SensorConfig    `json:"sensor,omitempty"`
}

// OrbitConfig describes a track looping at a fixed radius and speed
// around a center point, given as an offset in meters from the
// scenario's central location.
type OrbitConfig struct {
	OffsetNorthMeters float64 `json:"offsetNorthMeters,omitempty"`
	OffsetEastMeters  float64 `json:"offsetEastMeters,omitempty"`
	RadiusMeters      float64 `json:"radiusMeters"`
	SpeedMPS          float64 `json:"speedMps,omitempty"`
	SpeedKts          float64 `json:"speedKts,omitempty"` // alternative to speedMps; converted into SpeedMPS during Parse
	Clockwise         bool    `json:"clockwise,omitempty"`
	InitialBearingDeg float64 `json:"initialBearingDeg,omitempty"`
}

// RaceTrackConfig describes a track flying a stadium-shaped ("race
// track") loiter pattern: two straight legs joined by two 180-degree
// turns, centered on an offset in meters from the scenario's central
// location.
type RaceTrackConfig struct {
	OffsetNorthMeters float64 `json:"offsetNorthMeters,omitempty"`
	OffsetEastMeters  float64 `json:"offsetEastMeters,omitempty"`
	HeadingDeg        float64 `json:"headingDeg,omitempty"`
	LegLengthMeters   float64 `json:"legLengthMeters"`
	TurnRadiusMeters  float64 `json:"turnRadiusMeters"`
	SpeedMPS          float64 `json:"speedMps,omitempty"`
	SpeedKts          float64 `json:"speedKts,omitempty"` // alternative to speedMps; converted into SpeedMPS during Parse
	Clockwise         bool    `json:"clockwise,omitempty"`
}

// SensorConfig describes a track's steerable sensor field of view, kept
// aligned with the track's current direction of travel plus an optional
// offset (e.g. a side-looking sensor).
type SensorConfig struct {
	FOVDeg           float64 `json:"fovDeg"`
	RangeMeters      float64 `json:"rangeMeters"`
	AzimuthOffsetDeg float64 `json:"azimuthOffsetDeg,omitempty"`
}

// Scenario describes a full simulation run: how often to send position
// updates, and the tracks to simulate.
type Scenario struct {
	Description         string        `json:"description,omitempty"`
	TickIntervalSeconds float64       `json:"tickIntervalSeconds,omitempty"`
	Tracks              []TrackConfig `json:"tracks"`
}

type jsonSchema struct {
	fields  map[string]*jsonSchema
	element *jsonSchema
}

var (
	scalarJSON   = &jsonSchema{}
	sensorSchema = &jsonSchema{fields: map[string]*jsonSchema{
		"fovDeg": scalarJSON, "rangeMeters": scalarJSON, "azimuthOffsetDeg": scalarJSON,
	}}
	orbitSchema = &jsonSchema{fields: map[string]*jsonSchema{
		"offsetNorthMeters": scalarJSON, "offsetEastMeters": scalarJSON,
		"radiusMeters": scalarJSON, "speedMps": scalarJSON, "speedKts": scalarJSON,
		"clockwise": scalarJSON, "initialBearingDeg": scalarJSON,
	}}
	raceTrackSchema = &jsonSchema{fields: map[string]*jsonSchema{
		"offsetNorthMeters": scalarJSON, "offsetEastMeters": scalarJSON,
		"headingDeg": scalarJSON, "legLengthMeters": scalarJSON, "turnRadiusMeters": scalarJSON,
		"speedMps": scalarJSON, "speedKts": scalarJSON, "clockwise": scalarJSON,
	}}
	trackSchema = &jsonSchema{fields: map[string]*jsonSchema{
		"uid": scalarJSON, "callsign": scalarJSON, "type": scalarJSON,
		"offsetNorthMeters": scalarJSON, "offsetEastMeters": scalarJSON, "hae": scalarJSON,
		"courseDeg": scalarJSON, "speedMps": scalarJSON, "speedKts": scalarJSON,
		"orbit": orbitSchema, "raceTrack": raceTrackSchema, "sensor": sensorSchema,
	}}
	scenarioSchema = &jsonSchema{fields: map[string]*jsonSchema{
		"description": scalarJSON, "tickIntervalSeconds": scalarJSON,
		"tracks": {element: trackSchema},
	}}
)

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
	if err := validateJSONSchema(data); err != nil {
		return Scenario{}, fmt.Errorf("scenario: parsing JSON: %w", err)
	}

	var sc Scenario
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&sc); err != nil {
		return Scenario{}, fmt.Errorf("scenario: parsing JSON: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			err = fmt.Errorf("multiple JSON values")
		}
		return Scenario{}, fmt.Errorf("scenario: parsing JSON: %w", err)
	}

	if len(sc.Tracks) == 0 {
		return Scenario{}, fmt.Errorf("scenario: must define at least one track")
	}

	seen := make(map[string]bool, len(sc.Tracks))
	for i := range sc.Tracks {
		tr := &sc.Tracks[i]

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
		if tr.Orbit != nil && tr.RaceTrack != nil {
			return Scenario{}, fmt.Errorf("scenario: track %q: specify only one of orbit or raceTrack", tr.UID)
		}

		if tr.SpeedMPS != 0 && tr.SpeedKts != 0 {
			return Scenario{}, fmt.Errorf("scenario: track %q: specify only one of speedMps or speedKts", tr.UID)
		}
		if tr.SpeedKts != 0 {
			tr.SpeedMPS = knotsToMPS(tr.SpeedKts)
		}

		if tr.Orbit != nil {
			if tr.Orbit.RadiusMeters <= 0 {
				return Scenario{}, fmt.Errorf("scenario: track %q: orbit radiusMeters must be positive", tr.UID)
			}
			if tr.Orbit.SpeedMPS != 0 && tr.Orbit.SpeedKts != 0 {
				return Scenario{}, fmt.Errorf("scenario: track %q: specify only one of orbit speedMps or speedKts", tr.UID)
			}
			if tr.Orbit.SpeedKts != 0 {
				tr.Orbit.SpeedMPS = knotsToMPS(tr.Orbit.SpeedKts)
			}
			if tr.Orbit.SpeedMPS <= 0 {
				return Scenario{}, fmt.Errorf("scenario: track %q: orbit speedMps must be positive", tr.UID)
			}
		}

		if tr.RaceTrack != nil {
			if tr.RaceTrack.LegLengthMeters <= 0 {
				return Scenario{}, fmt.Errorf("scenario: track %q: raceTrack legLengthMeters must be positive", tr.UID)
			}
			if tr.RaceTrack.TurnRadiusMeters <= 0 {
				return Scenario{}, fmt.Errorf("scenario: track %q: raceTrack turnRadiusMeters must be positive", tr.UID)
			}
			if tr.RaceTrack.SpeedMPS != 0 && tr.RaceTrack.SpeedKts != 0 {
				return Scenario{}, fmt.Errorf("scenario: track %q: specify only one of raceTrack speedMps or speedKts", tr.UID)
			}
			if tr.RaceTrack.SpeedKts != 0 {
				tr.RaceTrack.SpeedMPS = knotsToMPS(tr.RaceTrack.SpeedKts)
			}
			if tr.RaceTrack.SpeedMPS <= 0 {
				return Scenario{}, fmt.Errorf("scenario: track %q: raceTrack speedMps must be positive", tr.UID)
			}
		}

		if tr.Sensor != nil {
			if tr.Sensor.FOVDeg <= 0 {
				return Scenario{}, fmt.Errorf("scenario: track %q: sensor fovDeg must be positive", tr.UID)
			}
			if tr.Sensor.RangeMeters <= 0 {
				return Scenario{}, fmt.Errorf("scenario: track %q: sensor rangeMeters must be positive", tr.UID)
			}
		}
	}

	if sc.TickIntervalSeconds <= 0 {
		sc.TickIntervalSeconds = defaultTickIntervalSeconds
	}

	return sc, nil
}

func validateJSONSchema(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := validateJSONValue(decoder, scenarioSchema); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		if err == nil {
			return fmt.Errorf("multiple JSON values")
		}
		return err
	}
	return nil
}

func validateJSONValue(decoder *json.Decoder, schema *jsonSchema) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}

	switch delim {
	case '{':
		if schema == nil || schema.fields == nil {
			return fmt.Errorf("unexpected JSON object")
		}
		seen := make(map[string]bool)
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return fmt.Errorf("object field name is not a string")
			}
			if seen[key] {
				return fmt.Errorf("duplicate field %q", key)
			}
			seen[key] = true
			child, exists := schema.fields[key]
			if !exists {
				return fmt.Errorf("unknown field %q", key)
			}
			if err := validateJSONValue(decoder, child); err != nil {
				return err
			}
		}
		_, err = decoder.Token()
		return err
	case '[':
		if schema == nil || schema.element == nil {
			return fmt.Errorf("unexpected JSON array")
		}
		for decoder.More() {
			if err := validateJSONValue(decoder, schema.element); err != nil {
				return err
			}
		}
		_, err = decoder.Token()
		return err
	default:
		return fmt.Errorf("unexpected JSON delimiter %q", delim)
	}
}
