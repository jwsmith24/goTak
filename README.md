# goTak

A terminal tool for simulating air tracks on a TAK server.

Given a server IP, username, and password, it:

1. Enrolls for a client certificate via the TAK server's
   username/password-authenticated enrollment endpoint.
2. Opens an mTLS connection to the server's CoT streaming port using that
   certificate.
3. Simulates one or more tracks, sending a CoT position update
   for each on a fixed interval until you stop it.

## Requirements

- Go 1.24 or later
- Network access to your TAK server on ports `8446` (certificate
  enrollment) and `8089` (CoT streaming)
- A username/password valid for certificate enrollment on that server

## Running

From the repository root:

```sh
go run ./cmd/gotak -server <server-ip> -username <username> -password <password>
```

For example:

```sh
go run ./cmd/gotak -server 192.168.1.50 -username dev -password devpass
```

Press `Ctrl+C` to stop the simulation.

### Building a binary

```sh
go build -o gotak ./cmd/gotak
./gotak -server 192.168.1.50 -username dev -password devpass
```

### Flags

| Flag         | Description                                                        |
|--------------|---------------------------------------------------------------------|
| `-server`    | TAK server IP address or hostname                                  |
| `-username`  | Username for certificate enrollment                                |
| `-password`  | Password for certificate enrollment                                |
| `-location`  | Name of the central location to run from (optional; see below)     |
| `-scenario`  | Path to a JSON scenario file (optional; see below)                  |

`-server`, `-username`, and `-password` are required; the app reports any
that are missing. `-location` and `-scenario` are optional.

### Choosing a location

Before the scenario menu, `gotak` shows a menu of central locations to
run the simulation from — every track position in a scenario file is an
offset from whichever one you pick, so the same scenario can be flown
near Austin one run and near Fort Campbell the next:

```
+------------------------------------------------------------------------------+
| Select a location to run from:                                               |
|                                                                                |
| > Austin, TX                                                                  |
|   Fort Campbell, KY                                                          |
|   Wheeler Army Airfield, HI                                                  |
|                                                                                |
| (Use up/down arrows and Enter; q to cancel)                                   |
+------------------------------------------------------------------------------+
```

It uses the same arrow-key/plain-numbered-list behavior described below
for the scenario menu. Only an explicit `-location` flag skips it (for
scripted/non-interactive runs) — a `GOTAK_LOCATION` default from `.env`
does not; the menu still shows so you can pick interactively each time.

### Choosing a scenario

Unless you pass `-scenario` explicitly on the command line, and the
`scenarios/` directory has any valid scenario files in it, `gotak` shows
an interactive menu before enrolling. Use the up/down arrow keys to move
the highlight and press Enter to pick:

```
+------------------------------------------------------------------------------+
| Select a configuration to run:                                               |
|                                                                                |
|   Default (single built-in track)                                            |
| > austin-capitol-helicopters.json - Helicopters, UAS, friendly ground uni...  |
|   austin-capitol.json - Two air tracks crossing paths near the Texas Capitol  |
|                                                                                |
| (Use up/down arrows and Enter; q to cancel)                                   |
+------------------------------------------------------------------------------+
```

The highlighted entry (shown here as `>`, rendered in reverse video in a
real terminal) is what Enter selects; the selection wraps around at
either end. Press `q`, Esc, or Ctrl+C to cancel without starting a
simulation.

If stdin isn't an interactive terminal (for example, input piped in from
a script or CI), `gotak` falls back to a plain numbered list read one
line at a time instead — no raw terminal mode, no arrow keys:

```
Select a configuration to run:
  1) Default (single built-in track)
  2) austin-capitol-helicopters.json - Helicopters, UAS, friendly ground units, and hostile infantry near the Texas Capitol
  3) austin-capitol.json - Two air tracks crossing paths near the Texas Capitol
Enter choice:
```

Only an explicit `-scenario` flag skips the menu entirely (for scripted/
non-interactive runs) — a `GOTAK_SCENARIO` default from `.env` does not;
the menu still shows so you can pick interactively each time. A scenario
file's optional top-level `"description"` field is what shows up next to
its filename in the menu.

### Scenario files

A scenario file describes one or more tracks to simulate and how often to
update them. Track positions are given as an offset in meters
(`offsetNorthMeters`/`offsetEastMeters`) from whichever central location
you pick in the location menu, not as absolute lat/lon — the same
scenario file runs near Austin, Fort Campbell, or Wheeler AAF unchanged:

```json
{
  "tickIntervalSeconds": 2,
  "tracks": [
    {
      "uid": "gotak-austin-eagle01",
      "callsign": "EAGLE01",
      "type": "a-f-A",
      "offsetNorthMeters": 0,
      "offsetEastMeters": -1882.19,
      "hae": 1500,
      "courseDeg": 90,
      "speedMps": 120
    },
    {
      "uid": "gotak-austin-eagle02",
      "callsign": "EAGLE02",
      "offsetNorthMeters": -1634.57,
      "offsetEastMeters": 0,
      "hae": 2000,
      "courseDeg": 0,
      "speedMps": 100
    }
  ]
}
```

- `description` (top-level, optional) is shown next to the file in the
  interactive scenario menu (see below).
- `uid` and `callsign` are required and must be unique per track.
- `type` is the CoT type (e.g. `a-f-A` for friendly air); defaults to
  `a-f-A` when omitted.
- `offsetNorthMeters`/`offsetEastMeters` place the track relative to the
  selected central location (default `0`, i.e. right at that location);
  `hae` is height above the ellipsoid in meters; `courseDeg` is true
  course in degrees clockwise from north.
- Speed is `speedMps` (meters/second) or `speedKts` (knots) — give one or
  the other, not both; `speedKts` is converted to meters/second when the
  file is loaded.
- `tickIntervalSeconds` controls how often every track's position updates;
  defaults to 2 seconds when omitted which aligns with real PLI of a moving track.

[`scenarios/austin-capitol.json`](scenarios/austin-capitol.json) ships
with the repo: two air tracks crossing paths near the Texas Capitol in
downtown Austin. Run it with:

```sh
go run ./cmd/gotak -server 192.168.1.50 -username dev -password devpass \
  -scenario scenarios/austin-capitol.json
```

#### Orbiting tracks

A track can loop around a fixed point instead of flying a straight
course by giving it an `orbit` object instead of
`offsetNorthMeters`/`offsetEastMeters`/`courseDeg`/`speedMps`:

```json
{
  "uid": "gotak-austin-helo01",
  "callsign": "HELO01",
  "type": "a-f-A-M-H",
  "hae": 300,
  "orbit": {
    "offsetNorthMeters": 0,
    "offsetEastMeters": 0,
    "radiusMeters": 800,
    "speedMps": 35,
    "clockwise": true,
    "initialBearingDeg": 0
  }
}
```

- `offsetNorthMeters`/`offsetEastMeters` place the orbit's center point
  relative to the selected central location (default `0`); `radiusMeters`
  is the orbit radius. Tangential ground speed is `speedMps` or
  `speedKts` (give one or the other); both radius and speed must be
  positive.
- `clockwise` sets rotation direction (default `false`, counterclockwise).
- `initialBearingDeg` places the track's starting position on the circle,
  as a compass bearing from the center (default `0`, due north of center).
- The track's reported course and speed update every tick to reflect its
  current heading around the circle.

#### Race-track (stadium) patterns

A track can fly a stadium-shaped "race track" pattern — two straight
legs joined by two 180-degree turns, the classic ISR loiter pattern
flown by fixed-wing UAS — by giving it a `raceTrack` object instead of
`offsetNorthMeters`/`offsetEastMeters`/`courseDeg`/`speedMps`:

```json
{
  "uid": "gotak-austin-uas-racetrack",
  "callsign": "RQ01",
  "type": "a-f-A-M-F-Q",
  "hae": 1800,
  "raceTrack": {
    "offsetNorthMeters": 1000.75,
    "offsetEastMeters": 1728.54,
    "headingDeg": 60,
    "legLengthMeters": 3000,
    "turnRadiusMeters": 600,
    "speedKts": 70,
    "clockwise": true
  }
}
```

- `offsetNorthMeters`/`offsetEastMeters` place the pattern's center point
  relative to the selected central location (default `0`); `headingDeg`
  is the compass heading of the two straight legs.
- `legLengthMeters` (length of each straight leg) and `turnRadiusMeters`
  (radius of each 180-degree turn) must both be positive.
- Ground speed is `speedMps` or `speedKts` (give one or the other); must
  be positive.
- `clockwise` mirrors the pattern (default `false`).
- The track's reported course updates every tick as it flies each leg
  and turn.

[`scenarios/austin-capitol-helicopters.json`](scenarios/austin-capitol-helicopters.json)
ships a fuller composite scenario near the Capitol:

- Two rotary-wing helicopters (`HELO01`, `HELO02`) orbiting the Capitol
  on the same circle, 180° offset from each other so they stay on
  opposite sides of the loop, at slightly different altitudes.
- Two UAS: `SCAN1`, a small rotary-wing/VTOL UAS hovering in place
  (`speedMps: 0`) just north of the Capitol, and `RQ01`, a fixed-wing UAS
  flying a race-track loiter pattern a couple of kilometers out.
- Three friendly ground units further out from the Capitol: an emplaced
  field artillery battery (`ARTY1`, stationary), an armor element
  (`ARMOR1`) moving cross-country, and an infantry element (`INF1`) on
  foot.

```sh
go run ./cmd/gotak -server 192.168.1.50 -username dev -password devpass \
  -scenario scenarios/austin-capitol-helicopters.json
```

#### Sensor field of view

A track can carry a steerable sensor's field of view by adding a
`sensor` object:

```json
{
  "sensor": {
    "fovDeg": 30,
    "rangeMeters": 8000,
    "azimuthOffsetDeg": 0
  }
}
```

This renders as a CoT `<sensor azimuth fov range/>` detail element
(per the standard CoT Sensor schema) that ATAK/WinTAK draw as a
field-of-view wedge on the map. `fovDeg` (horizontal field of view) and
`rangeMeters` must both be positive. The sensor's azimuth is recomputed
every tick from the track's *current* course plus `azimuthOffsetDeg`
(default `0`, straight ahead), so the FOV wedge stays pointed in the
track's direction of travel as it turns — including around an orbit,
where the heading is constantly changing, and around a race-track
pattern. Every air track in the shipped scenarios carries a
forward-looking sensor; the ground units in
`austin-capitol-helicopters.json` don't.

### Using a .env file instead of flags

Instead of passing flags every time, copy `.env.example` to `.env` in the
repository root and fill in your values:

```sh
cp .env.example .env
```

```
GOTAK_SERVER=192.168.1.50
GOTAK_USERNAME=dev
GOTAK_PASSWORD=devpass
GOTAK_LOCATION=Austin, TX
GOTAK_SCENARIO=scenarios/austin-capitol.json
```

Neither `GOTAK_LOCATION` nor `GOTAK_SCENARIO` skips its interactive menu
described above — they only matter as a fallback for when no menu is
shown at all (for example, if `scenarios/` isn't found relative to your
working directory, for the scenario menu). Use `-location`/`-scenario`
instead when you want to skip a menu and run non-interactively.

Then just run:

```sh
go run ./cmd/gotak
```

or 

```sh
./run.sh
```

Any flag you do pass on the command line overrides the matching value from
`.env`. `.env` is gitignored — never commit real credentials.

## Running the tests

```sh
go test ./...
```