# goTak User Guide

goTak enrolls for an in-memory client certificate, connects to a TAK server's Cursor-on-Target (CoT) stream over mTLS, and sends simulated track positions until stopped.

## Requirements

- Go 1.24.7 or later.
- Network access to the TAK server on port `8446` for enrollment and port `8089` for CoT streaming.
- A username and password accepted by the TAK server's certificate-enrollment endpoint.

Run commands from the repository root. goTak resolves `.env` and the `scenarios/` directory relative to the current working directory.

## Run Interactively

Pass the required connection values:

```sh
go run ./cmd/gotak \
  -server 192.168.1.50 \
  -username dev \
  -password devpass
```

goTak presents a location menu followed by a scenario menu. In an interactive terminal, use the up and down arrow keys and Enter; press `q`, Escape, or Ctrl+C to cancel. When input is not attached to a terminal, the menus require a numbered choice instead.

After selection, goTak enrolls, connects, and starts sending track updates. Press Ctrl+C to stop.

## Run Non-Interactively

Supply both optional selection flags to skip the menus:

```sh
go run ./cmd/gotak \
  -server 192.168.1.50 \
  -username dev \
  -password devpass \
  -location "Austin, TX" \
  -scenario scenarios/austin-capitol.json
```

Available flags:

| Flag | Required | Description |
|---|---:|---|
| `-server` | Yes | TAK server IP address or hostname. |
| `-username` | Yes | Username used for certificate enrollment. |
| `-password` | Yes | Password used for certificate enrollment. |
| `-location` | No | Exact name of the scenario origin; skips the location menu when explicitly supplied. |
| `-scenario` | No | Path to a JSON scenario; skips the scenario menu when explicitly supplied. |

If no scenario is selected or configured, goTak uses one built-in track at the selected origin.

## Use a `.env` File

Create a local file from the example:

```sh
cp .env.example .env
```

Set the values:

```dotenv
GOTAK_SERVER=192.168.1.50
GOTAK_USERNAME=dev
GOTAK_PASSWORD=devpass
GOTAK_SCENARIO=scenarios/austin-capitol.json
```

Then run:

```sh
go run ./cmd/gotak
```

or:

```sh
./run.sh
```

Explicit flags override matching `.env` values. `GOTAK_SCENARIO` is a fallback only: unlike an explicit `-scenario` flag, it does not suppress an available interactive menu. Location selection is configured with the interactive menu or the explicit `-location` flag; the CLI does not currently use `GOTAK_LOCATION`.

`.env` can contain credentials and is gitignored. Do not commit it.

## Build a Binary

```sh
go build -o gotak ./cmd/gotak
```

Run the binary with the same flags:

```sh
./gotak -server 192.168.1.50 -username dev -password devpass
```

## Choose a Location

Scenario coordinates are offsets from a named origin, so the same scenario can run in different areas. The built-in locations are:

- `Austin, TX`
- `Fort Campbell, KY`
- `Wheeler Army Airfield, HI`

The `-location` value must match one of these names exactly.

## Choose a Scenario

Without an explicit `-scenario`, goTak discovers valid `.json` files in `scenarios/` and presents them alongside the built-in default track. A scenario's optional top-level `description` is displayed in the menu.

The repository includes:

- `scenarios/austin-capitol.json`: two crossing air tracks.
- `scenarios/austin-capitol-helicopters.json`: helicopters, UAS, friendly ground units, and hostile infantry using multiple supported motion models.

## Write a Scenario

A scenario controls the common update interval and one or more tracks. Positions use east/north meter offsets from the selected location, not fixed latitude/longitude.

Scenario JSON is strict: every field must use the exact case-sensitive name listed below. Unknown, duplicated, or misspelled fields cause the entire scenario to be rejected instead of being silently ignored. A file must contain exactly one JSON object.

```json
{
  "description": "Two tracks crossing near the selected origin",
  "tickIntervalSeconds": 2,
  "tracks": [
    {
      "uid": "gotak-eagle01",
      "callsign": "EAGLE01",
      "type": "a-f-A",
      "offsetNorthMeters": 0,
      "offsetEastMeters": -1800,
      "hae": 1500,
      "courseDeg": 90,
      "speedMps": 120
    }
  ]
}
```

Top-level fields:

| Field | Required | Meaning |
|---|---:|---|
| `description` | No | Text displayed beside the filename in the scenario menu. |
| `tickIntervalSeconds` | No | Common update interval; defaults to 2 seconds when omitted or nonpositive. |
| `tracks` | Yes | Nonempty array of track objects. |

Track fields:

| Field | Required | Meaning |
|---|---:|---|
| `uid` | Yes | Identifier, unique within the scenario. |
| `callsign` | Yes | Display callsign. |
| `type` | No | CoT type; defaults to friendly air (`a-f-A`). |
| `offsetNorthMeters` | No | Straight-track northward offset from the selected origin; defaults to `0`. |
| `offsetEastMeters` | No | Straight-track eastward offset from the selected origin; defaults to `0`. |
| `hae` | No | Height above the ellipsoid in meters; defaults to `0`. |
| `courseDeg` | No | Straight-track true course clockwise from north; defaults to `0`. |
| `speedMps` | No | Straight-track ground speed in meters per second. |
| `speedKts` | No | Alternative straight-track ground speed in knots. |
| `orbit` | No | Circular-orbit configuration described below. |
| `raceTrack` | No | Race-track configuration described below. |
| `sensor` | No | Sensor field-of-view configuration described below. |

Speed can be specified as `speedMps` or `speedKts`. If both values are nonzero, the scenario is rejected. Knots are converted to meters per second when the scenario is loaded.

The `orbit` and `raceTrack` objects are mutually exclusive. When either is present, its configuration determines position and motion; top-level straight-motion offsets, course, and speed are ignored.

On each tick, goTak advances every track first and then emits its event. The first event is therefore sent after one tick interval at the first advanced position, not immediately at the configured starting position.

### Straight Motion

Use the track-level position, course, and speed fields:

```json
{
  "uid": "gotak-eagle01",
  "callsign": "EAGLE01",
  "offsetNorthMeters": 0,
  "offsetEastMeters": -1800,
  "hae": 1500,
  "courseDeg": 90,
  "speedMps": 120
}
```

A straight track may use speed `0` to remain stationary.

### Circular Orbit

Use an `orbit` object:

```json
{
  "uid": "gotak-helo01",
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

The radius and speed must be positive. The offsets locate the orbit center. `initialBearingDeg` is a compass bearing from that center and defaults to north. `clockwise` defaults to `false`.

Orbit fields:

| Field | Required | Meaning |
|---|---:|---|
| `offsetNorthMeters` | No | Northward orbit-center offset; defaults to `0`. |
| `offsetEastMeters` | No | Eastward orbit-center offset; defaults to `0`. |
| `radiusMeters` | Yes | Positive orbit radius. |
| `speedMps` | Conditional | Positive tangential speed in meters per second. |
| `speedKts` | Conditional | Positive tangential speed in knots; mutually exclusive with `speedMps`. |
| `clockwise` | No | Rotation direction; defaults to counterclockwise. |
| `initialBearingDeg` | No | Initial compass bearing from the center; defaults to north. |

Exactly one positive speed field is required.

### Race-Track Pattern

Use a `raceTrack` object for two straight legs joined by semicircular turns:

```json
{
  "uid": "gotak-rq01",
  "callsign": "RQ01",
  "type": "a-f-A-M-F-Q",
  "hae": 1800,
  "raceTrack": {
    "offsetNorthMeters": 1000,
    "offsetEastMeters": 1700,
    "headingDeg": 60,
    "legLengthMeters": 3000,
    "turnRadiusMeters": 600,
    "speedKts": 70,
    "clockwise": true
  }
}
```

Leg length, turn radius, and speed must be positive. The offsets locate the pattern center. `headingDeg` controls the straight-leg direction, and `clockwise` mirrors the pattern.

Race-track fields:

| Field | Required | Meaning |
|---|---:|---|
| `offsetNorthMeters` | No | Northward pattern-center offset; defaults to `0`. |
| `offsetEastMeters` | No | Eastward pattern-center offset; defaults to `0`. |
| `headingDeg` | No | True heading of the straight legs; defaults to north. |
| `legLengthMeters` | Yes | Positive length of each straight leg. |
| `turnRadiusMeters` | Yes | Positive radius of each semicircular turn. |
| `speedMps` | Conditional | Positive ground speed in meters per second. |
| `speedKts` | Conditional | Positive ground speed in knots; mutually exclusive with `speedMps`. |
| `clockwise` | No | Mirrors the pattern when true; defaults to `false`. |

Exactly one positive speed field is required.

### Sensor Field of View

Add a `sensor` object to any track:

```json
{
  "sensor": {
    "fovDeg": 30,
    "rangeMeters": 8000,
    "azimuthOffsetDeg": 0
  }
}
```

`fovDeg` and `rangeMeters` must be positive. Sensor azimuth follows the track's current course plus `azimuthOffsetDeg`, including while the track turns through an orbit or race-track pattern.

Sensor fields:

| Field | Required | Meaning |
|---|---:|---|
| `fovDeg` | Yes | Positive horizontal field of view in degrees. |
| `rangeMeters` | Yes | Positive sensor range in meters. |
| `azimuthOffsetDeg` | No | Offset from current track course; defaults to `0`. |

## Connection Behavior

Enrollment uses port `8446` and intentionally disables server-certificate verification while bootstrapping trust. Because username/password credentials and the returned CA chain cross that unauthenticated TLS channel, enroll only over a trusted network or after establishing the server's identity out of band. The returned CA chain is then used to verify the TAK server on the port `8089` mTLS stream.

Enrollment and stream connection attempts time out after 30 seconds. Individual CoT writes time out after 10 seconds. Keys and certificates remain in memory and are not written to disk.

## Troubleshooting

- **Missing required fields:** supply `server`, `username`, and `password` through flags or `.env`.
- **Unknown location:** use an exact built-in location name or omit `-location` to choose interactively.
- **No scenario menu:** run from the repository root so `scenarios/` can be discovered, or pass `-scenario` explicitly.
- **Enrollment failure:** confirm credentials and access to port `8446`.
- **Stream connection failure:** confirm access to port `8089` and that enrollment returned a valid CA chain and client certificate.
- **Scenario parse failure:** check field names against the schema above, required identities, unique UIDs, positive orbit/race-track/sensor values, exclusive speed units, and that `orbit` and `raceTrack` are not both set.
