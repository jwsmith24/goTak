# goTak Agent Guide

## Product Intent

goTak is a lightweight Go CLI for development and testing against a TAK server. It enrolls for a client certificate, opens an mTLS Cursor-on-Target (CoT) stream, and repeatedly sends simulated track positions until interrupted.

Keep the tool small and easy to build. Prefer the Go standard library and existing code over new packages. Add an external dependency only when the required behavior cannot be implemented reasonably without it, and document the justification in the change.

## Code Map

- `cmd/gotak`: composition root and CLI lifecycle.
- `internal/config`: flags and optional `.env` fallback parsing.
- `internal/menu`: terminal and numbered-list location/scenario selection.
- `internal/location`: named scenario origins.
- `internal/scenario`: JSON schema, defaults, validation, and unit conversion.
- `internal/sim`: straight, orbit, and race-track motion plus the tick/send loop.
- `internal/cot`: CoT XML event and sensor detail generation.
- `internal/enroll`: TAK certificate enrollment and generated key/CSR handling.
- `internal/stream`: verified mTLS connection and CoT writes.
- `scenarios`: example relative-position scenarios used by the menu.

## Behavioral Invariants

- Flags override `.env`; `.env` is optional and must never be committed because it can contain credentials.
- Only an explicitly supplied `-location` or `-scenario` skips its menu. Values from `.env` remain fallbacks and do not suppress interactive selection.
- Menus support raw-terminal arrow navigation and a non-terminal numbered fallback. Preserve both paths and the shared buffered stdin reader across consecutive prompts.
- Scenario coordinates are east/north meter offsets from a selected named origin, not fixed latitude/longitude. Keep scenarios portable between locations.
- Keep each scenario track to one motion model: straight, orbit, or race track. Scenario parsing owns validation, defaults, and knots-to-meters-per-second conversion.
- On every tick, each track advances first, then emits one CoT event. Sensor azimuth follows the track's current course plus its configured offset.
- Enrollment on port `8446` intentionally skips server verification only while bootstrapping trust. The CoT stream on port `8089` must use the issued client certificate and verify the server with the CA chain returned by enrollment.
- Keep secrets and generated private keys in memory; do not log or persist them.

## Change Discipline

- Make the smallest coherent change and keep package boundaries aligned with the code map above.
- Use standard-library interfaces and small local abstractions for test seams; avoid frameworks, dependency injection containers, and speculative compatibility layers.
- Add or update focused table-driven tests beside changed behavior. Network code should use local test servers, fake connections, or narrow interfaces rather than a live TAK server.
- If scenario behavior or CLI usage changes, update `README.md`, `.env.example`, and shipped scenario files when applicable.
- Format Go changes with `gofmt`.

## Validation

Run from the repository root:

```sh
go test ./...
go vet ./...
go build ./cmd/gotak
```

Use the Go version declared in `go.mod`. A live TAK server is not required for the automated test suite.
