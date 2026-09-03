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

// Run advances a track's position by interval on every tick received
// from ticks, builds a CoT event for the new position, and sends it via
// sender. It runs until ctx is cancelled or a send fails, returning the
// resulting error.
func Run(ctx context.Context, uid, callsign string, initial TrackState, ticks <-chan time.Time, interval time.Duration, sender EventSender) error {
	state := initial
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case now := <-ticks:
			state = state.Advance(interval)

			track := cot.AirTrack{
				UID:       uid,
				Callsign:  callsign,
				Lat:       state.Lat,
				Lon:       state.Lon,
				HAE:       state.HAE,
				CourseDeg: state.CourseDeg,
				SpeedMPS:  state.SpeedMPS,
				Time:      now,
			}

			event, err := track.BuildEvent()
			if err != nil {
				return err
			}
			if err := sender.Send(event); err != nil {
				return err
			}
		}
	}
}
