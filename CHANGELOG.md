# Changelog

Notable changes are documented here in reverse chronological order. Releases use [Semantic Versioning](https://semver.org/).

## [0.2.0] - 2026-09-19

### Fixed

- Restored custom latitude/longitude entry in the location menu, including validation, local persistence, and later selection by saved name.
- Reject scenario tracks that define both `orbit` and `raceTrack` instead of silently choosing one motion model.
- Make certificate enrollment and CoT stream connection cancellable by SIGINT or SIGTERM.
- Bound enrollment and stream connection attempts to 30 seconds.
- Limit enrollment responses to 1 MiB for TLS configuration XML and 4 MiB for certificate-signing JSON.
- Apply a refreshed 10-second deadline to every CoT stream write so an unresponsive peer cannot block indefinitely.

### Changed

- Extracted the CLI lifecycle into testable orchestration while keeping signal interception after menu and scenario preparation.
- Scenario JSON parsing is now strict: field names are exact and case-sensitive, while unknown fields, duplicate fields, and additional trailing JSON values are rejected.
- Added focused coverage for in-flight enrollment cancellation and timeout, TLS handshake cancellation, oversized enrollment responses, and blocked stream writes.
- Kept signal interception after interactive selection so normal SIGINT/SIGTERM termination remains available while menus await input.
- Added `AGENTS.md` with durable architecture, security, dependency, testing, and documentation guidance.
- Added `USER_GUIDE.md` and condensed `README.md` into a project overview and quick-start entry point.

## [0.1.0] - 2026-09-10

### Added

- Provided a lightweight Go CLI that enrolls with a TAK server and streams simulated Cursor-on-Target positions over mTLS.
- Supported straight, circular-orbit, and race-track motion with optional sensor field-of-view details.
- Added reusable scenarios positioned by east/north offsets from selectable named origins.
- Added interactive terminal menus with a numbered fallback for non-terminal input.
- Supported command-line configuration with optional `.env` fallbacks.
- Covered configuration, menus, scenarios, motion, CoT generation, enrollment, and mTLS streaming with local automated tests.
