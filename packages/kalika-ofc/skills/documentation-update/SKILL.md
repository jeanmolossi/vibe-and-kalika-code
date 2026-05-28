---
name: documentation-update
description: Use this skill when an agent must update existing documentation or create new documentation based on an audited task, patch, or workflow change. Call it after documentation-audit.
---

# Documentation Update Skill

## Purpose

Apply concise documentation changes based on verified behavior and the documentation audit.

## Trigger

Call this skill when:

- `documentation-audit` reports `DOCS_NEEDED`.
- Existing documentation is stale.
- New behavior requires new docs.
- A task creates new commands, APIs, workflows, architecture, configuration, or operational procedures.
- A documentation workflow has enough evidence to update docs safely.

## Do not use when

Do not call this skill when:

- Documentation audit has not been performed.
- Behavior is speculative or unverified.
- The task does not require documentation.
- The agent is not allowed to edit docs.
- The required documentation source is unavailable.

## Required inputs

The agent must provide:

- Documentation audit result.
- Task summary.
- Verified behavior.
- Changed file list.
- Validation or review evidence.
- Target docs paths.
- Documentation style/conventions.

## Update rules

- Write only what is supported by evidence.
- Prefer targeted edits.
- Use clear headings.
- Include commands only when validated or explicitly provided.
- Mark assumptions clearly.
- Do not copy large code blocks unless necessary.
- Do not include secrets.
- Do not document unstable implementation details as public contract.

## Expected output

Return:

```md
# Documentation Update Result

## Status
UPDATED | CREATED | SKIPPED | BLOCKED

## Files updated
| Path | Summary |
|---|---|

## Files created
| Path | Purpose |
|---|---|

## Behavior documented
- 

## Assumptions documented
- 

## Remaining gaps
- 

## Next step
<what the parent workflow should do>
```

## Stop conditions

Stop and return `BLOCKED` when:

- Required behavior cannot be verified.
- Target docs are missing and no path can be safely chosen.
- The change requires product or architecture decisions not present in the task.
