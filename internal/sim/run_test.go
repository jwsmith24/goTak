package sim

import (
	"context"
	"encoding/xml"
	"errors"
	"testing"
	"time"

	"github.com/jwsmith24/goTak/internal/cot"
)

type fakeSender struct {
	sent [][]byte
	err  error
}

func (f *fakeSender) Send(event []byte) error {
	if f.err != nil {
		return f.err
	}
	f.sent = append(f.sent, append([]byte(nil), event...))
	return nil
}

func TestRun_AdvancesAndSendsOnEachTick(t *testing.T) {
	ticks := make(chan time.Time, 3)
	sender := &fakeSender{}
	ctx, cancel := context.WithCancel(context.Background())

	tracks := []*Track{
		{UID: "uid-1", Callsign: "SIM01", State: TrackState{Lat: 0, Lon: 0, CourseDeg: 90, SpeedMPS: 100}},
	}
	interval := time.Second

	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, tracks, ticks, interval, sender)
	}()

	ticks <- time.Now()
	ticks <- time.Now()

	deadline := time.Now().Add(time.Second)
	for len(sender.sent) < 2 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	cancel()

	err := <-done
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run returned %v, want context.Canceled", err)
	}

	if len(sender.sent) != 2 {
		t.Fatalf("sender received %d events, want 2", len(sender.sent))
	}

	var first, second struct {
		UID   string `xml:"uid,attr"`
		Point struct {
			Lon float64 `xml:"lon,attr"`
		} `xml:"point"`
	}
	if err := xml.Unmarshal(sender.sent[0], &first); err != nil {
		t.Fatalf("unmarshaling first event: %v", err)
	}
	if err := xml.Unmarshal(sender.sent[1], &second); err != nil {
		t.Fatalf("unmarshaling second event: %v", err)
	}

	if first.UID != "uid-1" {
		t.Errorf("first.UID = %q, want %q", first.UID, "uid-1")
	}
	if second.Point.Lon <= first.Point.Lon {
		t.Errorf("second.Point.Lon (%v) should be greater than first.Point.Lon (%v): track should keep moving east", second.Point.Lon, first.Point.Lon)
	}
}

func TestRun_SendsAllTracksOnEachTick(t *testing.T) {
	ticks := make(chan time.Time, 1)
	sender := &fakeSender{}
	ctx, cancel := context.WithCancel(context.Background())

	tracks := []*Track{
		{UID: "track-1", Callsign: "EAGLE01", State: TrackState{Lat: 30.27, Lon: -97.76, CourseDeg: 90, SpeedMPS: 120}},
		{UID: "track-2", Callsign: "EAGLE02", State: TrackState{Lat: 30.26, Lon: -97.74, CourseDeg: 0, SpeedMPS: 100}},
	}

	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, tracks, ticks, time.Second, sender)
	}()

	ticks <- time.Now()

	deadline := time.Now().Add(time.Second)
	for len(sender.sent) < 2 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	cancel()
	<-done

	if len(sender.sent) != 2 {
		t.Fatalf("sender received %d events, want 2 (one per track)", len(sender.sent))
	}

	var uids []string
	for _, raw := range sender.sent {
		var got struct {
			UID string `xml:"uid,attr"`
		}
		if err := xml.Unmarshal(raw, &got); err != nil {
			t.Fatalf("unmarshaling event: %v", err)
		}
		uids = append(uids, got.UID)
	}

	if uids[0] != "track-1" || uids[1] != "track-2" {
		t.Errorf("uids = %v, want [track-1 track-2] in track order", uids)
	}
}

func TestRun_IncludesSensorFOVAlignedWithCourse(t *testing.T) {
	ticks := make(chan time.Time, 1)
	sender := &fakeSender{}
	ctx, cancel := context.WithCancel(context.Background())

	tracks := []*Track{
		{
			UID:      "uid-1",
			Callsign: "SIM01",
			State:    TrackState{Lat: 0, Lon: 0, CourseDeg: 45, SpeedMPS: 0},
			Sensor:   &cot.SensorFOV{FOVDeg: 30, RangeMeters: 5000},
		},
	}

	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, tracks, ticks, time.Second, sender)
	}()

	ticks <- time.Now()

	deadline := time.Now().Add(time.Second)
	for len(sender.sent) < 1 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	cancel()
	<-done

	if len(sender.sent) != 1 {
		t.Fatalf("sender received %d events, want 1", len(sender.sent))
	}

	var got struct {
		Detail struct {
			Sensor struct {
				Azimuth float64 `xml:"azimuth,attr"`
				FOV     float64 `xml:"fov,attr"`
				Range   float64 `xml:"range,attr"`
			} `xml:"sensor"`
		} `xml:"detail"`
	}
	if err := xml.Unmarshal(sender.sent[0], &got); err != nil {
		t.Fatalf("unmarshaling event: %v", err)
	}

	if got.Detail.Sensor.Azimuth != 45 {
		t.Errorf("Sensor.Azimuth = %v, want 45 (track's course)", got.Detail.Sensor.Azimuth)
	}
	if got.Detail.Sensor.FOV != 30 || got.Detail.Sensor.Range != 5000 {
		t.Errorf("Sensor FOV/Range = (%v, %v), want (30, 5000)", got.Detail.Sensor.FOV, got.Detail.Sensor.Range)
	}
}

func TestRun_StopsOnSenderError(t *testing.T) {
	ticks := make(chan time.Time, 1)
	sendErr := errors.New("connection lost")
	sender := &fakeSender{err: sendErr}

	tracks := []*Track{
		{UID: "uid-1", Callsign: "SIM01", State: TrackState{Lat: 0, Lon: 0, CourseDeg: 90, SpeedMPS: 100}},
	}

	done := make(chan error, 1)
	go func() {
		done <- Run(context.Background(), tracks, ticks, time.Second, sender)
	}()

	ticks <- time.Now()

	select {
	case err := <-done:
		if !errors.Is(err, sendErr) {
			t.Errorf("Run returned %v, want %v", err, sendErr)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after a send error")
	}
}
