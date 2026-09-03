package scenario

import (
	"strings"
	"testing"
)

const sensorTrackJSON = `{
	"tracks": [
		{
			"uid": "track-1",
			"callsign": "EAGLE01",
			"lat": 30.2747,
			"lon": -97.76,
			"courseDeg": 90,
			"speedMps": 120,
			"sensor": {
				"fovDeg": 30,
				"rangeMeters": 8000,
				"azimuthOffsetDeg": 15
			}
		}
	]
}`

func TestParse_SensorTrack(t *testing.T) {
	sc, err := Parse([]byte(sensorTrackJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tr := sc.Tracks[0]
	if tr.Sensor == nil {
		t.Fatal("Sensor = nil, want a populated sensor config")
	}
	if tr.Sensor.FOVDeg != 30 {
		t.Errorf("FOVDeg = %v, want 30", tr.Sensor.FOVDeg)
	}
	if tr.Sensor.RangeMeters != 8000 {
		t.Errorf("RangeMeters = %v, want 8000", tr.Sensor.RangeMeters)
	}
	if tr.Sensor.AzimuthOffsetDeg != 15 {
		t.Errorf("AzimuthOffsetDeg = %v, want 15", tr.Sensor.AzimuthOffsetDeg)
	}
}

func TestParse_TrackWithoutSensorHasNilSensor(t *testing.T) {
	sc, err := Parse([]byte(twoTrackJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Tracks[0].Sensor != nil {
		t.Errorf("Sensor = %+v, want nil when not configured", sc.Tracks[0].Sensor)
	}
}

func TestParse_SensorRequiresPositiveFOV(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{
		"uid": "t1", "callsign": "C1", "lat": 1, "lon": 2,
		"sensor": {"fovDeg": 0, "rangeMeters": 5000}
	}]}`))
	if err == nil {
		t.Fatal("expected error for zero fovDeg, got nil")
	}
	if !strings.Contains(err.Error(), "fovDeg") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "fovDeg")
	}
}

func TestParse_SensorRequiresPositiveRange(t *testing.T) {
	_, err := Parse([]byte(`{"tracks": [{
		"uid": "t1", "callsign": "C1", "lat": 1, "lon": 2,
		"sensor": {"fovDeg": 30, "rangeMeters": 0}
	}]}`))
	if err == nil {
		t.Fatal("expected error for zero rangeMeters, got nil")
	}
	if !strings.Contains(err.Error(), "rangeMeters") {
		t.Errorf("error = %q, want it to mention %q", err.Error(), "rangeMeters")
	}
}
