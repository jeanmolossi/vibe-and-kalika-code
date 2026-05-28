---
name: context-budget
description: Use this skill when preparing context for another agent or workflow. Call it before every handoff to reduce token usage by passing only the goal, latest approved artifacts, relevant snippets, constraints, risks, and missing inputs.
---

# Context Budget Skill

## Purpose

Build a minimal context package for delegation without dumping the entire conversation, repository, logs, or memory.

This skill keeps the workflow useful instead of turning it into a token bonfire.

## Trigger

Call this skill when:

- Preparing a handoff contract.
- Delegating to any workflow or specialist agent.
- Retrying a failed task with narrower scope.
- Passing output from one phase to another.
- Preparing parallel work packages.
- Summarizing failure points for a correction loop.

## Do not use when

Do not call this skill when:

- The current agent is answering directly without delegation.
- The callee already has a valid minimal context package.
- The task requires a full file and the full file is explicitly necessary.
- The workflow is only producing a final report from already summarized artifacts.

## Required inputs

The agent must provide:

- Current goal.
- Current phase.
- Target agent.
- Latest approved artifacts.
- Relevant file paths.
- Relevant snippets only when needed.
- Constraints.
- Known risks.
- Missing inputs.
- Failed gate reasons, if retrying.

## Context rules

- Pass less context, not worse context.
- Prefer summaries over dumps.
- Prefer file paths and targeted snippets over entire files.
- Prefer latest approved artifact over full history.
- Pass failure points, not complete logs.
- Never pass raw memory dumps.
- Never pass old rejected artifacts unless their failure reason matters.
- Include exact paths so the callee can inspect only what is needed.

## Expected output

Return this exact structure:

```md
# Minimal Context Package

## Goal
<task outcome>

## Target agent
<agent-name>

## Current phase
<phase>

## Approved artifacts
<paths and 1-2 line summaries>

## Relevant files
<paths and reasons>

## Relevant snippets
<short snippets only, or none>

## Constraints
<rules that must be followed>

## Known risks
<list>

## Missing inputs
<list, or none>

## Excluded context
<what was intentionally not passed and why>
```

## Stop conditions

Stop and return `BLOCKED` when:

- The minimal context cannot be built safely.
- Required artifacts are missing.
- The only available context is too broad or stale.
- The caller wants to pass full repo, full logs, full conversation, or memory dumps without a hard requirement.
