---
name: codebase-design
description: Design or improve deep modules and their interfaces, seams, adapters, test surfaces, and dependency placement. Use for module design, testability, AI navigability, deepening shallow code, or comparing alternative interfaces.
---

# Codebase design

Design **deep modules**: substantial behavior behind a small interface at a clean seam, tested through that interface.

## Vocabulary

- **Module:** anything with an interface and implementation, at any scale.
- **Interface:** everything a caller must know: types, invariants, ordering, errors, configuration, and performance characteristics.
- **Implementation:** behavior hidden inside a module.
- **Depth:** leverage delivered per unit of interface a caller must learn.
- **Seam:** a location where behavior can vary without editing the caller.
- **Adapter:** a concrete implementation satisfying an interface at a seam.
- **Leverage:** capability callers gain from the interface.
- **Locality:** concentration of change, knowledge, bugs, and verification.

Prefer this vocabulary over ambiguous substitutes when module shape is the topic.

## Design tests

- Reduce methods and parameters while hiding coherent complexity.
- Apply the deletion test: deleting a useful deep module redistributes complexity across callers; deleting a pass-through removes complexity.
- Treat the interface as the caller and test surface. A need to test past it signals a weak shape.
- Accept dependencies and return observable results where the domain permits.
- Introduce a seam when real alternatives justify it, commonly production and test adapters, not for hypothetical variation.
- Keep test-only internal seams private to the implementation.

For dependency-specific deepening, read [DEEPENING.md](DEEPENING.md). For competing interface designs, read [DESIGN-IT-TWICE.md](DESIGN-IT-TWICE.md).

Design analysis is not authorization to edit code. Preview proposed file changes and obtain explicit approval before implementation; git and external actions remain separate.