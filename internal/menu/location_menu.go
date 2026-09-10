package menu

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/jwsmith24/goTak/internal/location"
)

// RunLocationMenu shows a menu of the given locations, plus an option
// to enter a custom lat/lon, and returns the one selected. A newly
// entered custom location is saved to location.CustomLocationsPath()
// so it's offered again on future runs. See choose for how the menu is
// presented and why input must be shared across every menu shown
// during a run.
func RunLocationMenu(stdin *os.File, input *bufio.Reader, stdout io.Writer, locations []location.Location) (location.Location, error) {
	path, err := location.CustomLocationsPath()
	if err != nil {
		return location.Location{}, fmt.Errorf("menu: %w", err)
	}
	return runLocationMenu(stdin, input, stdout, locations, path)
}

func runLocationMenu(stdin *os.File, input *bufio.Reader, stdout io.Writer, locations []location.Location, customLocationsPath string) (location.Location, error) {
	entries := make([]string, len(locations)+1)
	for i, loc := range locations {
		entries[i] = loc.Name
	}
	entries[len(locations)] = customLocationLabel

	idx, err := choose(stdin, input, stdout, "Select a location to run from:", entries)
	if err != nil {
		return location.Location{}, err
	}
	if idx < len(locations) {
		return locations[idx], nil
	}

	loc, err := promptCustomLocation(input, stdout)
	if err != nil {
		return location.Location{}, err
	}
	if err := location.SaveCustom(customLocationsPath, loc); err != nil {
		return location.Location{}, fmt.Errorf("menu: saving custom location: %w", err)
	}
	return loc, nil
}
