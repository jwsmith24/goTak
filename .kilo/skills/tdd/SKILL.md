---
name: tdd
description: Implement behavior test-first in vertical red-green-refactor cycles. Use when the user requests TDD, red-green-refactor, test-first feature work, a regression test, or integration-style behavior tests.
---

# Test-driven development

Read repository instructions, domain vocabulary, ADRs, test configuration, and nearby tests. Agree the public seam and intended behavior before editing. Preview the first slice and obtain implementation authorization; preserve unrelated and concurrent work.

## Cycle

Work in one narrow vertical slice at a time:

1. **Red:** write one behavior test at an agreed seam. Use an independent expected value. Run the smallest command that proves it fails for the intended reason; a compile error or unrelated failure is not red.
2. **Green:** write only enough production code to satisfy that behavior. Run the focused test and observe it pass.
3. **Small refactor:** make a behavior-preserving cleanup justified by the completed slice, or explicitly skip it. Keep the test green throughout. Do not add new behavior during refactoring.
4. Repeat with the next tracer-bullet behavior.

Tests specify observable behavior through public interfaces. Avoid private-method tests, internal call assertions, tautological expectations, snapshots without a reviewed oracle, and horizontal batches of imagined tests. Read [tests.md](tests.md) and [mocking.md](mocking.md) when those branches arise.

At the end, run focused checks and the relevant broader suite, inspect the complete diff, and perform a final review. Report commands and actual results. A passing cycle does not authorize commits, pushes, tracker updates, or publication.