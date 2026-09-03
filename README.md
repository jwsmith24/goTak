# goTak

A terminal tool for simulating air tracks on a development TAK server.

Given a server IP, username, and password, it:

1. Enrolls for a client certificate via the TAK server's
   username/password-authenticated enrollment endpoint (no preconfigured
   trust store required).
2. Opens an mTLS connection to the server's CoT streaming port using that
   certificate.
3. Simulates one moving air track, sending a CoT position update every 2
   seconds until you stop it.

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

| Flag         | Description                                  |
|--------------|-----------------------------------------------|
| `-server`    | TAK server IP address or hostname             |
| `-username`  | Username for certificate enrollment           |
| `-password`  | Password for certificate enrollment           |

All three are required; the app reports any that are missing.

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
