# Design it twice

Use this process when meaningful interface alternatives exist.

1. Frame constraints, dependency categories, domain vocabulary, and the behavior behind the seam. A small code sketch may clarify constraints but must not preselect a design.
2. Produce at least three radically different interfaces. When native read-only subagents are available, ask them in parallel for independent designs; otherwise perform separate direct passes. Useful constraints are minimum surface, maximum extension, easiest common case, and ports-and-adapters placement.
3. For each design show the complete interface contract, caller example, hidden behavior, dependency strategy, and trade-offs.
4. Compare depth, locality, seam placement, migration cost, and test surface. Recommend one design or a coherent hybrid.

Subagent output is research input. Verify it and synthesize in the parent task.