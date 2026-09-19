---
name: to-spec
description: Synthesize the current conversation and codebase evidence into a reviewable specification without reopening discovery. Use when the user asks to turn an already-resolved discussion into a spec, optionally for a configured GitLab or local Markdown tracker.
---

# To spec

Synthesize established decisions; do not restart the interview. Ask only when a missing fact makes the draft materially unsafe or contradictory.

1. Read relevant repository instructions, code, tests, domain docs, ADRs, prototypes, and configured tracker conventions.
2. Identify the highest practical public test seams, preferring existing seams. Present them for confirmation before finalizing testing decisions.
3. Draft:
   - `## Problem statement`
   - `## Solution`
   - `## User stories` — enough to cover meaningful behavior and edge cases, without padding
   - `## Implementation decisions`
   - `## Testing decisions`
   - `## Out of scope`
   - `## Further notes`
4. Describe durable behavior, interfaces, schemas, and constraints rather than line numbers or volatile file paths. Include a small prototype-derived state machine, reducer, schema, or type only when it records a settled decision more precisely than prose.
5. Show the complete draft and proposed destination. Publishing may mean a new GitLab issue or a local Markdown file according to `docs/agents/issue-tracker.md`; it is not implied by drafting.
6. Obtain explicit authorization for the exact write or remote mutation batch. Re-read local destinations before writing. Applying `ready-for-agent`, commenting, editing an existing issue, or creating a new issue must be included explicitly in the preview.

Report the resulting path or tracker URL/identifier. Do not modify or close a parent item, commit, push, or start implementation.