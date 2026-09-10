package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// New returns a Location, validating that lat/lon are real-world
// coordinates.
func New(name string, lat, lon float64) (Location, error) {
	if lat < -90 || lat > 90 {
		return Location{}, fmt.Errorf("location: latitude %v out of range [-90, 90]", lat)
	}
	if lon < -180 || lon > 180 {
		return Location{}, fmt.Errorf("location: longitude %v out of range [-180, 180]", lon)
	}
	return Location{Name: name, Lat: lat, Lon: lon}, nil
}

// CustomLocationsPath returns the path to the local file where
// user-entered custom locations are persisted across runs.
func CustomLocationsPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("location: finding user config directory: %w", err)
	}
	return filepath.Join(dir, "gotak", "custom_locations.json"), nil
}

// LoadCustom reads previously saved custom locations from path. A
// missing file is not an error; it returns an empty slice.
func LoadCustom(path string) ([]Location, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("location: reading %s: %w", path, err)
	}

	var locs []Location
	if err := json.Unmarshal(data, &locs); err != nil {
		return nil, fmt.Errorf("location: parsing %s: %w", path, err)
	}
	return locs, nil
}

// SaveCustom appends loc to the custom locations persisted at path,
// creating the file (and its parent directory) if they don't exist yet.
func SaveCustom(path string, loc Location) error {
	existing, err := LoadCustom(path)
	if err != nil {
		return err
	}
	existing = append(existing, loc)

	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return fmt.Errorf("location: encoding %s: %w", path, err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("location: creating directory for %s: %w", path, err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("location: writing %s: %w", path, err)
	}
	return nil
}
