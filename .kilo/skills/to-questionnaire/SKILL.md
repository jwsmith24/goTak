---
name: to-questionnaire
description: Create a discovery questionnaire for another person to answer. Use when the user explicitly asks to turn unresolved decisions or knowledge gaps into a Markdown questionnaire.
---

# To questionnaire

Turn something the user cannot answer alone into a **questionnaire**: a Markdown document they hand to one person to fill in asynchronously or complete together in a meeting. The recipient holds knowledge the user lacks; the questionnaire draws it out.

**Grill the send, not the subject.** Interview the user only about the _send_, which they can answer: who it goes to and what they need back. The questions in the document then target the **gap** between what the recipient knows and what the user needs.

1. **Who is it going to?** Ask, in one exchange, the recipient's role, expertise, and relationship to the user. This fixes the questionnaire's tone and how much context it must carry. Done when you know who the recipient is and what they know that the user does not.

2. **What do you need back?** Ask, in one exchange, which decisions or facts the user cannot resolve alone and needs from this person. Done when you have a concrete list of what the user must walk away able to do or decide.

3. **Write the questionnaire.** Draft questions aimed at the gap from steps 1-2, following the document structure below. Write it to `to-questionnaire-<slug>.md` in the current directory, deriving the slug from the topic, and report the absolute path. Done when the file exists and every item the user named in step 2 is covered by a question.

## Document structure

Frame the document as a **discovery questionnaire**: the user lacks context and the recipient holds it. Order questions most-important-first because asynchronous responses may provide only one pass. Once there are more than a handful, group questions under `##` headings by theme. Use this template:

<questionnaire-template>

# <Questionnaire title>

**Purpose:** why this questionnaire exists and the decision riding on it.

**From:** <the user>, **To:** <the recipient>, **How your answers will be used:** <where they go>

## Context

One paragraph orienting a recipient who was not in the user's head. Include enough to answer well, not a page.

## How to answer

State the deadline and rough effort. Partial answers and "I don't know" are useful: ask the recipient to flag uncertainty rather than skip the question.

## <Theme heading>

Use one `##` section per theme. Under each, put its questions most-important-first. Each question covers one idea, never a compound question, with an answer stub directly beneath it. Add a one-line _why this matters_ only where the question could be misread or invite a throwaway answer.

<question-example>
### What load is the system expected to handle at launch?

_Why this matters: it decides whether we provision for burst traffic now or defer it._

>
</question-example>

## Anything else?

Close with a catch-all: anything we did not ask that we should know?

</questionnaire-template>
