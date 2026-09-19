package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jwsmith24/goTak/internal/enroll"
	"github.com/jwsmith24/goTak/internal/location"
	"github.com/jwsmith24/goTak/internal/menu"
	"github.com/jwsmith24/goTak/internal/sim"
)

type fakeSession struct{ closed bool }

func (*fakeSession) Send([]byte) error { return nil }
func (s *fakeSession) Close() error    { s.closed = true; return nil }

func testApp(t *testing.T) (app, *bytes.Buffer) {
	t.Helper()
	stdin, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = stdin.Close() })
	out := &bytes.Buffer{}
	a := newApp(stdin, out)
	a.loadLocations = func() ([]location.Location, error) {
		return append([]location.Location{}, location.All...), nil
	}
	a.startContext = func() (context.Context, context.CancelFunc) {
		return context.WithCancel(context.Background())
	}
	a.enroll = func(context.Context, *http.Client, string, string, string) (enroll.EnrollmentResult, error) {
		return enroll.EnrollmentResult{ClientCertPEM: []byte("cert"), PrivateKeyPEM: []byte("key"), CACertsPEM: [][]byte{[]byte("ca")}}, nil
	}
	a.dial = func(context.Context, string, []byte, []byte, [][]byte) (eventSenderCloser, error) {
		return &fakeSession{}, nil
	}
	a.simulate = func(context.Context, []*sim.Track, time.Duration, sim.EventSender) error { return nil }
	return a, out
}

func explicitArgs(extra ...string) []string {
	args := []string{"-server", "tak.example", "-username", "alice", "-password", "secret", "-location", "Austin, TX", "-scenario", ""}
	return append(args, extra...)
}

func TestRun_MenuCancellationStopsBeforeContext(t *testing.T) {
	a, _ := testApp(t)
	contextStarted := false
	a.startContext = func() (context.Context, context.CancelFunc) {
		contextStarted = true
		return context.WithCancel(context.Background())
	}
	a.runLocationMenu = func(*os.File, *bufio.Reader, io.Writer, []location.Location) (location.Location, error) {
		return location.Location{}, menu.ErrCancelled
	}

	err := a.run([]string{"-server", "tak.example", "-username", "alice", "-password", "secret", "-scenario", ""})
	if !errors.Is(err, menu.ErrCancelled) || contextStarted {
		t.Fatalf("run returned %v, contextStarted=%v", err, contextStarted)
	}
}

func TestRun_ScenarioLoadFailureStopsBeforeContext(t *testing.T) {
	a, _ := testApp(t)
	contextStarted := false
	a.startContext = func() (context.Context, context.CancelFunc) {
		contextStarted = true
		return context.WithCancel(context.Background())
	}
	missing := filepath.Join(t.TempDir(), "missing.json")
	args := explicitArgs()
	args[len(args)-1] = missing

	err := a.run(args)
	if err == nil || !strings.Contains(err.Error(), "loading scenario") || contextStarted {
		t.Fatalf("run returned %v, contextStarted=%v", err, contextStarted)
	}
}

func TestRun_EnrollmentFailure(t *testing.T) {
	a, out := testApp(t)
	want := errors.New("unavailable")
	a.enroll = func(context.Context, *http.Client, string, string, string) (enroll.EnrollmentResult, error) {
		return enroll.EnrollmentResult{}, want
	}

	err := a.run(explicitArgs())
	if !errors.Is(err, want) || !strings.Contains(err.Error(), "enrollment failed") {
		t.Fatalf("run returned %v", err)
	}
	if !strings.Contains(out.String(), "Enrolling with https://tak.example:8446 as alice") {
		t.Fatalf("output = %q", out.String())
	}
}

func TestRun_ConnectionFailure(t *testing.T) {
	a, out := testApp(t)
	want := errors.New("connection refused")
	a.dial = func(_ context.Context, addr string, cert, key []byte, cas [][]byte) (eventSenderCloser, error) {
		if addr != "tak.example:8089" || string(cert) != "cert" || string(key) != "key" || len(cas) != 1 {
			t.Fatalf("unexpected dial inputs: %q %q %q %d", addr, cert, key, len(cas))
		}
		return nil, want
	}

	err := a.run(explicitArgs())
	if !errors.Is(err, want) || !strings.Contains(err.Error(), "connecting to CoT stream") {
		t.Fatalf("run returned %v", err)
	}
	if !strings.Contains(out.String(), "Enrollment succeeded") || !strings.Contains(out.String(), "Connecting to CoT stream") {
		t.Fatalf("output = %q", out.String())
	}
}

func TestRun_ContextCancellationIsSuccessful(t *testing.T) {
	a, out := testApp(t)
	ctx, cancel := context.WithCancel(context.Background())
	a.startContext = func() (context.Context, context.CancelFunc) { return ctx, func() {} }
	session := &fakeSession{}
	a.dial = func(context.Context, string, []byte, []byte, [][]byte) (eventSenderCloser, error) {
		cancel()
		return session, nil
	}
	a.simulate = func(ctx context.Context, _ []*sim.Track, _ time.Duration, _ sim.EventSender) error {
		return ctx.Err()
	}

	if err := a.run(explicitArgs()); err != nil {
		t.Fatalf("run returned %v", err)
	}
	if !session.closed || !strings.Contains(out.String(), "Simulation stopped.") {
		t.Fatalf("closed=%v output=%q", session.closed, out.String())
	}
}

func TestRun_SharesReaderAcrossMenus(t *testing.T) {
	a, _ := testApp(t)
	var first *bufio.Reader
	a.runLocationMenu = func(_ *os.File, r *bufio.Reader, _ io.Writer, _ []location.Location) (location.Location, error) {
		first = r
		return location.All[0], nil
	}
	a.runScenarioMenu = func(_ *os.File, r *bufio.Reader, _ io.Writer, _ []menu.Option) (string, error) {
		if r != first {
			t.Fatal("location and scenario menus received different readers")
		}
		return "", nil
	}
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, scenariosDir), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, scenariosDir, "test.json"), []byte(`{"tracks":[{"uid":"t","callsign":"T"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)

	if err := a.run([]string{"-server", "tak.example", "-username", "alice", "-password", "secret"}); err != nil {
		t.Fatalf("run returned %v", err)
	}
}

func TestRun_ExplicitSavedLocationResolvesWithoutMenu(t *testing.T) {
	a, _ := testApp(t)
	saved := location.Location{Name: "Home", Lat: 31, Lon: -98}
	a.loadLocations = func() ([]location.Location, error) {
		return append(append([]location.Location{}, location.All...), saved), nil
	}
	a.runLocationMenu = func(*os.File, *bufio.Reader, io.Writer, []location.Location) (location.Location, error) {
		t.Fatal("location menu should be skipped for explicit saved location")
		return location.Location{}, nil
	}
	a.simulate = func(_ context.Context, tracks []*sim.Track, _ time.Duration, _ sim.EventSender) error {
		if tracks[0].State.Lat != saved.Lat || tracks[0].State.Lon != saved.Lon {
			t.Fatalf("default track origin = (%v, %v), want (%v, %v)", tracks[0].State.Lat, tracks[0].State.Lon, saved.Lat, saved.Lon)
		}
		return nil
	}
	args := explicitArgs()
	args[7] = saved.Name

	if err := a.run(args); err != nil {
		t.Fatalf("run returned %v", err)
	}
}
