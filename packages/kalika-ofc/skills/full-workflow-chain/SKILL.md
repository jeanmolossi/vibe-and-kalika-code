---
name: full-workflow-chain
description: Use this skill when an agent must build or execute a multi-workflow pipeline such as spike -> refinement -> validation -> development -> review -> validation -> documentation -> learning. Call it before delegating the first workflow in an end-to-end task.
---

# Full Workflow Chain Skill

## Purpose

Build a deterministic multi-workflow pipeline with explicit gates, required artifacts, and stop conditions.

This skill prevents full-workflow execution from becoming uncontrolled agent chaos.

## Trigger

Call this skill when:

- The user asks for full work, complete execution, or end-to-end delivery.
- The task requires multiple workflows.
- The parent wants a chain such as spike -> refinement -> validation -> development -> review -> validation -> documentation -> learning.
- Existing artifacts may allow some phases to be skipped.
- The orchestrator must coordinate workflow agents in sequence.

## Do not use when

Do not call this skill when:

- A single workflow is enough.
- The task is only research, review, validation, or documentation.
- There is no parent orchestrator controlling the pipeline.
- The next phase is already explicitly selected and isolated.

## Required inputs

The agent must provide:

- User request.
- Desired pipeline, if provided.
- Existing artifacts.
- Constraints.
- Maximum retry count.
- Base branch.
- Validation requirements.
- Documentation requirements.
- Learning policy.

## Default chain

Use this chain unless the parent provides another:

```txt
workflow-spike
-> workflow-refinement
-> workflow-validation(strategy)
-> workflow-development
-> workflow-review
-> workflow-validation(active)
-> workflow-documentation
-> workflow-learning
```

## Skip policy

A phase may be skipped only if:

- An equivalent approved artifact already exists.
- The parent explicitly marks it as unnecessary.
- The phase is not applicable and the reason is documented.

Never skip:

- review after development
- active validation after review
- phase gates between workflows

## Expected output

Return:

```md
# Full Workflow Chain Plan

## Pipeline
| Order | Workflow | Mode | Required input | Required output | Gate |
|---|---|---|---|---|---|

## Skippable phases
| Workflow | Reason |
|---|---|

## Non-skippable phases
- 

## Retry policy
- 

## Context policy
- 

## Stop conditions
- 

## Next handoff
<first workflow handoff>
```

## Stop conditions

Stop and return `BLOCKED` when:

- Required input for the first phase is missing.
- The requested pipeline violates safety or repository constraints.
- The parent asks to skip mandatory review/validation without explicit acceptance of risk.
- The task requires implementation but no validation strategy can be created.
