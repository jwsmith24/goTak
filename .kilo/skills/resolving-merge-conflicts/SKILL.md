---
name: resolving-merge-conflicts
description: Resolve an in-progress Git merge or rebase conflict while preserving the intent of both sides. Use when the repository has unmerged paths or the user explicitly asks to resolve merge or rebase conflicts.
---

# Resolving merge conflicts

1. **Inspect read-only.** Determine whether a merge, rebase, or cherry-pick is active. Read status, history, unmerged stages, conflict hunks, and repository instructions. Do not abort, reset, checkout over work, stage, commit, or continue.
2. **Recover intent.** For each side, inspect the relevant commits, blame/history, tests, docs, configured issue tracker, and available local evidence. Remote reads are allowed only when already authenticated and non-mutating. If intent remains uncertain, say what is missing.
3. **Propose.** Explain each conflict, both intents, the proposed resolution, trade-offs, files to change, and checks to run. Preserve both intents when compatible. When they conflict, follow the operation's stated goal and existing requirements; do not invent unrelated behavior.
4. **Authorize edits.** Obtain explicit approval for the proposed resolution batch. Immediately before editing, re-read every target and compare it with the inspected state. Stop if it changed unexpectedly.
5. **Resolve and validate.** Remove conflict markers, preserve unrelated work, and run the project's focused checks. Fix only failures caused by the resolution. Review the resulting diff and confirm that no unmerged entries or conflict markers remain.
6. **Finish separately.** Present the final diff summary, check results, and exact proposed Git commands. Staging files and running `merge --continue`, `rebase --continue`, `cherry-pick --continue`, or a commit each require explicit authorization. Never push automatically.

If a safe resolution cannot be inferred, leave the conflict active and ask for the missing decision rather than choosing silently.