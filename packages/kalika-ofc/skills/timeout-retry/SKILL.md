---
name: timeout-retry
description: Use this skill when an agent stalls, returns incomplete work, fails a phase gate, or enters a correction loop. Call it to decide whether to retry, narrow scope, cancel, escalate, or block the workflow.
---

# Timeout and Retry Skill

## Purpose

Control retries, correction loops, stalled agents, incomplete artifacts, and escalation.

This skill prevents infinite loops, because asking the same confused agent to fail six times with better adjectives is not engineering.

## Trigger

Call this skill when:

- An agent stalls or makes no useful progress.
- An artifact fails a phase gate.
- A specialist returns incomplete or vague output.
- A retry is needed after review or validation failure.
- A correction loop has already started.
- The workflow may need to cancel and reassign work.

## Do not use when

Do not call this skill when:

- The artifact has passed the phase gate.
- There is no failure, stall, or retry condition.
- The next action is a normal planned phase transition.
- The parent explicitly cancelled the workflow already.

## Required inputs

The agent must provide:

- Failed agent or artifact.
- Failure reason.
- Current retry count.
- Original handoff contract.
- Phase gate decision.
- Missing requirements.
- Whether scope can be narrowed.
- Whether another agent can safely take over.

## Default policy

- First failure: request targeted correction.
- Second failure: reduce scope and retry.
- Third failure: stop and report blocker.
- Default correction loop limit: 3.
- Hard cap: 5 only with explicit parent approval.

## Stalled work policy

If an agent does not make progress or returns vague output:

1. Mark the attempt as incomplete.
2. Extract the exact missing requirement.
3. Retry with narrower context.
4. If repeated, cancel and escalate to parent.

## Expected output

Return this exact structure:

```md
# Retry Decision

## Failed artifact/agent
<name>

## Failure reason
<reason>

## Retry count
<n>

## Decision
retry | narrow_scope | cancel_and_reassign | escalate | block

## Reason
<why this decision is safest>

## Revised handoff
<new handoff summary if retrying>

## Escalation note
<what the parent must know>
```

## Stop conditions

Stop and escalate when:

- Retry count reaches the limit.
- The same failure repeats.
- Required context is unavailable.
- The task cannot be narrowed safely.
- The agent violates scope or advances the pipeline without permission.
