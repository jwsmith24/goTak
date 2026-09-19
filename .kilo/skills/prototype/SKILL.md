---
name: prototype
description: Build a disposable prototype to answer one design question about business logic, state transitions, data shape, or user-interface structure. Use when the user wants to test whether a model feels right or compare materially different UI directions.
---

# Prototype

A prototype is disposable code that answers one explicit question.

Choose the branch: use [LOGIC.md](LOGIC.md) for logic, state, or data-shape questions; use [UI.md](UI.md) for appearance and interaction structure. If ambiguous, inspect surrounding code and ask.

Before writing, state the question, proposed location, files, run method, disposal plan, and any data or network effects. Preview the change and obtain explicit authorization. Re-read targets before editing and preserve unrelated work.

Both branches:

- mark artifacts visibly as prototypes and follow existing project conventions;
- keep startup trivial and state inspectable;
- use memory, fixtures, or an explicitly disposable local store by default;
- avoid real mutations, production data, secrets, and unnecessary dependencies;
- optimize for learning rather than production completeness;
- record the answer after evaluation.

When the question is answered, preview how validated behavior will be reimplemented or folded into production-quality code and how prototype artifacts will be removed or retained. Those edits and deletions require approval. Branch creation, commits, browser opening, publication, and tracker updates are separate actions and are never automatic.