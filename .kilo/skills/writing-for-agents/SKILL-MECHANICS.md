# Skill mechanics

The skill-specific branch of [`writing-for-agents`](SKILL.md): what changes when the document is a Cline skill. Everything else about writing it is the universal reference in `SKILL.md`.

## Metadata and invocation

Every Cline skill is a directory containing `SKILL.md`. Its YAML frontmatter requires:

- `name`: exactly the same as the skill directory name.
- `description`: a model-facing context pointer, no more than 1024 characters, that states what the skill does and the distinct branches that should trigger it.

For every enabled skill, Cline loads the name and description as metadata. When a request matches the description, Cline can activate the skill and load its body. A human can also activate it explicitly with its slash command.

Cline does not support `disable-model-invocation` in skill frontmatter. Do not use that field or treat a description as human-only metadata. If a skill should not be available for automatic matching, disable it through Cline's Skills interface; a frontmatter-only user-invoked mode is not available.

Keep the description selective because every enabled skill spends metadata context. Put execution steps and detailed reference in the body, and disclose larger or branch-specific resources through relative links.

## Splitting skills

Split a skill when a distinct branch should trigger independently or when its body has separate paths that would otherwise create sprawl. Each new enabled skill adds an always-visible name and description, so the independent trigger must justify that context load.

Keep one skill when branches share the same trigger and most instructions. Use bundled reference files when only some branches need substantial detail.

## Discovery

Cline's slash-command suggestions provide a human-visible index of enabled skills. Prefer clear skill names and selective descriptions over a router skill created only to make other skills memorable.

If a router is still useful for a domain workflow, treat it as an ordinary model-visible Cline skill: give it a truthful description, account for its metadata load, and name the conditions for consulting each routed skill. Do not assume that another skill is hidden from model discovery.
