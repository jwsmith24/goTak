# Mocking reference

Prefer real in-process collaborators and faithful local substitutes. Use controlled adapters at true external seams such as third-party APIs, time, randomness, or unavailable infrastructure. Inject narrow behavior-specific ports rather than a generic conditional fetcher.

Mocking an owned internal module usually couples the test to implementation structure. Reconsider the seam before proceeding. Never let tests contact production services or use real credentials.