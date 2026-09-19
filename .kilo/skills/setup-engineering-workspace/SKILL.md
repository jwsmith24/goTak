---
name: setup-engineering-workspace
description: Configure a repository for the portable engineering skills: issue-tracker conventions, triage vocabulary, and domain-document layout. Use once before tracker-aware workflows or when those repository conventions change.
---

# Setup engineering workspace

Explore, preview, confirm, then write repository-local configuration.

## 1. Inspect

Read repository instructions, remotes, existing `docs/agents/`, domain docs, ADR directories, `.scratch/`, available skills, and monorepo signals. Use read-only commands only. Prefer an existing root `AGENTS.md`; if repository instructions designate another agent file, follow them. If no suitable file exists, ask which portable instruction file to create.

## 2. Resolve configuration

- **Tracker:** recommend GitLab for a GitLab remote or local Markdown when no supported remote is used. Also support a user-described tracker. Record commands as examples, never standing authorization.
- **Triage:** when `triage` is installed, recommend canonical roles `needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, and `wontfix`; map to existing labels when needed.
- **Domain docs:** default to one root `CONTEXT.md` and `docs/adr/`. Offer `CONTEXT-MAP.md` only when genuine multiple-context signals exist.

## 3. Preview

Show complete drafts for the `## Agent skills` block and every proposed `docs/agents/*.md` file. Base them on [issue-tracker-gitlab.md](issue-tracker-gitlab.md), [issue-tracker-local.md](issue-tracker-local.md), [triage-labels.md](triage-labels.md), and [domain.md](domain.md). For another tracker, capture read/list/create/comment/label/close operations and their authorization boundaries without inventing commands.

Obtain explicit authorization for the complete file batch. This setup never provisions labels, opens a browser, authenticates a CLI, or mutates a remote.

## 4. Write safely

Re-read every destination immediately before editing. Update an existing `## Agent skills` block in place; preserve surrounding content. Stop and reconcile unexpected changes. Omit triage configuration when `triage` is unavailable.

Report changed files and explain that future local writes, remote mutations, and Git operations still require their own preview and authorization.