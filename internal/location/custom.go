package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
)

var customLocationsMu sync.Mutex

func New(name string, lat, lon float64) (Location, error) {
	if math.IsNaN(lat) || math.IsInf(lat, 0) || lat < -90 || lat > 90 {
		return Location{}, fmt.Errorf("location: latitude %v out of range [-90, 90]", lat)
	}
	if math.IsNaN(lon) || math.IsInf(lon, 0) || lon < -180 || lon > 180 {
		return Location{}, fmt.Errorf("location: longitude %v out of range [-180, 180]", lon)
	}
	return Location{Name: name, Lat: lat, Lon: lon}, nil
}

func CustomLocationsPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("location: finding user config directory: %w", err)
	}
	return filepath.Join(dir, "gotak", "custom_locations.json"), nil
}

func LoadCustom(path string) ([]Location, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("location: reading %s: %w", path, err)
	}

	var locations []Location
	if err := json.Unmarshal(data, &locations); err != nil {
		return nil, fmt.Errorf("location: parsing %s: %w", path, err)
	}
	seen := make(map[string]bool, len(All)+len(locations))
	for _, loc := range All {
		seen[loc.Name] = true
	}
	for _, loc := range locations {
		if _, err := New(loc.Name, loc.Lat, loc.Lon); err != nil {
			return nil, fmt.Errorf("location: parsing %s: %w", path, err)
		}
		if seen[loc.Name] {
			return nil, fmt.Errorf("location: parsing %s: duplicate location name %q", path, loc.Name)
		}
		seen[loc.Name] = true
	}
	return locations, nil
}

func SaveCustom(path string, loc Location) error {
	customLocationsMu.Lock()
	defer customLocationsMu.Unlock()

	existing, err := LoadCustom(path)
	if err != nil {
		return err
	}
	for _, candidate := range append(append([]Location{}, All...), existing...) {
		if candidate.Name == loc.Name {
			return fmt.Errorf("location: name %q already exists", loc.Name)
		}
	}
	existing = append(existing, loc)

	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return fmt.Errorf("location: encoding %s: %w", path, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("location: creating directory for %s: %w", path, err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".custom_locations-*")
	if err != nil {
		return fmt.Errorf("location: creating temporary file for %s: %w", path, err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("location: securing temporary file for %s: %w", path, err)
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("location: writing temporary file for %s: %w", path, err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("location: syncing temporary file for %s: %w", path, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("location: closing temporary file for %s: %w", path, err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("location: replacing %s: %w", path, err)
	}
	return nil
}
