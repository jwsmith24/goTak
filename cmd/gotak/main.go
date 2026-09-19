// Command gotak simulates air tracks on a TAK server for development use.
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
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

type eventSenderCloser interface {
	sim.EventSender
	io.Closer
}

type app struct {
	stdin  *os.File
	stdout io.Writer

	startContext  func() (context.Context, context.CancelFunc)
	loadLocations func() ([]location.Location, error)

	runLocationMenu func(*os.File, *bufio.Reader, io.Writer, []location.Location) (location.Location, error)
	runScenarioMenu func(*os.File, *bufio.Reader, io.Writer, []menu.Option) (string, error)
	enroll          func(context.Context, *http.Client, string, string, string) (enroll.EnrollmentResult, error)
	dial            func(context.Context, string, []byte, []byte, [][]byte) (eventSenderCloser, error)
	simulate        func(context.Context, []*sim.Track, time.Duration, sim.EventSender) error
}

const (
	cotStreamPort        = "8089"
	defaultTickInterval  = 2 * time.Second
	defaultTrackUID      = "gotak-sim-1"
	defaultTrackCallsign = "SIM01"
	scenariosDir         = "scenarios"
	enrollmentTimeout    = 30 * time.Second
	streamConnectTimeout = 30 * time.Second
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

func newApp(stdin *os.File, stdout io.Writer) app {
	return app{
		stdin:  stdin,
		stdout: stdout,
		startContext: func() (context.Context, context.CancelFunc) {
			return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		},
		loadLocations: func() ([]location.Location, error) {
			path, err := location.CustomLocationsPath()
			if err != nil {
				return nil, err
			}
			custom, err := location.LoadCustom(path)
			if err != nil {
				return nil, err
			}
			return append(append([]location.Location{}, location.All...), custom...), nil
		},
		runLocationMenu: menu.RunLocationMenu,
		runScenarioMenu: menu.RunScenarioMenu,
		enroll:          enroll.Enroll,
		dial: func(ctx context.Context, addr string, cert, key []byte, cas [][]byte) (eventSenderCloser, error) {
			return stream.Dial(ctx, addr, cert, key, cas)
		},
		simulate: func(ctx context.Context, tracks []*sim.Track, interval time.Duration, sender sim.EventSender) error {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			return sim.Run(ctx, tracks, ticker.C, interval, sender)
		},
	}
}

func (a app) run(args []string) error {
	cfg, err := config.ParseFlags(args)
	if err != nil {
		return err
	}

	stdinReader := bufio.NewReader(a.stdin)
	availableLocations, err := a.loadLocations()
	if err != nil {
		return err
	}

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
			return fmt.Errorf("unknown location %q", cfg.LocationName)
		}
	} else {
		chosen, err := a.runLocationMenu(a.stdin, stdinReader, a.stdout, availableLocations)
		if err != nil {
			return err
		}
		origin = chosen
	}

	// An explicit -scenario flag always skips the menu, so scripted/
	// non-interactive runs are unaffected. A GOTAK_SCENARIO default from
	// .env is just a convenience and should not suppress the menu.
	if !cfg.ScenarioFromFlag {
		if scenarios := menu.DiscoverScenarios(scenariosDir); len(scenarios) > 0 {
			chosen, err := a.runScenarioMenu(a.stdin, stdinReader, a.stdout, scenarios)
			if err != nil {
				return err
			}
			cfg.ScenarioPath = chosen
		}
	}

	tracks, tickInterval, err := loadTracks(cfg.ScenarioPath, origin)
	if err != nil {
		return err
	}
	ctx, stop := a.startContext()
	defer stop()

	baseURL := enroll.DefaultBaseURL(cfg.ServerAddress)
	fmt.Fprintf(a.stdout, "Enrolling with %s as %s...\n", baseURL, cfg.Username)

	enrollCtx, cancelEnroll := context.WithTimeout(ctx, enrollmentTimeout)
	result, err := a.enroll(enrollCtx, enroll.InsecureHTTPClient(), baseURL, cfg.Username, cfg.Password)
	cancelEnroll()
	if err != nil {
		return fmt.Errorf("enrollment failed: %w", err)
	}
	fmt.Fprintf(a.stdout, "Enrollment succeeded: received client certificate and %d CA certificate(s).\n", len(result.CACertsPEM))

	streamAddr := cfg.ServerAddress + ":" + cotStreamPort
	fmt.Fprintf(a.stdout, "Connecting to CoT stream at %s...\n", streamAddr)
	connectCtx, cancelConnect := context.WithTimeout(ctx, streamConnectTimeout)
	sender, err := a.dial(connectCtx, streamAddr, result.ClientCertPEM, result.PrivateKeyPEM, result.CACertsPEM)
	cancelConnect()
	if err != nil {
		return fmt.Errorf("connecting to CoT stream: %w", err)
	}
	defer sender.Close()

	fmt.Fprintf(a.stdout, "Simulating %d track(s), updating every %s. Press Ctrl+C to stop.\n", len(tracks), tickInterval)

	if err := a.simulate(ctx, tracks, tickInterval, sender); err != nil && ctx.Err() == nil {
		return fmt.Errorf("simulation stopped: %w", err)
	}

	fmt.Fprintln(a.stdout, "Simulation stopped.")
	return nil
}

func main() {
	err := newApp(os.Stdin, os.Stdout).run(os.Args[1:])
	if errors.Is(err, menu.ErrCancelled) {
		fmt.Fprintln(os.Stdout, "Cancelled.")
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "gotak:", err)
		os.Exit(1)
	}
}
