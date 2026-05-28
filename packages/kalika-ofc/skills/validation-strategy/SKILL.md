---
name: validation-strategy
description: Use this skill when a plan, implementation, or risky change needs a concrete validation strategy. Call it before implementation, before review acceptance, or when the workflow must map acceptance criteria to tests, commands, manual checks, and pass/fail rules.
---

# Validation Strategy Skill

## Purpose

Create a concrete validation strategy mapped to acceptance criteria, project commands, risk areas, and regression coverage.

This skill exists because “run the tests” is not a strategy. It is a sentence people say before discovering the test suite checks basically nothing.

## Trigger

Call this skill when:

- A plan is ready but implementation has not started.
- A workflow needs test strategy before coding.
- A change affects security, performance, data, API contracts, or legacy behavior.
- A validator must actively verify an implementation.
- The orchestrator needs pass/fail criteria before accepting work.
- Review or validation found coverage gaps.

## Do not use when

Do not call this skill when:

- There is no plan or change scope to validate.
- The task is only initial routing.
- A valid validation strategy already exists and the scope did not change.
- The task is documentation-only and has no behavioral validation requirement.

## Required inputs

The agent must provide:

- Approved plan or implementation summary.
- Acceptance criteria.
- Changed or affected files/modules.
- Known commands from the project.
- Risk areas.
- Existing tests, if known.
- Manual verification needs, if any.

## Validation dimensions

Consider:

- Unit tests.
- Integration tests.
- Contract tests.
- End-to-end tests.
- Static analysis.
- Type checks.
- Lint.
- Build.
- Security checks.
- Performance checks.
- Manual verification steps.

## Rules

- Use project-native commands.
- Do not invent scripts unless explicitly asked.
- Mark manual checks explicitly.
- Include coverage gaps.
- Map every acceptance criterion to at least one validation method or explicitly mark it as unvalidated.
- Prefer deterministic validation over subjective inspection.

## Expected output

Return this exact structure:

```md
# Validation Strategy

## Acceptance criteria mapping
| Criteria | Validation method | Required | Evidence expected |
|---|---|---:|---|

## Commands
| Command | Purpose | Evidence source |
|---|---|---|

## Regression areas
<list>

## Test data
<needed test data, or none>

## Manual checks
<list, or none>

## Coverage gaps
<what cannot be validated automatically>

## Pass/fail criteria
<clear rules>

## Blocking validation requirements
<items that must pass before acceptance>
```

## Stop conditions

Stop and return `BLOCKED` when:

- No reliable validation method exists for a critical requirement.
- Required project commands are unknown and cannot be inferred from files.
- The implementation scope is too vague to validate.
- Security/data risks lack validation coverage.
