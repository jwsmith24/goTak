package menu

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jwsmith24/goTak/internal/location"
)

// customLocationLabel is the entry appended to a location menu that
// lets the user enter their own lat/lon instead of picking one of the
// named locations.
const customLocationLabel = "Custom location (enter lat/lon)"

const defaultCustomLocationName = "Custom location"

// promptCustomLocation reads a name and lat/lon from r, prompting on w,
// and returns the resulting Location. An empty name defaults to
// defaultCustomLocationName.
func promptCustomLocation(r *bufio.Reader, w io.Writer) (location.Location, error) {
	name, err := promptLine(r, w, "Name for this location (optional): ")
	if err != nil {
		return location.Location{}, err
	}
	if name == "" {
		name = defaultCustomLocationName
	}

	lat, err := promptFloat(r, w, "Latitude (decimal degrees): ")
	if err != nil {
		return location.Location{}, err
	}

	lon, err := promptFloat(r, w, "Longitude (decimal degrees): ")
	if err != nil {
		return location.Location{}, err
	}

	return location.New(name, lat, lon)
}

// promptLine writes prompt to w, reads one line from r, and returns it
// with surrounding whitespace trimmed.
func promptLine(r *bufio.Reader, w io.Writer, prompt string) (string, error) {
	fmt.Fprint(w, prompt)

	line, err := r.ReadString('\n')
	if line == "" && err != nil {
		return "", fmt.Errorf("menu: reading input: %w", err)
	}
	return strings.TrimSpace(line), nil
}

// promptFloat is promptLine followed by float parsing.
func promptFloat(r *bufio.Reader, w io.Writer, prompt string) (float64, error) {
	line, err := promptLine(r, w, prompt)
	if err != nil {
		return 0, err
	}

	value, convErr := strconv.ParseFloat(line, 64)
	if convErr != nil {
		return 0, fmt.Errorf("menu: invalid number %q", line)
	}
	return value, nil
}
