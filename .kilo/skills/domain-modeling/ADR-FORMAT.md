# ADR format

Store ADRs in the configured context's `docs/adr/` directory as `NNNN-slug.md`. Scan existing ADRs and increment the highest number; never overwrite or reuse a number silently.

```markdown
# <Short decision title>

<One to three sentences stating the context, decision, and why.>
```

Add status, considered options, or consequences only when they carry information a future reader needs. Keep an ADR compact. Its value is the durable decision and rationale, not template completion.