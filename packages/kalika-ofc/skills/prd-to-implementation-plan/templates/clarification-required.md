# Clarification Required

Status: BLOCKED
Source PRD: <path-or-inline-reference>
Planning target: <feature/system/change>

## Summary

The PRD is not ready for implementation planning because one or more blocking decisions are missing.

## Extracted facts

- REQ-001: <fact from PRD>
- AC-001: <acceptance criterion from PRD>
- NFR-001: <non-functional requirement from PRD>

## Blocking questions

Answer these before the implementation plan is generated.

### BQ-001: <question>

Why this blocks planning: <architecture/data/API/security/rollout/acceptance impact>
Affected items: <REQ/AC/NFR/FLOW/API/DATA IDs>
Required answer shape: <exact expected format>

### BQ-002: <question>

Why this blocks planning: <impact>
Affected items: <IDs>
Required answer shape: <exact expected format>

## Assumption candidates

These are optional. The human must explicitly accept them before they can be used.

### ASM-001: <assumption>

Why it may be acceptable: <reason>
Risk if wrong: <risk>
Rollback/mitigation: <mitigation>
Affected items: <IDs>

## Non-blocking unknowns

These can be handled as bounded discovery tasks inside the plan.

- UNK-001: <unknown>
  - Resolution task: <future TASK-ID placeholder or discovery action>
  - Why non-blocking: <reason>

## Required human response format

```md
# Clarification Answers

## Blocking questions

BQ-001: <answer>
BQ-002: <answer>

## Accepted assumptions

ASM-001: ACCEPTED | REJECTED | REPLACE_WITH: <new assumption>

## Additional constraints

- <constraint or none>
```
