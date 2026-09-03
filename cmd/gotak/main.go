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
	"github.com/jwsmith24/goTak/internal/enroll"
	"github.com/jwsmith24/goTak/internal/sim"
	"github.com/jwsmith24/goTak/internal/stream"
)

const (
	cotStreamPort = "8089"
	tickInterval  = 2 * time.Second
	trackUID      = "gotak-sim-1"
	trackCallsign = "SIM01"
)

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
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

	initial := sim.TrackState{
		Lat: 34.05, Lon: -118.25, HAE: 3048,
		CourseDeg: 90, SpeedMPS: 128,
	}

	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	fmt.Printf("Simulating air track %s (%s), updating every %s. Press Ctrl+C to stop.\n", trackUID, trackCallsign, tickInterval)

	if err := sim.Run(ctx, trackUID, trackCallsign, initial, ticker.C, tickInterval, sender); err != nil && ctx.Err() == nil {
		fmt.Fprintln(os.Stderr, "gotak: simulation stopped:", err)
		os.Exit(1)
	}

	fmt.Println("Simulation stopped.")
}
