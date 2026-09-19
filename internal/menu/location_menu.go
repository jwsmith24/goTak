package menu

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/jwsmith24/goTak/internal/location"
)

// RunLocationMenu shows named locations plus a custom-coordinate option.
func RunLocationMenu(stdin *os.File, input *bufio.Reader, stdout io.Writer, locations []location.Location) (location.Location, error) {
	path, err := location.CustomLocationsPath()
	if err != nil {
		return location.Location{}, fmt.Errorf("menu: %w", err)
	}
	return runLocationMenu(stdin, input, stdout, locations, path)
}

func runLocationMenu(stdin *os.File, input *bufio.Reader, stdout io.Writer, locations []location.Location, customLocationsPath string) (location.Location, error) {
	return runLocationMenuWithChooser(stdin, input, stdout, locations, customLocationsPath, choose)
}

func runLocationMenuWithChooser(stdin *os.File, input *bufio.Reader, stdout io.Writer, locations []location.Location, customLocationsPath string, selectEntry func(*os.File, *bufio.Reader, io.Writer, string, []string) (int, error)) (location.Location, error) {
	entries := make([]string, len(locations)+1)
	for i, loc := range locations {
		entries[i] = loc.Name
	}
	entries[len(locations)] = customLocationLabel

	idx, err := selectEntry(stdin, input, stdout, "Select a location to run from:", entries)
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
