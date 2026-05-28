---
name: knowledge-extraction
description: Use this skill after a session completes and only when durable reusable knowledge should be saved. Call it after the final report to extract project conventions, recurring bugs, validated commands, architecture decisions, or reusable workflow lessons without loading or saving memory dumps.
---

# Knowledge Extraction Skill

## Purpose

Extract durable reusable lessons after a completed session while avoiding memory bloat, secrets, stale task details, and useless trivia.

This skill exists because memory should be a sharp tool, not a landfill with timestamps.

## Trigger

Call this skill when:

- A workflow or orchestration session is complete.
- The final report contains reusable learning.
- A new project convention was discovered.
- A recurring bug or failure pattern was confirmed.
- A command or validation method was verified.
- An architecture decision or module boundary was confirmed.
- A user preference relevant to future work was explicitly stated.

## Do not use when

Do not call this skill when:

- The session is still in progress.
- The final report is missing.
- The finding is speculative.
- The detail is temporary or one-off.
- The information contains secrets, credentials, private data, raw logs, or full diffs.
- The lesson is already documented in memory or project docs.

## Required inputs

The agent must provide:

- Final report.
- Relevant accepted artifacts.
- Source paths for evidence.
- Existing memory index or known memory references, if available.
- Candidate lessons.

## Save only durable knowledge

Good memory candidates:

- Project conventions.
- Recurring bugs.
- Architecture decisions.
- Validated commands.
- Reusable workflow lessons.
- Known testing strategy.
- Stable file ownership or module boundaries.
- Confirmed constraints that affect future implementation.

Do not save:

- Secrets.
- Credentials.
- Private data.
- Raw logs.
- Temporary task details.
- Speculation.
- Full diffs.
- Full files.
- One-off failures with no reusable value.
- Anything copied blindly from the conversation.

## Extraction process

1. Read the final report.
2. Identify candidate lessons.
3. Verify each lesson is durable and reusable.
4. Check whether it already exists in memory or docs.
5. Save short memory entries only.
6. Include source artifact path.
7. Do not load, rewrite, or summarize full memory dumps.

## Expected output

Return this exact structure:

```md
# Knowledge Extraction Report

## Decision
SAVE | SKIP

## Reason
<why saving is or is not justified>

## Candidate memories
| Title | Type | Why durable | Source |
|---|---|---|---|

## Memory entries to save

### <short title>

Memory entry:

~~~md
---
type: project | feedback | reference | user
source: <artifact path>
created_at: <YYYY-MM-DD>
---

# <short title>

## Summary
<1-3 sentences>

## When to use
<when future agents should apply this>

## Evidence
<artifact/file path>
~~~

## Skipped items
<items intentionally not saved and why>
```

## Stop conditions

Stop and return `SKIP` when:

- No durable lesson exists.
- Evidence is weak.
- The information is sensitive or temporary.
- Saving would duplicate existing memory.
