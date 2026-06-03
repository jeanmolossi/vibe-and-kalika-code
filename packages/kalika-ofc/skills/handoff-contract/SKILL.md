---
name: handoff-contract
description: Use this skill when one agent must delegate work to another agent. Call it before invoking any workflow agent or specialist agent so the callee receives a strict goal, scope, inputs, constraints, expected artifacts, and return rule.
---

# Handoff Contract Skill

## Purpose

Create a strict delegation contract between orchestrator, workflow agents, and specialist agents.

This skill stops agents from freelancing, expanding scope, skipping artifacts, or advancing the pipeline without permission.

## Trigger

Call this skill when:

- The orchestrator delegates to a workflow agent.
- A workflow agent delegates to a specialist agent.
- A retry or correction loop needs a narrower task.
- Parallel work packages need separate instructions.
- A previous handoff produced vague, incomplete, or out-of-scope output.
- The caller needs to define exactly what the callee must return.

## Do not use when

Do not call this skill when:

- The agent is not delegating work.
- The next step is only local artifact evaluation.
- The required handoff has already been written and accepted.
- The caller cannot identify a specific callee and goal yet. Use `workflow-routing` first.

## Required inputs

The agent must provide:

- Caller agent name.
- Callee agent name.
- One clear goal.
- Minimal relevant context.
- Required artifacts and file paths.
- Scope and out-of-scope boundaries.
- Acceptance criteria.
- Constraints and risks.
- Expected return format.

## Required contract fields

Every handoff must include:

```md
# Handoff Contract

## Caller
<agent name>

## Callee
<agent name>

## Goal
<one clear outcome>

## Inputs
<artifacts, files, snippets, constraints>

## Scope
<what is included>

## Out of scope
<what must not be done>

## Required output artifacts
<files/reports expected>

## Acceptance criteria
<checklist>

## Constraints
<security, performance, dependency, architecture, style>

## Context budget
<what context is allowed and what must not be loaded>

## Return rule
Return to caller only. Do not call another agent. Do not advance the pipeline.
```

## Process

1. Identify the single outcome the callee must produce.
2. Remove unrelated context.
3. Define exact scope and prohibited actions.
4. Define required artifacts and output format.
5. Define acceptance criteria.
6. Add return rule.
7. If the handoff is still broad, split it or route back to planning.
8. Prefer terse, LLM-friendly bullets. Grammar can bend if clarity improves and the meaning stays unambiguous.
9. Pass only the delta the callee needs; do not restate history or duplicate artifacts.

## Expected output

Return a completed `# Handoff Contract` only. Keep it compact, LLM-friendly, and low-token. Do not include implementation, research, review, or validation content inside the contract.

## Stop conditions

Stop and return `BLOCKED` to the parent when:

- The caller cannot define a clear goal.
- The caller does not know which agent should receive the handoff.
- Required artifacts are missing.
- Scope cannot be bounded safely.
