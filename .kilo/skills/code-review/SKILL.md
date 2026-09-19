---
name: code-review
description: Review a branch, merge request, worktree, or diff from a fixed point along separate standards and specification axes. Use for code review, change review, merge-request review, or “review since” a commit, tag, branch, or merge-base.
---

# Code review

This is a read-only review. Do not edit, stage, commit, comment remotely, or publish findings unless separately requested and authorized.

## 1. Pin the comparison

Resolve the user-supplied fixed point. If absent, ask for one or propose the repository's merge-base and confirm it. Capture the exact diff and commit range once. Include working-tree changes when the requested scope requires them. Stop on a bad ref or empty scope.

## 2. Gather sources

Find the originating spec from supplied paths, configured tracker references, commit messages, or matching repository specs. Read repository instructions, coding standards, domain docs, ADRs, and tool configuration. If no spec exists, say so rather than inventing one.

## 3. Run independent axes

Use parallel read-only subagents when available; otherwise perform two independent direct passes:

- **Standards:** report documented-standard violations with citations. Also flag possible Mysterious Name, Duplicated Code, Feature Envy, Data Clumps, Primitive Obsession, Repeated Switches, Shotgun Surgery, Divergent Change, Speculative Generality, Message Chains, Middle Man, and Refused Bequest. Repository rules override these heuristics; omit issues enforced by successful tooling.
- **Spec:** report missing or partial requirements, scope creep, and behavior that appears to implement a requirement incorrectly. Quote the requirement for each finding.

Verify every finding against the diff and source. Include file and line/hunk locators, impact, evidence, and a concrete correction. Avoid style preferences without a governing standard and avoid claims about code outside the reviewed scope.

## 4. Report

Keep `## Standards` and `## Spec` separate; do not merge their rankings. Order findings by severity within each axis, then list uncertainties and checks reviewed. End with counts and the worst issue in each axis. If no findings remain, state that explicitly with residual test gaps.