---
name: phase-gate
description: Use this skill when an agent must decide whether a workflow can advance to the next phase. Call it after receiving an artifact from a workflow or specialist agent and before allowing research, planning, validation, implementation, review, or documentation to continue.
---

# Phase Gate Skill

## Purpose

Evaluate whether a produced artifact is good enough for the workflow to advance.

This skill prevents the orchestrator from accepting pretty nonsense just because it was formatted confidently.

## Trigger

Call this skill when:

- A specialist returns an artifact.
- A workflow phase claims to be complete.
- The orchestrator needs to approve or reject progression.
- A plan, research report, validation strategy, implementation summary, review, or documentation artifact is submitted.
- A correction loop needs to decide whether the fix is sufficient.

## Do not use when

Do not call this skill when:

- No artifact has been produced yet.
- The task is still being routed.
- The agent is preparing a handoff contract.
- The user explicitly asked for a draft without validation.

## Required inputs

The agent must provide:

- Artifact name or path.
- Artifact content or concise summary.
- Original handoff contract.
- Acceptance criteria.
- Current workflow phase.
- Known constraints and risks.

## Gate states

- `ACCEPTED`: artifact meets the requirements and the workflow may advance.
- `REVISION_REQUIRED`: artifact is useful but incomplete; targeted correction is needed.
- `BLOCKED`: workflow cannot continue safely due to missing information, missing access, or unresolved risk.
- `REJECTED`: artifact violates scope, quality, or required format.

## Gate checklist

Evaluate:

- Required artifact exists.
- Output format matches the contract.
- Scope was respected.
- Assumptions are explicit.
- Risks and blockers are documented.
- Acceptance criteria are covered.
- Missing inputs are visible.
- Next action is clear.
- No unauthorized pipeline advancement happened.
- No implementation happened in research/planning/review-only phases.

## Expected output

Return this exact structure:

```md
# Phase Gate Decision

## Artifact
<name/path>

## Phase
<current phase>

## Decision
ACCEPTED | REVISION_REQUIRED | BLOCKED | REJECTED

## Reasons
<ordered list>

## Missing requirements
<list, or none>

## Required correction
<targeted correction, or none>

## Next action
<what the parent should do next>
```

## Stop conditions

Stop and escalate to parent when:

- The artifact cannot be evaluated from available context.
- The artifact contradicts the original goal.
- The artifact creates security, data, or production risk.
- Required evidence is absent.
