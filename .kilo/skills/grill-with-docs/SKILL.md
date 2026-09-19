---
name: grill-with-docs
description: Stress-test a plan or design while capturing resolved domain vocabulary and durable architecture decisions. Use when the user wants a grilling interview whose conclusions should update CONTEXT.md or ADRs as they crystallize.
---

# Grill with docs

Compose two installed workflows:

1. Invoke `grilling` to map and interview the decision tree.
2. Invoke `domain-modeling` for terminology, scenario testing, glossary updates, and qualifying ADRs.

During the interview, identify resolved terms and durable decisions immediately, but keep discussion separate from persistence. For each proposed document change, show the destination and exact content and obtain the authorization required by `domain-modeling`. One document approval does not authorize later edits, implementation, tracker mutations, or Git operations.

If either skill is unavailable, explain the limitation and follow its behavior directly only when its full instructions are already available; never claim it was invoked.