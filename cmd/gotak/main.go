// Command gotak simulates air tracks on a TAK server for development use.
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jwsmith24/goTak/internal/config"
	"github.com/jwsmith24/goTak/internal/cot"
	"github.com/jwsmith24/goTak/internal/enroll"
	"github.com/jwsmith24/goTak/internal/location"
	"github.com/jwsmith24/goTak/internal/menu"
	"github.com/jwsmith24/goTak/internal/scenario"
	"github.com/jwsmith24/goTak/internal/sim"
	"github.com/jwsmith24/goTak/internal/stream"
)

const (
	cotStreamPort        = "8089"
	defaultTickInterval  = 2 * time.Second
	defaultTrackUID      = "gotak-sim-1"
	defaultTrackCallsign = "SIM01"
	scenariosDir         = "scenarios"
)

// loadTracks returns the tracks to simulate and how often to update them.
// Track positions in a scenario file are offsets in meters from origin,
// the selected central location; the built-in default track (no scenario
// file) starts at origin itself. With no scenario file, it falls back to
// a single default track so the tool still runs with just
// server/username/password.
func loadTracks(scenarioPath string, origin location.Location) ([]*sim.Track, time.Duration, error) {
	if scenarioPath == "" {
		return []*sim.Track{
			{
				UID:      defaultTrackUID,
				Callsign: defaultTrackCallsign,
				State: sim.TrackState{
					Lat: origin.Lat, Lon: origin.Lon, HAE: 3048,
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
		switch {
		case tc.Orbit != nil:
			centerLat, centerLon := sim.OffsetLatLon(origin.Lat, origin.Lon, tc.Orbit.OffsetEastMeters, tc.Orbit.OffsetNorthMeters)
			state = sim.NewOrbitTrackState(tc.HAE, sim.OrbitState{
				CenterLat:            centerLat,
				CenterLon:            centerLon,
				RadiusMeters:         tc.Orbit.RadiusMeters,
				SpeedMPS:             tc.Orbit.SpeedMPS,
				Clockwise:            tc.Orbit.Clockwise,
				BearingFromCenterDeg: tc.Orbit.InitialBearingDeg,
			})
		case tc.RaceTrack != nil:
			centerLat, centerLon := sim.OffsetLatLon(origin.Lat, origin.Lon, tc.RaceTrack.OffsetEastMeters, tc.RaceTrack.OffsetNorthMeters)
			state = sim.NewRaceTrackTrackState(tc.HAE, sim.RaceTrackState{
				CenterLat:        centerLat,
				CenterLon:        centerLon,
				HeadingDeg:       tc.RaceTrack.HeadingDeg,
				LegLengthMeters:  tc.RaceTrack.LegLengthMeters,
				TurnRadiusMeters: tc.RaceTrack.TurnRadiusMeters,
				SpeedMPS:         tc.RaceTrack.SpeedMPS,
				Clockwise:        tc.RaceTrack.Clockwise,
			})
		default:
			lat, lon := sim.OffsetLatLon(origin.Lat, origin.Lon, tc.OffsetEastMeters, tc.OffsetNorthMeters)
			state = sim.TrackState{
				Lat: lat, Lon: lon, HAE: tc.HAE,
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

	// Shared across every interactive menu shown this run: each wraps
	// stdin exactly once, so a sequence of prompts doesn't drop bytes
	// buffered-but-unread by an earlier one.
	stdinReader := bufio.NewReader(os.Stdin)

	customLocationsPath, err := location.CustomLocationsPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "gotak:", err)
		os.Exit(1)
	}
	customLocations, err := location.LoadCustom(customLocationsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gotak:", err)
		os.Exit(1)
	}
	availableLocations := append(append([]location.Location{}, location.All...), customLocations...)

	// An explicit -location flag always skips the menu, so scripted/
	// non-interactive runs are unaffected. A GOTAK_LOCATION default from
	// .env is just a convenience and should not suppress the menu.
	origin := availableLocations[0]
	if cfg.LocationFromFlag {
		found := false
		for _, loc := range availableLocations {
			if loc.Name == cfg.LocationName {
				origin = loc
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "gotak: unknown location %q\n", cfg.LocationName)
			os.Exit(1)
		}
	} else {
		chosen, err := menu.RunLocationMenu(os.Stdin, stdinReader, os.Stdout, availableLocations)
		if errors.Is(err, menu.ErrCancelled) {
			fmt.Println("Cancelled.")
			os.Exit(0)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "gotak:", err)
			os.Exit(1)
		}
		origin = chosen
	}

	// An explicit -scenario flag always skips the menu, so scripted/
	// non-interactive runs are unaffected. A GOTAK_SCENARIO default from
	// .env is just a convenience and should not suppress the menu.
	if !cfg.ScenarioFromFlag {
		if scenarios := menu.DiscoverScenarios(scenariosDir); len(scenarios) > 0 {
			chosen, err := menu.RunScenarioMenu(os.Stdin, stdinReader, os.Stdout, scenarios)
			if errors.Is(err, menu.ErrCancelled) {
				fmt.Println("Cancelled.")
				os.Exit(0)
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "gotak:", err)
				os.Exit(1)
			}
			cfg.ScenarioPath = chosen
		}
	}

	tracks, tickInterval, err := loadTracks(cfg.ScenarioPath, origin)
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

	fmt.Printf("Simulating %d track(s), updating every %s. Press Ctrl+C to stop.\n", len(tracks), tickInterval)

	if err := sim.Run(ctx, tracks, ticker.C, tickInterval, sender); err != nil && ctx.Err() == nil {
		fmt.Fprintln(os.Stderr, "gotak: simulation stopped:", err)
		os.Exit(1)
	}

	fmt.Println("Simulation stopped.")
}
