---
name: diagnosing-bugs
description: Diagnose hard bugs, failures, and performance regressions with a tight reproducible feedback loop, ranked falsifiable hypotheses, targeted instrumentation, and regression verification. Use when asked to diagnose or debug broken, failing, throwing, flaky, or slow behavior.
---

# Diagnosing bugs

Read relevant repository instructions, domain docs, ADRs, code, and tests. Redact secrets, authentication material, protected data, and unnecessary personal data from commands, artifacts, and reports. Keep credentials in environment variables. Obtain explicit authorization before adding harnesses, instrumentation, tests, or fixes.

## 1. Build a tight loop

Create a fast pass/fail signal for the exact reported symptom. Prefer, in order: focused test, HTTP/CLI script, browser automation, captured-trace replay, minimal harness, fuzz loop, bisection, or differential comparison. Use [scripts/hitl-loop.template.sh](scripts/hitl-loop.template.sh) only as a user-run last resort; interactive commands cannot be delegated through a non-interactive tool.

For flaky bugs, raise and measure the reproduction rate. If no red-capable loop is possible, report attempts and request a redacted artifact, environment access, or separately authorized temporary instrumentation. Do not form a confident cause without a loop.

## 2. Reproduce and minimize

Run the loop and confirm it catches the user's symptom. Remove inputs, callers, configuration, data, and steps one at a time until every remainder is load-bearing.

## 3. Rank hypotheses

Generate three to five falsifiable hypotheses before testing. For each state the predicted observation. Show the ranking to the user; incorporate supplied knowledge, then test one variable at a time.

## 4. Instrument narrowly

Prefer debugger/REPL inspection, then targeted boundary logs. Tag temporary instrumentation with one unique `[DEBUG-...]` prefix. For performance, establish a baseline and use profiling or query plans before changing code.

## 5. Fix and lock down

At the correct seam, turn the minimized reproduction into a failing regression test, observe red, apply the smallest supported fix, observe green, and rerun the original loop. If no correct seam exists, document that architectural finding rather than adding a misleading test.

## 6. Clean up and review

Remove authorized temporary instrumentation and throwaway artifacts, rerun focused and relevant broader checks, inspect the diff, and report the supported cause and residual risk. Cleanup that deletes files requires preview and approval. Commits, pushes, production instrumentation, browser actions, and remote changes require separate authorization.