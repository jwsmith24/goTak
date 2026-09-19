---
name: to-tickets
description: Break an approved plan, specification, or discussion into dependency-aware tracer-bullet tickets for a configured GitLab or local Markdown tracker. Use when the user asks for ticket decomposition, vertical slices, blockers, or an implementation backlog.
---

# To tickets

1. Read the complete source material, repository conventions, domain docs, ADRs, and relevant code. Treat existing tracker content as read-only evidence.
2. Draft narrow, complete, independently verifiable vertical slices. Each ticket must fit a fresh agent context and state title, delivered behavior, acceptance criteria, blockers, and out-of-scope boundaries.
3. Use a preparatory refactor ticket only when it makes later behavior changes safer. For one mechanical change with a repository-wide blast radius, use expand → migration batches → contract, with blockers preserving a green state.
4. Present the numbered dependency graph. Ask whether granularity and blocking edges are correct; revise until approved.
5. Preview the complete publication batch and exact destinations. For local Markdown, create one `.scratch/<feature>/issues/NN-slug.md` file per ticket in dependency order. For GitLab, create issues in dependency order and use documented blocking links when available. Include proposed labels and parent references.
6. Obtain explicit authorization before any local writes or remote mutations. If identifiers are needed for second-pass links, preview that second mutation batch after creation. Preserve concurrent edits and stop if tracker state has changed materially.

Never modify or close the source/parent item automatically. Ticket approval does not authorize assignment, implementation, Git operations, or publication elsewhere.

## Ticket body

```markdown
## What to build
<End-to-end behavior from the user's perspective.>

## Acceptance criteria
- [ ] <Independent observable criterion>

## Blocked by
<Named ticket references or “None”.>

## Out of scope
<Explicit boundary.>
```