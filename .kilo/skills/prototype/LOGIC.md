# Logic prototype

Build one self-contained HTML file that a non-developer can open locally and use to drive the model.

1. Put the question in a visible introduction.
2. Isolate the behavior in a pure reducer, explicit state machine, pure function set, or small state-owning module. Keep DOM access in the page shell.
3. Render the complete relevant state in domain language after every action.
4. Provide free-play actions and deterministic guided scenarios for the happy path, a difficult edge case, and an illegal action.
5. Keep all HTML, CSS, and JavaScript inline; avoid network access and external packages.
6. Hand over the absolute path and instructions. Opening the file or a browser requires explicit authorization.

The shell is disposable. Any validated logic must still be reviewed and implemented to the production codebase's standards.