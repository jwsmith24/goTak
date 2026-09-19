---
name: wayfinder
description: Chart a large, uncertain effort as a destination, decision map, dependency frontier, and fog that clears over multiple sessions. Use when the work is too uncertain or large for one agent session and needs configured GitLab or local Markdown decision tickets.
---

# Wayfinder

Wayfinding resolves decisions until the route to a named destination is clear. It plans by default; implementation belongs after the map unless the map explicitly records an approved exception.

## Model

- **Map:** low-resolution index containing Destination, Notes, Decisions so far, Not yet specified, and Out of scope.
- **Decision ticket:** one question sized for one fresh session, typed `research`, `prototype`, `grilling`, or prerequisite `task`.
- **Frontier:** open, unblocked, unclaimed tickets.
- **Fog:** in-scope questions visible but not yet sharp enough to ticket.
- **Out of scope:** work deliberately beyond the destination; it never graduates from fog.

Refer to tickets by linked title, not bare identifiers. Keep each decision's detail in its ticket and only a one-line linked gist in the map.

## Chart a map

1. Invoke `grill-with-docs` to name the destination.
2. Explore breadth-first for the initial frontier and fog. If the whole route is already clear and fits one session, recommend a normal plan instead.
3. Draft the complete map, currently specifiable tickets, types, blockers, and tracker operations according to `docs/agents/issue-tracker.md`.
4. Preview and obtain explicit authorization before creating local files or GitLab issues. Create identifiers first; preview relationship/label mutations as a second batch when required. Do not launch implementation or create research branches.

## Work one ticket

Resolve at most one non-research ticket per session.

1. Read the map and query the current frontier. If the user names no ticket, recommend the first frontier item.
2. Preview any claim/assignment mutation and obtain authorization before applying it. Re-query immediately first; stop if another session claimed or changed it.
3. Resolve using the ticket type. Use optional read-only subagents for research only; the parent verifies and drafts persistence. Human-in-the-loop tickets never answer for the human.
4. Preview the answer, closure/status, map gist, fog graduation, new tickets, blockers, and invalidated-ticket changes as explicit local/remote mutation batches.
5. Apply only authorized batches, checking concurrent state between dependent passes. Never silently delete history; link superseded or out-of-scope decisions.

Tracker authorization does not authorize code implementation, provisioning, browser actions, secrets handling, Git operations, or publication elsewhere.