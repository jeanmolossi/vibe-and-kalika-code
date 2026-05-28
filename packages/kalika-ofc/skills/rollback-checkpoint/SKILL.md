---
name: rollback-checkpoint
description: Use this skill when planning or executing large, risky, behavior-changing, security-sensitive, migration-heavy, or legacy refactor work. Call it before implementation to define reversible checkpoints, validation points, and rollback actions.
---

# Rollback Checkpoint Skill

## Purpose

Force large or risky changes into smaller reversible phases with validation and rollback plans.

This skill exists because “big bang refactor” is just production roulette wearing a hoodie.

## Trigger

Call this skill when:

- Many files will change.
- Public APIs or contracts may change.
- Data migrations are involved.
- Authentication or authorization behavior changes.
- Performance-critical paths change.
- Legacy refactor affects core behavior.
- Rollback would be hard after full implementation.
- The plan mixes refactor and behavior changes.

## Do not use when

Do not call this skill when:

- The change is small, isolated, and easily reverted.
- No behavior, schema, interface, config, or security-sensitive path changes.
- The workflow is only research or documentation with no change execution.
- A validated checkpoint plan already exists.

## Required inputs

The agent must provide:

- Approved plan.
- Files/modules affected.
- Risk areas.
- Validation strategy.
- Deployment or runtime constraints, if known.
- Current rollback options.
- Any migration or data safety notes.

## Checkpoint rules

- Prefer small reversible phases.
- Avoid mixing refactor and behavior changes.
- Validate before moving to the next phase.
- Make rollback explicit for each phase.
- Keep checkpoints aligned with acceptance criteria.
- Do not hide irreversible risk.

## Expected output

Return this exact structure:

```md
# Rollback Checkpoint Plan

## Risk level
low | medium | high | critical

## Reason
<why checkpointing is or is not required>

## Checkpoints

### Checkpoint 1: <name>

#### Expected state
<what should be true>

#### Validation
<how to verify>

#### Rollback
<how to revert safely>

#### Risk
<remaining risk>

## Non-reversible changes
<list, or none>

## Required approval before implementation
yes | no
```

## Stop conditions

Stop and return `BLOCKED` when:

- A risky change has no rollback path.
- Validation is insufficient for the risk level.
- Data loss or security risk is possible and unaddressed.
- The plan is too large to checkpoint safely.
