---
name: documentation-audit
description: Use this skill when an agent must find stale, missing, or misleading documentation affected by a task or patch. Call it before updating or creating documentation.
---

# Documentation Audit Skill

## Purpose

Identify documentation that should be updated, created, or intentionally left unchanged based on the current task and patch.

## Trigger

Call this skill when:

- A task changes behavior, architecture, commands, APIs, workflows, configuration, or operational procedures.
- Documentation may be stale.
- The user asks to update project docs.
- Documentation workflow starts.
- A review report indicates documentation gaps.

## Do not use when

Do not call this skill when:

- The task has no behavior or operational impact.
- There is no task summary or patch summary.
- The agent is only reviewing code.
- Documentation scope has already been audited in the current workflow.

## Required inputs

The agent must provide:

- Task summary.
- Implementation summary.
- Changed file list.
- Patch summary.
- Existing documentation paths, if available.
- Repository documentation conventions.
- Review report, if available.

## Audit targets

Inspect likely relevant docs:

- Root README.
- `docs/**/*.md`.
- ADRs.
- API documentation.
- Runbooks.
- Local development docs.
- Testing/validation docs.
- Agent/workflow/skill docs.
- Configuration/deployment docs.

## Expected output

Return:

```md
# Documentation Audit Result

## Status
DOCS_NEEDED | NO_DOCS_NEEDED | BLOCKED

## Relevant documentation
| Path | Current state | Required action |
|---|---|---|

## Missing documentation
| Topic | Suggested path | Reason |
|---|---|---|

## Stale documentation
| Path | Problem | Required update |
|---|---|---|

## Do not update
| Path | Reason |
|---|---|

## Next step
<update existing docs, create new docs, or skip>
```

## Stop conditions

Stop and return `BLOCKED` when:

- The implementation summary is missing.
- The patch scope is unknown.
- Existing docs cannot be inspected but are required.
- The task behavior is ambiguous.
