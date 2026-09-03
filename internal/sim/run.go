package sim

import (
	"context"
	"time"

	"github.com/jwsmith24/goTak/internal/cot"
)

// EventSender delivers a single CoT XML event, as implemented by
// stream.Sender.
type EventSender interface {
	Send(event []byte) error
}

// Track is one simulated air track's identity and current kinematic
// state.
type Track struct {
	UID      string
	Callsign string
	Type     string // CoT type; empty defaults to a friendly air track
	State    TrackState
}

// Run advances every track's position by interval on each tick received
// from ticks, builds a CoT event per track, and sends them via sender in
// order. It runs until ctx is cancelled or a send fails, returning the
// resulting error.
func Run(ctx context.Context, tracks []*Track, ticks <-chan time.Time, interval time.Duration, sender EventSender) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case now := <-ticks:
			for _, tr := range tracks {
				tr.State = tr.State.Advance(interval)

				event, err := cot.AirTrack{
					UID:       tr.UID,
					Callsign:  tr.Callsign,
					Type:      tr.Type,
					Lat:       tr.State.Lat,
					Lon:       tr.State.Lon,
					HAE:       tr.State.HAE,
					CourseDeg: tr.State.CourseDeg,
					SpeedMPS:  tr.State.SpeedMPS,
					Time:      now,
				}.BuildEvent()
				if err != nil {
					return err
				}
				if err := sender.Send(event); err != nil {
					return err
				}
			}
		}
	}
}
