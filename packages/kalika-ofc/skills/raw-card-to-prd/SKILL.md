---
name: raw-card-to-prd
description: >
  Use this skill when the user provides any of the following:
   - a product card
   - a Jira ticket
   - a Linear issue
   - a vague feature request
   - an A/B test request
   - a business requirement
   - a stakeholder request
   - a PM request
   - a low-quality LLM-generated task
   - a task with unclear business context, scope, actors, or acceptance criteria
   Use this skill before:
   - technical discovery
   - technical planning
   - validation strategy
   - implementation
   - code generation
   - API design
   - database modeling
   - test planning
---

# raw-card-to-prd

## Purpose

Transform a vague, incomplete, ambiguous, or low-quality product card into a clear Product Requirements Document using a human-in-the-loop refinement process.

This skill exists to prevent technical discovery, planning, implementation, or code generation from starting before the product intent is clear enough to support a safe delivery.

The skill must clarify business intent, actors, scope, behavior, business rules, risks, and acceptance criteria before generating the final PRD.

## Core Principle

Do not generate the PRD until all blocking questions have been answered by a human or explicitly accepted as assumptions by the human.

No clarification, no PRD.
No approved PRD, no technical planning.
No technical planning, no implementation.

## When Not to Use

Do not use this skill when:

- the task is a small technical fix with clear reproduction steps
- the task already has a complete and approved PRD
- the task is purely mechanical and has no product ambiguity
- the user explicitly asks to skip product refinement
- the task is a refactor with clear technical scope and no business behavior change

If business impact is unclear, use this skill anyway.

## Required Files

This skill expects the following files to exist:

```txt
raw-card-to-prd/SKILL.md
raw-card-to-prd/templates/PRD.md
raw-card-to-prd/templates/clarification-required.md
```

## Operating Modes

This skill has two modes:

1. `CLARIFICATION_MODE`
2. `PRD_GENERATION_MODE`

The agent must always start in `CLARIFICATION_MODE`.

The agent may switch to `PRD_GENERATION_MODE` only when there are no remaining blocking questions.

## Determinism Rules

The agent must:

- follow the operating modes in order
- classify every uncertainty explicitly
- ask objective questions with stable IDs
- keep the human response format explicit
- avoid mixing clarification output with PRD output
- avoid generating partial PRDs while blocking questions remain
- use templates exactly as output contracts

The agent must not:

- invent business context
- silently assume critical behavior
- convert assumptions into facts
- hide uncertainty
- generate implementation details
- generate architecture decisions
- generate code
- create database schemas
- define API contracts
- create a technical implementation plan
- produce generic product fluff

## Context Management Rules

The agent must keep context small and useful.

### Allowed Context

Use only:

- the raw card
- user-provided comments
- links or references explicitly provided by the user
- human clarification answers
- accepted assumptions
- artifacts generated in the current session

### Forbidden Context Behavior

The agent must not:

- repeat the full raw card in every response
- duplicate large sections unnecessarily
- include irrelevant implementation speculation
- load unrelated repository context unless requested
- use hidden assumptions as facts
- produce long generic explanations instead of actionable artifacts

## Question Labels

Use only these labels for refinement:

- `BLOCKING_QUESTION`
- `NON_BLOCKING_QUESTION`
- `ASSUMPTION_CANDIDATE`

## BLOCKING_QUESTION

Use `BLOCKING_QUESTION` when the answer changes the meaning, scope, behavior, acceptance criteria, or delivery boundary of the feature.

Examples:

- What business problem is being solved?
- Who is the primary actor?
- Is this replacing an existing flow or running in parallel?
- Who is eligible for the new behavior?
- Which current rules must be preserved?
- What is explicitly in scope?
- What is explicitly out of scope?
- What happens in fallback or failure cases?
- What must be true for the task to be considered complete?

If any `BLOCKING_QUESTION` remains unresolved, do not generate the PRD.

## NON_BLOCKING_QUESTION

Use `NON_BLOCKING_QUESTION` when the answer is useful but not required to produce the first reliable PRD.

Examples:

- final UI copy
- secondary metrics
- exact experiment name
- minor wording decisions
- later rollout details
- final analytics naming

Non-blocking questions may appear in the PRD as open items.

## ASSUMPTION_CANDIDATE

Use `ASSUMPTION_CANDIDATE` when the agent can continue safely if the human explicitly accepts the assumption.

Every assumption candidate must include:

- the assumption
- the risk if wrong
- whether it is reversible

Example:

```txt
ASSUMPTION_CANDIDATE:
The new flow will run only for users included in the A/B test group.

Risk:
If wrong, rollout, scope, analytics, and technical impact may change.

Reversible:
Yes
```

Assumptions must not be used unless accepted by the human.

## Clarification Policy

The agent must use a human-in-the-loop clarification process.

The agent must:

- analyze the raw card
- detect ambiguities
- classify uncertainty
- ask only questions required to produce a reliable PRD
- ask questions in batches
- wait for human answers before generating the PRD
- re-analyze after each human response
- ask another clarification round if blockers remain
- if question can be answered by code base analysis, scan the code base and answer it without asking the human

The agent must not ask cosmetic questions during refinement unless they block the PRD.

## Clarification Round Rules

Each clarification round must ask no more than 7 blocking questions.

Prioritize questions in this order:

1. Business problem
2. Affected actors
3. Current flow
4. Proposed behavior
5. Scope boundaries
6. Business rules
7. Acceptance criteria
8. Risks and edge cases

If more than 7 blocking questions exist, ask the most important 7 first and continue in another round after receiving answers.

## CLARIFICATION_MODE

### Goal

Clarify the business request until the PRD can be generated safely.

### Input

The agent receives:

- raw card
- user comments
- optional references
- optional previous clarification answers

### Process

The agent must:

1. summarize current understanding
2. extract explicit requirements
3. identify implied requirements
4. identify missing context
5. identify ambiguous terms
6. identify scope uncertainty
7. identify missing business rules
8. identify missing acceptance criteria
9. produce a clarification request using the template

### Output

When clarification is needed, output only:

```txt
raw-card-to-prd/templates/clarification-required.md
```

Do not output a PRD in the same response.

## Human Response Handling

After receiving human answers, the agent must:

1. map each answer to its question ID
2. update the current understanding
3. identify contradictions
4. identify unanswered blocking questions
5. identify rejected assumptions
6. ask another clarification round if blockers remain
7. switch to `PRD_GENERATION_MODE` only when blockers are resolved

If the human gives partial answers, ask only for the missing blocking information.

If the human rejects an assumption candidate, ask for the correct behavior unless the answer already provides it.

If the human accepts an assumption candidate, record it in the PRD under `Assumptions`.

## PRD_GENERATION_MODE

### Entry Conditions

Enter `PRD_GENERATION_MODE` only when:

- all blocking questions are answered
- all critical assumptions are accepted or replaced with facts
- the business problem is clear
- the primary actors are clear
- the current flow is understood or explicitly marked as not applicable
- the proposed behavior is clear
- scope boundaries are clear
- business rules can be listed
- acceptance criteria can be generated
- no unresolved contradiction remains

### Goal

Generate a complete Product Requirements Document from:

- raw card
- human answers
- accepted assumptions
- clarified business rules
- clarified scope
- clarified acceptance criteria

### Output

When PRD generation is ready, output only:

```txt
raw-card-to-prd/templates/PRD.md
```

Do not include a clarification request in the same response.

## PRD Generation Rules

The PRD must:

- be product-focused
- avoid technical implementation design
- separate facts from assumptions
- preserve unknown non-blocking items as `UNKNOWN`
- include testable acceptance criteria
- include edge cases
- include risks
- include open non-blocking questions
- include decision log
- be ready for technical discovery

The PRD must not:

- write code
- design database schemas
- define API contracts
- choose frameworks
- create implementation plans
- hide uncertainty
- treat assumptions as confirmed facts
- use vague acceptance criteria like "should work" or "must be intuitive"

## PRD Readiness Rules

The PRD can be marked as `APPROVED_FOR_TECHNICAL_DISCOVERY` only if:

- the problem is clear
- the primary actor is identified
- the current flow is described or explicitly marked as not applicable
- the proposed flow is described
- in-scope items are documented
- out-of-scope items are documented
- functional requirements are listed
- business rules are testable
- acceptance criteria are observable
- blocking questions are resolved
- accepted assumptions are documented
- non-blocking questions are listed separately

If any readiness rule fails, return to `CLARIFICATION_MODE`.

## Failure Conditions

Reject PRD generation and continue clarification if:

- the problem is vague
- the target actor is unknown
- the requested behavior has multiple unresolved interpretations
- business rules conflict
- scope boundaries are unclear
- acceptance criteria cannot be produced
- critical assumptions are not accepted
- the card describes only a solution and the actual problem remains unknown
- the feature changes production behavior but rollout is unknown
- eligibility rules are unknown and affect behavior
- fallback behavior is required but undefined

## Anti-Patterns

Reject or rewrite outputs that contain:

- generic product fluff
- invented business metrics
- fake certainty
- architecture decisions before discovery
- code suggestions
- vague acceptance criteria
- requirements without source or confidence
- duplicated context dumps
- hidden assumptions
- "implementation-first" thinking

## Final Response Contract

If clarification is still needed, output a clarification request using:

```txt
raw-card-to-prd/templates/clarification-required.md
```

If the PRD is ready, output the PRD using:

```txt
raw-card-to-prd/templates/PRD.md
```

Never return both at the same time.

Never generate a partial PRD while blocking questions remain.
