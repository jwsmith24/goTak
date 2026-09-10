// Package location provides the named central locations a scenario's
// tracks are positioned relative to.
package location

// Location is a named reference point on the earth's surface.
type Location struct {
	Name string
	Lat  float64
	Lon  float64
}

// All is the fixed list of locations a user can run a scenario from.
var All = []Location{
	{Name: "Austin, TX", Lat: 30.2747, Lon: -97.7404},                 // Texas State Capitol
	{Name: "Fort Campbell, KY", Lat: 36.6722, Lon: -87.4925},          // Campbell Army Airfield (KHOP)
	{Name: "Wheeler Army Airfield, HI", Lat: 21.4816, Lon: -158.0371}, // Wheeler AAF (PHHI), Oahu
}
