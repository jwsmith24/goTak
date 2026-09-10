package menu

import (
	"bufio"
	"io"
	"os"

	"github.com/jwsmith24/goTak/internal/location"
)

// RunLocationMenu shows a menu of the given locations and returns the
// one selected. See choose for how the menu is presented and why input
// must be shared across every menu shown during a run.
func RunLocationMenu(stdin *os.File, input *bufio.Reader, stdout io.Writer, locations []location.Location) (location.Location, error) {
	entries := make([]string, len(locations))
	for i, loc := range locations {
		entries[i] = loc.Name
	}

	idx, err := choose(stdin, input, stdout, "Select a location to run from:", entries)
	if err != nil {
		return location.Location{}, err
	}
	return locations[idx], nil
}
