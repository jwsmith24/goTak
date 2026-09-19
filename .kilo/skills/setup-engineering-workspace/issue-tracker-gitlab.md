# Issue tracker: GitLab

Issues live in GitLab. Infer the project from the configured remote and use the installed `glab` CLI for reads and authorized mutations.

## Operations

- Read/list: `glab issue view <iid> --comments`, `glab issue list -F json`
- Create: `glab issue create --title ... --description ...`
- Comment: `glab issue note <iid> --message ...`
- Labels: `glab issue update <iid> --label ... --unlabel ...`
- Close: post the approved explanation, then `glab issue close <iid>`

Merge requests are not a triage request surface unless this file is explicitly changed to say they are.

## Authorization

Reads may run without mutation approval. Before any create, comment, label, assignment, edit, close, or other remote mutation, preview the exact targets and content and obtain explicit authorization for that batch. Closing and publication remain consequential even when earlier edits were approved. Never infer authorization from a command example.

## Publication language

“Publish to the tracker” means create a GitLab issue only after approval. “Fetch the ticket” means read the issue and relevant comments.