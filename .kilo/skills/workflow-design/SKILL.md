---
name: workflow-design
description: Discover and specify a recurring human or technical workflow through a decision-tree interview. Use when the user wants to identify delegable loops, define triggers and checkpoints, or create an implementation-ready workflow specification.
---

# Workflow design

A **loop** is a recurring pattern. A **workflow** is its durable specification. Invoke `grilling` to expose the decision tree and use this vocabulary only where it helps:

- **Trigger:** the event or schedule that starts a run; prefer an event when it is the true cause.
- **Checkpoint:** a human verification or decision point.
- **Push right:** perform safe preparatory work before the checkpoint so the human reviews once, late, with useful context.
- **Brief:** the compact, decision-ready checkpoint summary with links to evidence or artifacts.

Do not assume a workflow needs AI, automation, a schedule, or a checkpoint. Map the real current loop, failure modes, inputs, outputs, systems, data sensitivity, authority boundaries, retries, idempotency, observability, ownership, and recovery. A workflow is implementation-ready only when an implementer can build it without guessing consequential behavior.

Use existing workspace conventions for workflow specs and notes. Before creating or editing a spec, show the absolute destination and complete proposed content and obtain explicit authorization. Re-read the destination before writing and preserve concurrent edits. Rewrites and deletions require a separate preview and confirmation. Specification does not authorize implementation, provisioning, scheduling, external mutations, Git operations, or publication.