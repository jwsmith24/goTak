---
name: improve-codebase-architecture
description: Scan a codebase for high-value deepening opportunities, present a visual architecture report, and explore one selected candidate. Use when the user asks for architecture improvement, deep-module opportunities, testability, or easier codebase navigation.
---

# Improve codebase architecture

Use `codebase-design` vocabulary and project domain language.

1. **Scope.** Follow the user's named area. Otherwise use read-only history to identify recently changing hot spots before widening the scan. Read relevant domain docs and ADRs.
2. **Explore.** Use optional read-only subagents for independent areas, or inspect directly. Look for shallow modules, scattered concept knowledge, leaking seams, tests forced through internals, and repeated caller complexity. Apply the deletion test. Prefer demonstrated friction over theoretical cleanup.
3. **Report.** Draft a self-contained offline HTML report following [HTML-REPORT.md](HTML-REPORT.md). Each candidate includes files/modules, observed friction, proposed deepening, test impact, before/after diagram, ADR conflicts, and strength: `Strong`, `Worth exploring`, or `Speculative`. Recommend one candidate. Do not design its final interface yet.
4. **Persist.** Show the proposed temporary absolute path and report summary, then obtain explicit authorization before writing. Use a collision-resistant file in the OS temporary directory. Opening it in a browser requires separate authorization.
5. **Explore one candidate.** After the user chooses, invoke `grilling` for the decision tree and `domain-modeling` for authorized glossary/ADR capture. Use `codebase-design`'s design-it-twice process when competing interfaces would help.

The report is analysis, not authorization to refactor code, delete tests, create branches, or perform Git operations.