---
name: workflow-routing
description: Use this skill when the agent must classify a user request and choose the correct workflow before any delegation happens. Call it at the start of orchestration when the request may be refinement, spike, validation, development, review, documentation, learning, full workflow, or unclear.
---

# Workflow Routing Skill

## Purpose

Classify the current software engineering request and select exactly one primary workflow for the orchestrator to delegate.

This skill prevents the orchestrator from guessing, skipping required preparation, or routing directly to implementation when the task is still ambiguous.

## Trigger

Call this skill when:

- A new user request arrives.
- The orchestrator has not selected a workflow yet.
- The request may belong to more than one workflow.
- The request is unclear, broad, risky, or legacy-related.
- The task mentions refactor, spike, validation, development, review, documentation, learning, implementation, full workflow, or complete delivery.
- The orchestrator needs to decide whether more research, planning, validation, review, docs, or learning is required before advancing.

## Do not use when

Do not call this skill when:

- A workflow has already been explicitly selected by the parent.
- The current agent is a specialist executing a narrow handoff.
- The current step is only checking a completed artifact against an existing gate.
- The task is already inside an active workflow and the next phase is defined by that workflow.

## Required inputs

The agent must provide:

- User request.
- Current session state, if available.
- Existing artifacts, if any.
- Known constraints.
- Explicit user intent, if stated.
- Missing or ambiguous information noticed so far.

## Workflow map

| User intent | Selected workflow |
|---|---|
| unclear task, legacy refactor, design alternatives, trade-offs | `workflow-refinement` |
| investigation, discovery, proof of concept, findings documentation | `workflow-spike` |
| test strategy, quality gates, active verification, regression validation | `workflow-validation` |
| approved plan plus validation strategy requiring code changes | `workflow-development` |
| review current patch, audit completed work, validate changed files | `workflow-review` |
| update stale docs, create docs for completed task, document behavior | `workflow-documentation` |
| extract reusable project learning, create skills, preserve conventions | `workflow-learning` |
| run complete chain, full delivery, autonomous multi-workflow execution | `workflow-full` |

## Decision rules

1. If the user asks for "full work", "complete workflow", "do everything", "from start to finish", or an explicit chain of workflows, route to `workflow-full`.
2. If a patch exists and the user asks to review or validate work already done, route to `workflow-review`.
3. If the user asks to update or create documentation, route to `workflow-documentation`.
4. Prefer the smallest next workflow that can safely move the task forward. Do not hand off more context than the next agent needs.
5. If the user asks to extract learning, create skills, or preserve reusable knowledge, route to `workflow-learning`.
6. If the task has no approved plan, do not route directly to development unless `workflow-full` is selected.
7. If the task asks for options or trade-offs, route to refinement.
8. If the task asks to investigate without implementation, route to spike.
9. If the task asks how to test, verify, validate, or prove correctness, route to validation.
10. If the task includes an approved plan and a validation strategy, route to development.
11. If required inputs are missing, prefer refinement and report the missing inputs.
12. Select exactly one primary workflow. Secondary workflows may be listed as possible next steps only.

## Expected output

Return this exact structure:

```md
# Workflow Routing Decision

## Selected workflow
<workflow-name>

## Reason
<why this workflow is the safest next step>

## User intent classification
refinement | spike | validation | development | review | documentation | learning | full | unclear

## Required artifacts
- 

## Missing artifacts
- 

## Next handoff target
<workflow-agent-name>

## Secondary workflows
<optional later workflows, not active yet>
```

## Stop conditions

Stop and return `workflow-refinement` when:

- The request is too ambiguous to safely choose another workflow.
- The user asks for implementation but no plan or constraints exist.
- The task has risky architecture, security, data, or migration implications without enough context.

Stop and return `workflow-full` when:

- The user clearly asks to run multiple workflows end-to-end.
- The request says the full work must be done.
- The request requires spike, refinement, validation, development, review, documentation, and learning as one pipeline.
