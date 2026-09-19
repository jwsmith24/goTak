---
name: wizard
description: Generate a reviewable interactive Bash guide for a manual procedure only a human can complete. Use for credential setup, unfamiliar third-party dashboards, infrastructure provisioning, one-off migrations, or cutovers—not for steps the agent can safely perform directly.
---

# Wizard

A wizard guides a human through a manual procedure. It is generated code with consequential capabilities, so design and review it before writing or running it.

## 1. Scope from evidence

Read repository instructions, environment examples, setup docs, service configuration, CI references, and migration state. For each stage identify:

- the exact human action and authoritative documentation;
- inputs and outputs;
- whether each value is secret, sensitive, or public;
- where a value may be persisted;
- local, remote, irreversible, and retry effects;
- verification, recovery, and safe restart behavior.

Never invent dashboard paths or commands. Do not request secret values in chat.

## 2. Preview the procedure

Show ordered stages, destinations, external systems, values captured, commands/actions, confirmation gates, and rollback limitations. Confirm scope with the user. Then preview the target script path and meaningful generated stage code and obtain explicit authorization before creating it.

## 3. Author safely

Copy [templates/wizard.sh](templates/wizard.sh) and replace only the marked stages. Keep strict mode. Use its helpers for instructions, hidden capture, URL confirmation, and `.env` writes. Each local write, browser open, remote mutation, provisioning step, migration, and irreversible action must have an in-script confirmation immediately before it. Commands must use argument arrays or direct quoted arguments, never `eval`. Do not print secrets or place them in command lines when stdin or a protected file is supported.

Generated wizards are ephemeral by default and must preserve existing configuration. Make reruns safe or state exactly where they are not. Do not add dependencies without approval.

## 4. Validate and hand off

Run `bash -n`; run ShellCheck when installed. Statically trace every captured value and mutation, and inspect file permissions and cleanup paths. Do not run the wizard end-to-end: it is interactive and consequential. Changing executable bits, running it, retaining it in the repository, linking it from docs, or performing Git operations each requires separate authorization.