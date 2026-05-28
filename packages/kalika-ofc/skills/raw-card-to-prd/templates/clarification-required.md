# Clarification Required

## Summary of Current Understanding

<!-- Summarize what is currently understood from the raw card and previous human answers. Do not repeat the full raw card. -->

...

## Explicit Requirements Identified

| ID | Requirement | Source |
|---|---|---|
| ER-001 | ... | Raw card / Human answer |

## Implied Requirements Identified

| ID | Requirement | Reason | Confidence |
|---|---|---|---|
| IR-001 | ... | ... | Low / Medium / High |

## Ambiguities Detected

| ID | Ambiguity | Why it matters |
|---|---|---|
| AMB-001 | ... | ... |

## Blocking Questions

| ID | Question | Why it blocks the PRD |
|---|---|---|
| BQ-001 | ... | ... |

## Assumption Candidates

| ID | Assumption | Risk if wrong | Reversible |
|---|---|---|---|
| ASM-001 | ... | ... | Yes / No |

## Non-Blocking Questions

| ID | Question | Why it matters | Can proceed without answer |
|---|---|---|---|
| NBQ-001 | ... | ... | Yes |

## Required Human Response Format

Please answer the blocking questions using this format:

```txt
BQ-001: ...
BQ-002: ...
BQ-003: ...
```

For assumption candidates, answer using this format:

```txt
ASM-001: accepted/rejected/adjusted - ...
ASM-002: accepted/rejected/adjusted - ...
```

## Next Step

After receiving the human answers, the agent must re-analyze the remaining uncertainty.

If blocking questions remain, the agent must generate another clarification request.

If no blocking questions remain, the agent must generate the PRD using `raw-card-to-prd/templates/PRD.md`.
