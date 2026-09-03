# goTak

A terminal tool for simulating air tracks on a development TAK server.

Given a server IP, username, and password, it:

1. Enrolls for a client certificate via the TAK server's
   username/password-authenticated enrollment endpoint (no preconfigured
   trust store required).
2. Opens an mTLS connection to the server's CoT streaming port using that
   certificate.
3. Simulates one or more moving air tracks, sending a CoT position update
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
| `-scenario`  | Path to a JSON scenario file (optional; see below)                  |

`-server`, `-username`, and `-password` are required; the app reports any
that are missing. `-scenario` is optional — without it, the app simulates
a single default track.

### Scenario files

A scenario file describes one or more tracks to simulate and how often to
update them:

```json
{
  "tickIntervalSeconds": 2,
  "tracks": [
    {
      "uid": "gotak-austin-eagle01",
      "callsign": "EAGLE01",
      "type": "a-f-A",
      "lat": 30.2747,
      "lon": -97.76,
      "hae": 1500,
      "courseDeg": 90,
      "speedMps": 120
    },
    {
      "uid": "gotak-austin-eagle02",
      "callsign": "EAGLE02",
      "lat": 30.26,
      "lon": -97.7404,
      "hae": 2000,
      "courseDeg": 0,
      "speedMps": 100
    }
  ]
}
```

- `uid` and `callsign` are required and must be unique per track.
- `type` is the CoT type (e.g. `a-f-A` for friendly air); defaults to
  `a-f-A` when omitted.
- `lat`/`lon` are decimal degrees; `hae` is height above the ellipsoid in
  meters; `courseDeg` is true course in degrees clockwise from north;
  `speedMps` is ground speed in meters/second.
- `tickIntervalSeconds` controls how often every track's position updates;
  defaults to 2 seconds when omitted.

[`scenarios/austin-capitol.json`](scenarios/austin-capitol.json) ships
with the repo: two air tracks crossing paths near the Texas Capitol in
downtown Austin. Run it with:

```sh
go run ./cmd/gotak -server 192.168.1.50 -username dev -password devpass \
  -scenario scenarios/austin-capitol.json
```

#### Orbiting tracks

A track can loop around a fixed point instead of flying a straight
course by giving it an `orbit` object instead of `lat`/`lon`/`courseDeg`/
`speedMps`:

```json
{
  "uid": "gotak-austin-helo01",
  "callsign": "HELO01",
  "type": "a-f-A-M-H",
  "hae": 300,
  "orbit": {
    "centerLat": 30.2747,
    "centerLon": -97.7404,
    "radiusMeters": 800,
    "speedMps": 35,
    "clockwise": true,
    "initialBearingDeg": 0
  }
}
```

- `centerLat`/`centerLon` are the orbit's center point; `radiusMeters` is
  the orbit radius; `speedMps` is the tangential ground speed (both must
  be positive).
- `clockwise` sets rotation direction (default `false`, counterclockwise).
- `initialBearingDeg` places the track's starting position on the circle,
  as a compass bearing from the center (default `0`, due north of center).
- The track's reported course and speed update every tick to reflect its
  current heading around the circle.

[`scenarios/austin-capitol-helicopters.json`](scenarios/austin-capitol-helicopters.json)
ships two rotary-wing helicopters (`HELO01`, `HELO02`) orbiting the
Capitol on the same circle, 180° offset from each other so they stay on
opposite sides of the loop, at slightly different altitudes:

```sh
go run ./cmd/gotak -server 192.168.1.50 -username dev -password devpass \
  -scenario scenarios/austin-capitol-helicopters.json
```

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
GOTAK_SCENARIO=scenarios/austin-capitol.json
```

Then just run:

```sh
go run ./cmd/gotak
```

Any flag you do pass on the command line overrides the matching value from
`.env`. `.env` is gitignored — never commit real credentials.

## Running the tests

```sh
go test ./...
```
