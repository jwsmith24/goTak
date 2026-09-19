# Offline HTML report

Produce one static HTML file with inline CSS and inline SVG/HTML diagrams. Do not load scripts, fonts, styles, Mermaid, analytics, or other resources from a network.

Include a compact legend and one card per candidate:

- title, recommendation strength, and dependency category;
- involved modules/files;
- side-by-side before/after schematic;
- one-sentence problem and solution;
- short leverage, locality, and test-surface gains;
- a visible ADR conflict when applicable.

End with one top recommendation and rationale. Escape repository-derived text before embedding it in HTML. Keep diagrams understandable without JavaScript or network access.