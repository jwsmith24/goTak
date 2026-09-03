package cot

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

func TestAirTrack_BuildEvent_NoSensorByDefault(t *testing.T) {
	track := AirTrack{
		UID:      "gotak-sim-1",
		Callsign: "SIM01",
		Time:     time.Now(),
	}

	xmlBytes, err := track.BuildEvent()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(string(xmlBytes), "<sensor") {
		t.Errorf("expected no <sensor> element when Sensor is nil, got:\n%s", xmlBytes)
	}
}

func TestAirTrack_BuildEvent_SensorAzimuthTracksCourse(t *testing.T) {
	track := AirTrack{
		UID:       "gotak-sim-1",
		Callsign:  "SIM01",
		CourseDeg: 90,
		Time:      time.Now(),
		Sensor: &SensorFOV{
			FOVDeg:      30,
			RangeMeters: 5000,
		},
	}

	xmlBytes, err := track.BuildEvent()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got eventXML
	if err := xml.Unmarshal(xmlBytes, &got); err != nil {
		t.Fatalf("BuildEvent produced invalid XML: %v\n%s", err, xmlBytes)
	}

	if got.Detail.Sensor == nil {
		t.Fatal("Detail.Sensor = nil, want a populated sensor element")
	}
	if got.Detail.Sensor.Azimuth != 90 {
		t.Errorf("Sensor.Azimuth = %v, want 90 (track's course, no offset)", got.Detail.Sensor.Azimuth)
	}
	if got.Detail.Sensor.FOV != 30 {
		t.Errorf("Sensor.FOV = %v, want 30", got.Detail.Sensor.FOV)
	}
	if got.Detail.Sensor.Range != 5000 {
		t.Errorf("Sensor.Range = %v, want 5000", got.Detail.Sensor.Range)
	}
}

func TestAirTrack_BuildEvent_SensorAzimuthAppliesOffsetAndWraps(t *testing.T) {
	track := AirTrack{
		UID:       "gotak-sim-1",
		Callsign:  "SIM01",
		CourseDeg: 350,
		Time:      time.Now(),
		Sensor: &SensorFOV{
			FOVDeg:           20,
			RangeMeters:      2000,
			AzimuthOffsetDeg: 30,
		},
	}

	xmlBytes, err := track.BuildEvent()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got eventXML
	if err := xml.Unmarshal(xmlBytes, &got); err != nil {
		t.Fatalf("BuildEvent produced invalid XML: %v", err)
	}

	// 350 + 30 = 380, wraps to 20.
	if got.Detail.Sensor.Azimuth != 20 {
		t.Errorf("Sensor.Azimuth = %v, want 20 (350+30 wrapped)", got.Detail.Sensor.Azimuth)
	}
}

func TestAirTrack_BuildEvent_SensorUpdatesAsCourseChanges(t *testing.T) {
	base := AirTrack{
		UID:      "gotak-sim-1",
		Callsign: "SIM01",
		Time:     time.Now(),
		Sensor:   &SensorFOV{FOVDeg: 30, RangeMeters: 5000},
	}

	first := base
	first.CourseDeg = 45
	firstXML, err := first.BuildEvent()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := base
	second.CourseDeg = 200
	secondXML, err := second.BuildEvent()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var gotFirst, gotSecond eventXML
	if err := xml.Unmarshal(firstXML, &gotFirst); err != nil {
		t.Fatalf("unmarshaling first event: %v", err)
	}
	if err := xml.Unmarshal(secondXML, &gotSecond); err != nil {
		t.Fatalf("unmarshaling second event: %v", err)
	}

	if gotFirst.Detail.Sensor.Azimuth != 45 {
		t.Errorf("first Sensor.Azimuth = %v, want 45", gotFirst.Detail.Sensor.Azimuth)
	}
	if gotSecond.Detail.Sensor.Azimuth != 200 {
		t.Errorf("second Sensor.Azimuth = %v, want 200", gotSecond.Detail.Sensor.Azimuth)
	}
}
