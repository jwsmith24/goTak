// Command gotak simulates air tracks on a TAK server for development use.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jwsmith24/goTak/internal/config"
	"github.com/jwsmith24/goTak/internal/cot"
	"github.com/jwsmith24/goTak/internal/enroll"
	"github.com/jwsmith24/goTak/internal/scenario"
	"github.com/jwsmith24/goTak/internal/sim"
	"github.com/jwsmith24/goTak/internal/stream"
)

const (
	cotStreamPort        = "8089"
	defaultTickInterval  = 2 * time.Second
	defaultTrackUID      = "gotak-sim-1"
	defaultTrackCallsign = "SIM01"
)

// loadTracks returns the tracks to simulate and how often to update them.
// With no scenario file, it falls back to a single default track so the
// tool still runs with just server/username/password.
func loadTracks(scenarioPath string) ([]*sim.Track, time.Duration, error) {
	if scenarioPath == "" {
		return []*sim.Track{
			{
				UID:      defaultTrackUID,
				Callsign: defaultTrackCallsign,
				State: sim.TrackState{
					Lat: 34.05, Lon: -118.25, HAE: 3048,
					CourseDeg: 90, SpeedMPS: 128,
				},
			},
		}, defaultTickInterval, nil
	}

	sc, err := scenario.Load(scenarioPath)
	if err != nil {
		return nil, 0, fmt.Errorf("loading scenario %s: %w", scenarioPath, err)
	}

	tracks := make([]*sim.Track, len(sc.Tracks))
	for i, tc := range sc.Tracks {
		var state sim.TrackState
		if tc.Orbit != nil {
			state = sim.NewOrbitTrackState(tc.HAE, sim.OrbitState{
				CenterLat:            tc.Orbit.CenterLat,
				CenterLon:            tc.Orbit.CenterLon,
				RadiusMeters:         tc.Orbit.RadiusMeters,
				SpeedMPS:             tc.Orbit.SpeedMPS,
				Clockwise:            tc.Orbit.Clockwise,
				BearingFromCenterDeg: tc.Orbit.InitialBearingDeg,
			})
		} else {
			state = sim.TrackState{
				Lat: tc.Lat, Lon: tc.Lon, HAE: tc.HAE,
				CourseDeg: tc.CourseDeg, SpeedMPS: tc.SpeedMPS,
			}
		}

		var sensor *cot.SensorFOV
		if tc.Sensor != nil {
			sensor = &cot.SensorFOV{
				FOVDeg:           tc.Sensor.FOVDeg,
				RangeMeters:      tc.Sensor.RangeMeters,
				AzimuthOffsetDeg: tc.Sensor.AzimuthOffsetDeg,
			}
		}

		tracks[i] = &sim.Track{
			UID:      tc.UID,
			Callsign: tc.Callsign,
			Type:     tc.Type,
			State:    state,
			Sensor:   sensor,
		}
	}
	return tracks, sc.TickInterval(), nil
}

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "gotak:", err)
		os.Exit(1)
	}

	tracks, tickInterval, err := loadTracks(cfg.ScenarioPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gotak:", err)
		os.Exit(1)
	}

	baseURL := enroll.DefaultBaseURL(cfg.ServerAddress)
	fmt.Printf("Enrolling with %s as %s...\n", baseURL, cfg.Username)

	result, err := enroll.Enroll(context.Background(), enroll.InsecureHTTPClient(), baseURL, cfg.Username, cfg.Password)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gotak: enrollment failed:", err)
		os.Exit(1)
	}
	fmt.Printf("Enrollment succeeded: received client certificate and %d CA certificate(s).\n", len(result.CACertsPEM))

	streamAddr := cfg.ServerAddress + ":" + cotStreamPort
	fmt.Printf("Connecting to CoT stream at %s...\n", streamAddr)
	sender, err := stream.Dial(context.Background(), streamAddr, result.ClientCertPEM, result.PrivateKeyPEM, result.CACertsPEM)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gotak: connecting to CoT stream:", err)
		os.Exit(1)
	}
	defer sender.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	fmt.Printf("Simulating %d air track(s), updating every %s. Press Ctrl+C to stop.\n", len(tracks), tickInterval)

	if err := sim.Run(ctx, tracks, ticker.C, tickInterval, sender); err != nil && ctx.Err() == nil {
		fmt.Fprintln(os.Stderr, "gotak: simulation stopped:", err)
		os.Exit(1)
	}

	fmt.Println("Simulation stopped.")
}
