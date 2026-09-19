---
name: domain-modeling
description: Build or sharpen a project's domain vocabulary and durable architecture decisions. Use when resolving terminology, creating or editing CONTEXT.md or CONTEXT-MAP.md, or deciding whether to record an ADR.
---

# Domain modeling

Actively sharpen the model rather than merely reading its vocabulary.

## Explore

Read the relevant code, `CONTEXT.md`, `CONTEXT-MAP.md`, ADRs, and `docs/agents/domain.md` when present. Use a root `CONTEXT.md` for a single context. In a multi-context repository, follow `CONTEXT-MAP.md` to the relevant context and use root `docs/adr/` only for system-wide decisions.

During discussion:

- call out conflicts with established glossary language;
- replace vague or overloaded words with a proposed canonical term;
- test relationships with concrete edge-case scenarios;
- compare claims with code and surface contradictions;
- distinguish domain vocabulary from implementation details.

## Capture resolved terms

Use [CONTEXT-FORMAT.md](CONTEXT-FORMAT.md). `CONTEXT.md` is a glossary, not a specification or scratchpad. Create files lazily only after a term is resolved.

Preview each proposed glossary append or update, including its destination and surrounding text, and obtain explicit authorization before writing. Re-read the destination immediately before editing. Preserve user edits and stop on unexpected changes. One approval may cover a clearly previewed batch; later terms require another preview.

## Record durable decisions

Offer an ADR only when the decision is all three:

1. costly to reverse;
2. surprising without its context;
3. the outcome of a real trade-off.

Use [ADR-FORMAT.md](ADR-FORMAT.md). Preview the complete ADR, path, and any map change before requesting write authorization. Recording an ADR does not authorize implementation, git operations, or publication.