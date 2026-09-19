---
name: engineering-workflows
description: Route an engineering or writing task to the smallest useful sequence of installed workflow skills. Use when the user asks which skill or flow fits, how the skills connect, or what to do next.
---

# Engineering workflows

Recommend, do not automatically invoke, the smallest useful path. A skill does not inherit authorization from an earlier phase: preview and approval requirements still apply at every write, Git, remote, browser, publication, or destructive boundary.

## Main engineering path

1. **Clarify:** use `grilling` for a stateless decision-tree interview or `grill-with-docs` when durable vocabulary and ADR candidates should be captured. Use `domain-modeling` for focused language and architectural decisions. Use `research` for authoritative external evidence and `to-questionnaire` when another person holds the missing knowledge.
2. **Answer runnable design questions:** use `prototype` for one disposable logic, state, data-shape, or UI question. Use `handoff` only when work must move to another session, agent, directory, or person.
3. **Plan:** use `to-spec` after discovery is resolved. For a multi-part build, use `to-tickets` for dependency-aware vertical slices. Use `wayfinder` instead when the destination is known but the route is still too uncertain for a fixed plan.
4. **Build:** use `implement` for an approved bounded change; it can apply `tdd` in red-green-refactor slices. Use `tdd` directly for a concrete behavior that does not need a full specification.
5. **Review:** use `code-review` against an explicit fixed point. Review findings do not authorize fixes, commits, pushes, or merge-request actions.

## On-ramps and specialist paths

- Incoming raw issues: `triage` → `implement` after the issue is accepted and authorized. Do not re-triage tickets already produced by `to-tickets`.
- Hard failure: `diagnosing-bugs` → `tdd` or `implement` once a cause and approved fix are clear.
- Architecture: `improve-codebase-architecture` → `codebase-design` → approved planning or implementation.
- Repository conventions: `setup-engineering-workspace` before tracker-aware workflows when conventions are absent or changing.
- In-progress merge/rebase conflict: `resolving-merge-conflicts`; keep Git continuation and other mutations explicitly authorized.
- Human-only setup or cutover: `wizard` to generate a reviewable interactive guide; generation does not authorize execution.
- Recurring process: `workflow-design` to discover and specify the loop.
- Session learning: `agent-retrospective` to recommend evidence-backed environment improvements.

## Writing path

Use `writing-fragments` to explore and collect raw material, then `writing-shape` to mine a fixed pile into an article one agreed block at a time. Use `writing-for-agents` for skills, AGENTS.md, rules, and linked agent-facing references. Use `wait-what` only to re-explain the immediately preceding response.

## Phase boundaries

Change context only at a natural boundary, not mid-investigation or mid-slice. Prefer, in order:

1. continue when the next phase needs the current reasoning and capacity remains;
2. start fresh when prior context is irrelevant;
3. use `handoff` when portable state must travel;
4. use an optional read-only research subagent only for tightly scoped independent evidence, with direct inspection as fallback;
5. otherwise create a concise user-reviewed continuation summary.

Summaries are lossy. Preserve primary evidence, decisions, unresolved questions, validation state, and authority boundaries. Never advertise or invoke a skill that is not installed.