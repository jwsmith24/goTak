package sim

import (
	"math"
	"testing"
)

func TestOffsetLatLon_ZeroOffsetReturnsOrigin(t *testing.T) {
	lat, lon := OffsetLatLon(30.2747, -97.7404, 0, 0)
	if lat != 30.2747 || lon != -97.7404 {
		t.Errorf("OffsetLatLon(origin, 0, 0) = (%v, %v), want the origin unchanged", lat, lon)
	}
}

func TestOffsetLatLon_NorthOffsetIncreasesLatitude(t *testing.T) {
	lat, lon := OffsetLatLon(0, 0, 0, earthRadiusMeters*degreesToRadians(1))
	if !almostEqual(lat, 1.0, 1e-6) {
		t.Errorf("lat = %v, want ~1.0", lat)
	}
	if !almostEqual(lon, 0.0, 1e-9) {
		t.Errorf("lon = %v, want ~0.0", lon)
	}
}

func TestOffsetLatLon_EastOffsetIncreasesLongitudeAlongEquator(t *testing.T) {
	lat, lon := OffsetLatLon(0, 0, earthRadiusMeters*degreesToRadians(1), 0)
	if !almostEqual(lat, 0.0, 1e-9) {
		t.Errorf("lat = %v, want ~0.0", lat)
	}
	if !almostEqual(lon, 1.0, 1e-6) {
		t.Errorf("lon = %v, want ~1.0", lon)
	}
}

func TestOffsetLatLon_RoundTripsWithGreatCircleDistance(t *testing.T) {
	lat, lon := OffsetLatLon(36.6722, -87.4925, 1500, -2000)
	d := greatCircleDistanceMeters(36.6722, -87.4925, lat, lon)
	want := math.Hypot(1500, -2000)
	if !almostEqual(d, want, 1.0) {
		t.Errorf("distance from origin = %v, want ~%v", d, want)
	}
}
