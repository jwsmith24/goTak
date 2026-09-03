// Package cot builds Cursor-on-Target (CoT) XML events representing
// simulated air tracks.
package cot

import (
	"encoding/xml"
	"errors"
	"time"
)

const (
	// defaultType marks the track as a friendly air asset (affiliation
	// "f", dimension "A" for air) per the CoT type taxonomy.
	defaultType = "a-f-A"
	// defaultHow marks the position as machine-generated from a GPS-like
	// source, matching what a real air track feed would report.
	defaultHow = "m-g"
	// defaultStaleWindow is how long a position report stays valid
	// before the receiving system should consider it stale.
	defaultStaleWindow = 5 * time.Minute
	// defaultCircularError and defaultLinearError are placeholder
	// position/altitude error estimates for simulated tracks.
	defaultCircularError = 10.0
	defaultLinearError   = 10.0

	cotTimeLayout = "2006-01-02T15:04:05.000Z"
)

// AirTrack describes one position report for a simulated air track.
type AirTrack struct {
	UID         string
	Callsign    string
	Type        string // CoT type, e.g. "a-f-A"; defaults to a friendly air track
	Lat         float64
	Lon         float64
	HAE         float64 // height above the ellipsoid, in meters
	CourseDeg   float64
	SpeedMPS    float64
	Time        time.Time
	StaleWindow time.Duration // defaults to 5 minutes
}

type eventXML struct {
	XMLName xml.Name  `xml:"event"`
	Version string    `xml:"version,attr"`
	UID     string    `xml:"uid,attr"`
	Type    string    `xml:"type,attr"`
	Time    string    `xml:"time,attr"`
	Start   string    `xml:"start,attr"`
	Stale   string    `xml:"stale,attr"`
	How     string    `xml:"how,attr"`
	Point   pointXML  `xml:"point"`
	Detail  detailXML `xml:"detail"`
}

type pointXML struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	HAE float64 `xml:"hae,attr"`
	CE  float64 `xml:"ce,attr"`
	LE  float64 `xml:"le,attr"`
}

type detailXML struct {
	Contact contactXML `xml:"contact"`
	Track   trackXML   `xml:"track"`
}

type contactXML struct {
	Callsign string `xml:"callsign,attr"`
}

type trackXML struct {
	Course float64 `xml:"course,attr"`
	Speed  float64 `xml:"speed,attr"`
}

// BuildEvent renders the track as a CoT XML event, ready to be sent to a
// TAK server.
func (t AirTrack) BuildEvent() ([]byte, error) {
	if t.UID == "" {
		return nil, errors.New("cot: UID is required")
	}

	cotType := t.Type
	if cotType == "" {
		cotType = defaultType
	}

	staleWindow := t.StaleWindow
	if staleWindow == 0 {
		staleWindow = defaultStaleWindow
	}

	when := t.Time.UTC()
	timeStr := when.Format(cotTimeLayout)

	event := eventXML{
		Version: "2.0",
		UID:     t.UID,
		Type:    cotType,
		Time:    timeStr,
		Start:   timeStr,
		Stale:   when.Add(staleWindow).Format(cotTimeLayout),
		How:     defaultHow,
		Point: pointXML{
			Lat: t.Lat,
			Lon: t.Lon,
			HAE: t.HAE,
			CE:  defaultCircularError,
			LE:  defaultLinearError,
		},
		Detail: detailXML{
			Contact: contactXML{Callsign: t.Callsign},
			Track:   trackXML{Course: t.CourseDeg, Speed: t.SpeedMPS},
		},
	}

	return xml.Marshal(event)
}
