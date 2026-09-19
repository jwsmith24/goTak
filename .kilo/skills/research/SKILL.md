---
name: research
description: Investigate a focused question against high-trust primary sources and produce a cited Markdown research note. Use when the user asks for research, authoritative documentation or API facts, or evidence-gathering that should be saved for later work.
---

# Research

Investigate a focused question and produce one evidence-backed Markdown note.

1. **Frame.** Restate the question, scope, intended reader, and what decision the research should support. Resolve only ambiguities that would materially change the investigation.
2. **Locate evidence.** Prefer primary sources: official documentation, specifications, source code, release notes, and first-party APIs. Trace important secondary claims back to the source that owns them. Record versions and retrieval dates when facts can drift.
3. **Investigate.** Use native read-only subagents for independent codebase questions when available and worthwhile; otherwise inspect directly. Subagents are optional researchers, not writers: they cannot browse or modify files. Verify their findings against cited evidence.
4. **Synthesize.** Separate sourced facts, reasoned conclusions, uncertainties, and recommendations. Cite each consequential claim with a stable URL, repository path and revision, or other precise locator. Quote sparingly and never fabricate quotations.
5. **Preview.** Choose the repository's existing research-note location, or propose a sensible path when none exists. Re-read an existing destination, then show the proposed absolute path and complete draft. Redact credentials, protected data, unnecessary personal data, and confidential material.
6. **Persist.** Write only after explicit authorization. Preserve concurrent user edits; if the destination changed after preview, stop and reconcile. Report the path and remaining evidence gaps.

Research authorization does not authorize code changes, git operations, tracker updates, publication, or other external actions.