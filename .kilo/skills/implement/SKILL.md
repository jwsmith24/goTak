---
name: implement
description: Implement an approved specification, ticket, or clearly bounded change with incremental validation and final review. Use when the user explicitly asks to implement defined work rather than only plan or analyze it.
---

# Implement

1. **Load the contract.** Read the full spec or tickets, repository instructions, domain docs, ADRs, relevant code, tests, and configured tracker context. Separate requirements, acceptance criteria, dependencies, and out-of-scope work. Ask only about blocking ambiguity.
2. **Plan and preview.** Propose vertical slices, affected interfaces/files, test seams, commands, migrations, and risks. Identify destructive, external, privileged, generated, or dependency-changing steps. Obtain explicit authorization for the implementation batch before editing.
3. **Implement incrementally.** Preserve unrelated and concurrent work. Re-read targets before edits. Use `tdd` at agreed seams where behavior can be tested: red, green, small behavior-preserving refactor. Otherwise make the smallest coherent slice and add proportionate verification.
4. **Validate continuously.** Run focused tests and static checks after each slice. Run relevant broader checks at the end. Never claim an unrun check passed.
5. **Review.** Inspect the final diff and invoke `code-review` against the agreed fixed point. Address in-scope findings with the same edit authorization; report deferred findings honestly.
6. **Report.** Map completed work to acceptance criteria, list changed files and actual check results, and identify residual risks.

Implementation approval does not authorize dependency installation with side effects, database or production mutations, remote tracker changes, commits, branches, pushes, deployment, or publication. Preview and request those separately.