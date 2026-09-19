# `CONTEXT.md` format

```markdown
# <Context name>

<One or two sentences describing this context.>

## Language

**Order**:
<One or two sentences defining the term.>
_Avoid_: Purchase, transaction
```

## Rules

- Pick one canonical word and list meaningful rejected synonyms under `_Avoid_`.
- Define what the concept is in one or two sentences; keep implementation details out.
- Include only concepts specific to this domain, not general programming terms.
- Add subheadings only when natural clusters emerge.

For multiple contexts, keep a root `CONTEXT-MAP.md` listing each context, its relative `CONTEXT.md` path, and relationships between contexts. Ask when the relevant context is ambiguous.