---
name: arbitration
description: Use this skill when multiple agents return competing outputs or when parallel results must be compared before acceptance. Call it before selecting a plan, implementation, research conclusion, review decision, or documentation version.
---

# Arbitration Skill

## Purpose

Compare competing outputs and select the safest, most correct, most maintainable result.

This skill exists because parallel agents sometimes produce three different answers with the confidence of a senior engineer and the consistency of a broken toaster.

## Trigger

Call this skill when:

- Multiple agents return competing outputs.
- Parallel implementers modify related behavior.
- Two plans propose different technical approaches.
- Research findings disagree.
- Reviewers disagree on whether something is blocking.
- Validation results conflict.
- The orchestrator must choose one candidate before advancing.

## Do not use when

Do not call this skill when:

- There is only one artifact and no competing alternative.
- A phase gate can decide acceptance directly.
- The difference is purely formatting.
- The candidates cannot be compared from available evidence.

## Required inputs

The agent must provide:

- Candidate artifacts or summaries.
- Original goal.
- Acceptance criteria.
- Known constraints.
- Validation evidence, if available.
- Risk and rollback notes.
- Any reviewer or validator findings.

## Compare by

Evaluate each candidate by:

1. Correctness.
2. Acceptance criteria coverage.
3. Security.
4. Performance.
5. Maintainability.
6. Simplicity.
7. Local convention fit.
8. Risk and rollback.
9. Testability.
10. Compatibility with approved plan.

## Expected output

Return this exact structure:

```md
# Arbitration Report

## Candidates
| Candidate | Strengths | Weaknesses | Risk | Evidence |
|---|---|---|---|---|

## Decision
<selected candidate or no winner>

## Reason
<why selected>

## Rejected candidates
<why each candidate was rejected>

## Required merge/correction
<actions before acceptance>

## Gate recommendation
ACCEPTED | REVISION_REQUIRED | BLOCKED | REJECTED
```

## Stop conditions

Stop and return `BLOCKED` when:

- Evidence is insufficient.
- Candidates affect different scopes and cannot be compared fairly.
- All candidates violate acceptance criteria.
- The safest result is to request a targeted correction.
