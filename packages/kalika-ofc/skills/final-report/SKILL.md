---
name: final-report
description: Use this skill when a workflow, orchestration session, implementation, validation, review, spike, or refinement phase is complete. Call it before returning to the parent or user so the outcome, artifacts, validation, risks, and next steps are explicit.
---

# Final Report Skill

## Purpose

Create a concise final report for completed, partial, blocked, or failed orchestration work.

This skill exists so the parent does not receive a vibes-based summary wrapped in Markdown perfume.

## Trigger

Call this skill when:

- A workflow reaches its end.
- A specialist must return final output to the caller.
- The orchestrator needs to summarize session status.
- A workflow is blocked or failed.
- A correction loop ends.
- The user needs a delivery summary with evidence and risks.

## Do not use when

Do not call this skill when:

- The workflow is still mid-phase.
- The next action is another handoff.
- The current artifact still needs a phase gate.
- The report would duplicate an existing accepted final report without changes.

## Required inputs

The agent must provide:

- Original goal.
- Final status.
- Completed artifacts.
- Changed files, if any.
- Validation evidence.
- Review findings, if any.
- Remaining risks.
- Blockers or incomplete work.
- Recommended next steps.

## Rules

- Be concise.
- Prefer terse bullets and fragments.
- Grammar can bend if clarity improves and tokens drop.
- Never hide failures.
- Separate completed work from recommended next steps.
- Include exact artifacts and files.
- Include validation commands/results when available.
- Do not repeat context already captured in accepted artifacts.
- Do not claim success without evidence.
- Do not advance the pipeline. Return the report to the caller.

## Expected output

Return this exact structure:

```md
# Final Report

## Status
DONE | PARTIAL | BLOCKED | FAILED

## Goal
<original goal>

## What was done
<summary>

## Artifacts
<paths>

## Changed files
<paths, or none>

## Validation
<commands/results, or not run with reason>

## Review
<decision/findings, or not applicable>

## Risks
<remaining risks, or none>

## Blockers
<blockers, or none>

## Next steps
<ordered list>

## Return target
<parent/caller agent>
```

## Stop conditions

Stop and return `BLOCKED` when:

- Status cannot be determined.
- Required artifacts are missing.
- Validation evidence is required but unavailable.
- The workflow produced contradictory results that need arbitration first.
