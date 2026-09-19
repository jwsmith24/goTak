# Deepening modules

Classify dependencies before moving behavior behind a smaller interface:

1. **In-process:** merge coherent pure or in-memory behavior and test directly through the new interface.
2. **Local-substitutable:** use a faithful local implementation such as an in-memory filesystem or local database; keep that seam internal.
3. **Remote but owned:** define a port at the network seam; inject transport adapters while the deep module owns policy.
4. **External:** inject a narrow port around the third-party behavior and test with a controlled adapter.

Replace shallow tests with behavior tests at the new interface only when coverage is equivalent and deletion is included in the approved change. Tests should assert observable outcomes and survive internal refactors.