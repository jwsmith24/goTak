---
name: handoff
description: Transfer current work to a fresh agent or Cline task in a compact continuation document. Use only when the user explicitly requests transfer to a new agent or session.
---

# Handoff

Create a compact Markdown handoff so a fresh agent or Cline task can continue the current work without replaying the full conversation.

Use this skill only when the user explicitly requests a handoff or session transfer.

## Process

1. Determine the intended focus of the next session from the user's request. If no focus is supplied, use continuation of the current task.
2. Resolve the workspace root and create `.kilo/handoff/` there when needed. Save the handoff as `.kilo/handoff/handoff-YYYYMMDD-HHMMSS.md`, using the local date and time so repeated handoffs do not collide.
3. Summarize only the state needed to continue:
   - intended outcome and scope;
   - current status;
   - decisions made and their rationale;
   - relevant workspace, repository, branch, revision, and diff identity when known;
   - files, URLs, issues, commands, and other evidence already inspected;
   - changes already made;
   - checks run and their actual results;
   - unresolved questions, blockers, risks, and untested areas;
   - the smallest useful next action.
4. Do not duplicate content already captured in specs, plans, ADRs, issues, commits, diffs, or other durable artifacts. Reference each by path, revision, identifier, or URL instead.
5. Include a `Suggested skills` section only when relevant skills are actually available in the current Cline setup. Name each skill and explain why it may help. Instruct the next agent to invoke it with Cline's `use_skill` tool or its slash command. Do not invent skills.
6. Redact credentials, secrets, tokens, private keys, protected data, and unnecessary personally identifiable information. Use `<REDACTED>` and retain only enough context to continue safely.
7. If a fact, completed action, or successful check cannot be supported by the conversation or an artifact, label it unknown or unverified rather than inferring it.
8. Report the handoff's absolute path and a one-sentence summary of its intended next session.
9. End the response with a fenced `text` code block containing a self-contained prompt the user can copy into a new task. Direct the follow-on agent to read the handoff at its absolute path and continue the intended next session. Place nothing after the code block.

Creating the handoff does not authorize continuing the work, editing project files, committing, pushing, publishing, deploying, or taking external actions.
