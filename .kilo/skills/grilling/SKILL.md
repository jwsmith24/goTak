---
name: grilling
description: Stress-test a plan, decision, design, or idea through a relentless decision-tree interview. Use when the user asks to be grilled, says “grill me,” or wants hidden assumptions and dependent decisions exposed before action.
---

# Grilling

Interview the user until you reach a shared understanding. Map the subject as a **design tree**: every decision branches into the decisions that depend on it.

Work the tree in **rounds**. The **frontier** is every unresolved decision whose prerequisites are settled. Ask the whole current frontier in one round, numbering each question and giving a recommended answer. Then wait for the user's answers before recomputing the frontier. If one answer depends on another unresolved question, put it in a later round.

Format each question like this:

```text
❓ **Q1** - **<question title>**: <question body and choices, when useful>

➡️ <recommended answer and brief rationale>
```

**Facts are the agent's job; decisions are the user's.** Before asking a question that depends on inspectable facts, use the available filesystem, search, command, or web tools to gather them. Run independent lookups together when possible. Ask the user only for facts that are unavailable through those sources, and state what was checked. Put each actual decision to the user.

The interview is complete only when the frontier is empty: every branch has been visited and no assumption remains silent. Summarize the resulting shared understanding and ask the user to confirm it. End there; implementation or other action requires the user's confirmation and authorization.