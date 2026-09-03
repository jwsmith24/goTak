package cot

import (
	"encoding/xml"
	"testing"
	"time"
)

func TestAirTrack_BuildEvent_DefaultsAndFields(t *testing.T) {
	when := time.Date(2026, 9, 3, 14, 0, 0, 0, time.UTC)

	track := AirTrack{
		UID:       "gotak-sim-1",
		Callsign:  "SIM01",
		Lat:       34.05,
		Lon:       -118.25,
		HAE:       3048.0,
		CourseDeg: 90,
		SpeedMPS:  120,
		Time:      when,
	}

	xmlBytes, err := track.BuildEvent()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got eventXML
	if err := xml.Unmarshal(xmlBytes, &got); err != nil {
		t.Fatalf("BuildEvent produced invalid XML: %v\n%s", err, xmlBytes)
	}

	if got.Version != "2.0" {
		t.Errorf("Version = %q, want %q", got.Version, "2.0")
	}
	if got.UID != "gotak-sim-1" {
		t.Errorf("UID = %q, want %q", got.UID, "gotak-sim-1")
	}
	if got.Type != "a-f-A" {
		t.Errorf("Type = %q, want default %q", got.Type, "a-f-A")
	}
	if got.How != "m-g" {
		t.Errorf("How = %q, want %q", got.How, "m-g")
	}

	wantTime := "2026-09-03T14:00:00.000Z"
	if got.Time != wantTime {
		t.Errorf("Time = %q, want %q", got.Time, wantTime)
	}
	if got.Start != wantTime {
		t.Errorf("Start = %q, want %q", got.Start, wantTime)
	}
	wantStale := "2026-09-03T14:05:00.000Z"
	if got.Stale != wantStale {
		t.Errorf("Stale = %q, want %q (default 5m stale window)", got.Stale, wantStale)
	}

	if got.Point.Lat != 34.05 {
		t.Errorf("Point.Lat = %v, want %v", got.Point.Lat, 34.05)
	}
	if got.Point.Lon != -118.25 {
		t.Errorf("Point.Lon = %v, want %v", got.Point.Lon, -118.25)
	}
	if got.Point.HAE != 3048.0 {
		t.Errorf("Point.HAE = %v, want %v", got.Point.HAE, 3048.0)
	}

	if got.Detail.Contact.Callsign != "SIM01" {
		t.Errorf("Detail.Contact.Callsign = %q, want %q", got.Detail.Contact.Callsign, "SIM01")
	}
	if got.Detail.Track.Course != 90 {
		t.Errorf("Detail.Track.Course = %v, want %v", got.Detail.Track.Course, 90)
	}
	if got.Detail.Track.Speed != 120 {
		t.Errorf("Detail.Track.Speed = %v, want %v", got.Detail.Track.Speed, 120)
	}
}

func TestAirTrack_BuildEvent_CustomTypeAndStaleWindow(t *testing.T) {
	track := AirTrack{
		UID:         "gotak-sim-2",
		Callsign:    "SIM02",
		Type:        "a-h-A",
		Time:        time.Date(2026, 9, 3, 14, 0, 0, 0, time.UTC),
		StaleWindow: time.Minute,
	}

	xmlBytes, err := track.BuildEvent()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got eventXML
	if err := xml.Unmarshal(xmlBytes, &got); err != nil {
		t.Fatalf("BuildEvent produced invalid XML: %v", err)
	}

	if got.Type != "a-h-A" {
		t.Errorf("Type = %q, want %q", got.Type, "a-h-A")
	}
	wantStale := "2026-09-03T14:01:00.000Z"
	if got.Stale != wantStale {
		t.Errorf("Stale = %q, want %q", got.Stale, wantStale)
	}
}

func TestAirTrack_BuildEvent_RequiresUID(t *testing.T) {
	track := AirTrack{Callsign: "SIM01", Time: time.Now()}

	_, err := track.BuildEvent()
	if err == nil {
		t.Fatal("expected error for missing UID, got nil")
	}
}
