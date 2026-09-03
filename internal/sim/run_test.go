package sim

import (
	"context"
	"encoding/xml"
	"errors"
	"testing"
	"time"
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

	initial := TrackState{Lat: 0, Lon: 0, CourseDeg: 90, SpeedMPS: 100}
	interval := time.Second

	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, "uid-1", "SIM01", initial, ticks, interval, sender)
	}()

	ticks <- time.Now()
	ticks <- time.Now()

	// Give the goroutine a moment to process both ticks before cancelling.
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

func TestRun_StopsOnSenderError(t *testing.T) {
	ticks := make(chan time.Time, 1)
	sendErr := errors.New("connection lost")
	sender := &fakeSender{err: sendErr}

	initial := TrackState{Lat: 0, Lon: 0, CourseDeg: 90, SpeedMPS: 100}

	done := make(chan error, 1)
	go func() {
		done <- Run(context.Background(), "uid-1", "SIM01", initial, ticks, time.Second, sender)
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
