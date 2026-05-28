---
name: prd-to-implementation-plan
description: Use this skill when you have a finalized PRD, usually produced by raw-card-to-prd, and need to transform it into a deterministic, executable implementation plan with clear phases, tasks, validation gates, risks, dependencies, and human clarification only when blockers remain.
---

# prd-to-implementation-plan

## Purpose

Transform a finalized Product Requirements Document into a deterministic implementation plan that can be executed by an implementer, validator, reviewer, or multi-agent workflow without guessing what the product team meant.

This skill does not implement code. It converts product intent into an execution-ready technical plan.

## Trigger

Use this skill when:

- A PRD exists and needs to become an implementation plan.
- The PRD was produced by `raw-card-to-prd`.
- A user asks for a delivery plan, implementation plan, technical plan, work breakdown, execution plan, or agent-ready plan based on a PRD.
- A coordinator needs a canonical `plan.md` before development starts.

## Do not use this skill when

- There is no PRD or equivalent product specification.
- The user is still trying to clarify the product requirements.
- The task is to write code directly.
- The task is only to review an existing implementation.
- The task is only to estimate time or staffing without implementation details.

## Hard gates

### Gate 1: No PRD, no plan

If no PRD is provided or discoverable, stop and ask for the PRD path/content.

Do not invent requirements.

### Gate 2: No unresolved blocking ambiguity

If a requirement affects architecture, data model, API contract, rollout, security, compliance, user flow, or acceptance criteria and remains unclear, enter `CLARIFICATION_MODE`.

Do not generate the final implementation plan until all blocking questions have either:

- been answered by a human, or
- been explicitly accepted by a human as assumptions.

### Gate 3: No implementation before validation strategy

The implementation plan must define validation before implementation tasks are considered ready.

Every implementation slice must include:

- acceptance criteria,
- test strategy,
- verification command or manual validation path,
- rollback or mitigation note.

## Operating modes

The skill has four modes. Follow them in order.

1. `PRD_INTAKE_MODE`
2. `PRD_ANALYSIS_MODE`
3. `CLARIFICATION_MODE`
4. `PLAN_GENERATION_MODE`

Never mix outputs from `CLARIFICATION_MODE` with the final implementation plan.

---

## Mode 1: PRD_INTAKE_MODE

### Goal

Collect only the context required to plan implementation.

### Allowed inputs

- The PRD content.
- The PRD file path.
- Human clarification answers.
- Accepted assumptions.
- Explicit technical constraints provided by the user.
- Existing architecture docs, business rules docs, API references, ADRs, or repository maps explicitly provided or requested.
- Prior artifacts from the same session.

### Forbidden context behavior

Do not:

- load unrelated repository context,
- repeat the full PRD in every response,
- summarize away acceptance criteria,
- use unstated product behavior as fact,
- silently convert uncertainty into requirements,
- broaden scope beyond the PRD.

### Output

A compact context summary:

```md
## Context Summary

PRD source: <path-or-inline>
Planning target: <feature/system/change>
Known constraints:
- <constraint>

Available artifacts:
- <artifact>

Missing artifacts:
- <artifact-or-none>
```

---

## Mode 2: PRD_ANALYSIS_MODE

### Goal

Extract implementation-relevant facts from the PRD and classify uncertainty.

### Required extraction

Produce stable IDs using the following format:

- `REQ-001` for product requirements.
- `AC-001` for acceptance criteria.
- `NFR-001` for non-functional requirements.
- `FLOW-001` for user/system flows.
- `DATA-001` for data model requirements.
- `API-001` for API/integration requirements.
- `SEC-001` for security/compliance requirements.
- `OBS-001` for observability requirements.
- `RISK-001` for risks.
- `UNK-001` for unknowns.
- `BQ-001` for blocking questions.
- `ASM-001` for assumption candidates.

### Classification rules

Every uncertainty must be classified as one of:

- `BLOCKING_QUESTION`: must be answered before planning.
- `ASSUMPTION_CANDIDATE`: can be accepted by the human if planning needs to continue.
- `NON_BLOCKING_UNKNOWN`: can be handled inside a bounded discovery task.

### Blocking ambiguity examples

Treat as blocking when unclear:

- user flow changes,
- API contract,
- persistence model,
- authorization or permission behavior,
- external integration behavior,
- rollout strategy for risky changes,
- feature flag or A/B test mechanics,
- definition of success,
- acceptance criteria,
- production data migration behavior,
- security/compliance constraints.

### Non-blocking ambiguity examples

Treat as non-blocking when it can be resolved by a bounded task, such as:

- exact file paths to inspect,
- current component location,
- current test command,
- current lint command,
- local naming conventions.

---

## Mode 3: CLARIFICATION_MODE

### Goal

Ask the smallest possible set of questions needed to remove blockers.

Use `templates/clarification-required.md` exactly.

### Rules

- Ask only blocking questions first.
- Use stable IDs: `BQ-001`, `BQ-002`.
- Do not ask vague questions.
- Do not ask questions that can be answered by inspecting provided artifacts.
- Provide assumption candidates only when they are safe and reversible.
- Require the human to answer using the exact response format.

### Stop condition

Exit `CLARIFICATION_MODE` only when every blocking question has either:

- a direct answer, or
- an explicit accepted assumption.

---

## Mode 4: PLAN_GENERATION_MODE

### Goal

Generate an execution-ready implementation plan.

Use `templates/implementation-plan.md` exactly.

### Planning principles

The plan must be:

- deterministic,
- ordered,
- executable,
- testable,
- scoped to the PRD,
- safe for handoff to another agent or developer,
- clear enough for a low-context implementer to follow without inventing behavior.

### Required plan sections

The implementation plan must include:

1. Plan metadata.
2. Source PRD reference.
3. Context summary.
4. Requirement mapping.
5. Assumptions accepted by human.
6. Out-of-scope items.
7. Architecture impact.
8. Data model impact.
9. API/integration impact.
10. Frontend/UI impact.
11. Backend/domain impact.
12. Security and compliance notes.
13. Observability notes.
14. Rollout strategy.
15. Test strategy.
16. Implementation phases.
17. Atomic task breakdown.
18. Validation gates.
19. Risks and mitigations.
20. Handoff contract.
21. Final readiness check.

### Implementation phases

Each phase must have:

```md
## Phase <N>: <name>

ID: PHASE-<NNN>
Goal: <one sentence>
Depends on: <phase IDs or none>
Inputs:
- <artifact or requirement ID>
Outputs:
- <artifact or code outcome>
Tasks:
- TASK-<NNN>: <atomic task>
Validation:
- <command or manual check>
Done when:
- <binary condition>
Rollback/mitigation:
- <note>
```

### Task rules

Each task must be atomic.

A task is atomic when:

- it has one clear objective,
- it has one owner role,
- it has explicit inputs,
- it has explicit outputs,
- it can be validated,
- it does not require hidden product decisions.

Avoid tasks like:

- "implement backend",
- "adjust frontend",
- "fix tests",
- "review everything",
- "improve performance",
- "handle edge cases".

Use tasks like:

- "Add request DTO for `<endpoint>` containing fields `<fields>` and validation rules `<rules>`."
- "Create unit tests for `<domain function>` covering `<case A>`, `<case B>`, `<case C>`."
- "Add feature flag guard around `<flow>` and verify disabled behavior remains unchanged."

### Spike task rules

A discovery/spike task is allowed only when:

- it is bounded,
- it has a concrete output artifact,
- it does not replace product clarification,
- it cannot be answered from the PRD alone.

Example:

```md
TASK-003: Inspect current storage flow for wizard data.
Owner role: researcher
Input: repository files related to localStorage persistence
Output: `.ai/sessions/<session>/research/storage-flow.md`
Done when: artifact lists files, flow summary, constraints, and risks.
```

---

## Output artifacts

Recommended session paths:

```txt
.ai/sessions/<session-id>/prd.md
.ai/sessions/<session-id>/prd-analysis.md
.ai/sessions/<session-id>/clarification-required.md
.ai/sessions/<session-id>/clarification-answers.md
.ai/sessions/<session-id>/plan.md
.ai/sessions/<session-id>/planning-state.md
```

Session ID format:

```txt
YYYY-MM-DD--HH-mm_<short-headline>
```

Example:

```txt
2026-05-11--14-30_job-title-ai-ab-test
```

---

## Determinism rules

- Use stable IDs.
- Keep original PRD requirement wording in the requirement map when possible.
- Preserve product intent.
- Do not rename requirements unless adding a stable ID.
- Do not merge unrelated requirements.
- Do not split one requirement without preserving traceability.
- Every task must trace back to at least one requirement, acceptance criterion, risk, or non-functional requirement.
- Every assumption must be explicitly marked as human-accepted before it influences the final plan.
- Every risk must have a mitigation.
- Every validation gate must have pass/fail criteria.

---

## Context minimization rules

When handing the plan to another agent, pass only:

- the final implementation plan,
- the source PRD path,
- accepted assumptions,
- relevant referenced artifacts,
- the exact phase/task IDs assigned to that agent.

Do not pass:

- the full chat history,
- unrelated PRD sections,
- raw brainstorming,
- rejected assumptions,
- unrelated repository dumps.

---

## Quality gate

Before finalizing the plan, run `checklists/plan-quality-gate.md`.

If using a local environment, optionally run:

```bash
./scripts/validate-plan.sh .ai/sessions/<session-id>/plan.md
```

If the plan fails the quality gate, do not present it as final.

---

## Final response format

When the plan is complete, respond with:

```md
Implementation plan generated.

Artifacts:
- PRD analysis: <path>
- Implementation plan: <path>

Readiness: READY | BLOCKED
Blocking items: <none-or-list>
Next recommended step: <one sentence>
```

When blocked, respond with:

```md
Implementation planning blocked.

Reason: unresolved product or technical ambiguity.
Required human answers:
- BQ-001: <question>
- BQ-002: <question>

Use the clarification template before generating the plan.
```
