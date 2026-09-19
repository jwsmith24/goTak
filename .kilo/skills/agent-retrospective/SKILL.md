---
name: agent-retrospective
description: Review a coding-agent session and recommend evidence-backed improvements to navigation, checks, standards, instructions, tools, and information access. Use when the user asks for a retrospective on the current or identified engineering session.
---

# Agent retrospective

Invoke `writing-for-agents` before proposing agent-facing documents.

1. **Establish evidence.** Use the current session unless another is identified. Read only accessible session artifacts, diffs, commands, check output, repository instructions, CI, tool configuration, and relevant logs. Ask before accessing sensitive or out-of-workspace logs. Redact secrets and personal data.
2. **Reconstruct.** Distinguish observed facts from inference. Identify delays, mistakes, repeated searches, failed checks, expensive calls, unavailable evidence, and successful guardrails.
3. **Find improvements:**
   - navigation pointers for hard-to-find knowledge;
   - deterministic checks for mechanical failures;
   - reviewer standards for genuine judgment calls;
   - pruning or relocation of bloated/no-op steering instructions;
   - more efficient local tools and commands;
   - safe read-only access to missing information.
4. **Rank.** Order candidates by demonstrated impact, recurrence, confidence, implementation cost, and maintenance burden. Prefer fixing an existing broken or unwired check over adding a duplicate. Prefer deterministic enforcement over prose for mechanical rules.
5. **Report.** For each recommendation cite session evidence, explain the future behavior change, identify the target artifact, and state validation and rollback.

The retrospective is analysis only. Show exact proposed changes and obtain explicit authorization before editing files, installing tools, changing hooks/CI, accessing services, or performing Git/remote actions. Do not manufacture a recommendation when evidence shows no meaningful improvement.