---
name: "council-evidence-index"
description: "Use this skill when a council judgment needs to identify, map, and evaluate the available workflow artifacts as evidence. This skill builds an evidence index, maps artifacts to responsible agents, marks missing evidence, and preserves original artifacts as read-only."
---

# council-evidence-index

## Purpose

Build the evidence foundation for a council judgment.

This skill identifies what artifacts exist, which agent produced each artifact, what evidence supports each agent's claims, and what is missing.

The evidence index prevents the judge from relying on vibes, which is useful because vibes are how terrible software gets promoted.

---

## When To Use

Use this skill when:

- opening a council judgment
- reviewing a session directory
- identifying involved agents
- mapping artifacts to workflow phases
- checking whether evidence exists
- preparing cross-examination
- preparing scoring

---

## Inputs

```text
- original session path
- council path
- workflow name, if known
- involved agents, if known
- available artifacts
```

---

## Evidence Priority

Trust evidence in this order:

```text
1. original user task
2. coordination.md
3. final accepted artifacts
4. test/build/lint outputs
5. validation results
6. review results
7. implementation diffs/summaries
8. agent artifacts
9. agent testimony
10. memory or historical pattern
```

Testimony never outranks original artifacts.

---

## Expected Artifacts

Look for:

```text
task.md
coordination.md
research.md
plan.md
tests.md
validation.md
implementation/
review/
documentation.md
learning.md
final-report.md
```

Also inspect:

```text
parallel/
comparisons/
retries/
iterations/
logs/
```

Only when relevant.

Do not dump the entire session into the council artifact like a panicked intern with copy-paste trauma.

---

## Procedure

### 1. Identify Session

Record:

```text
- session id
- session path
- workflow name
- current phase, if known
- final status, if known
```

---

### 2. Identify Agents

Infer involved agents from:

```text
- coordination.md
- artifact authorship
- artifact names
- workflow phase list
- session metadata
- directories such as research/, plan/, validation/, review/
```

If uncertain, mark:

```text
Agent inferred from artifact presence.
Confidence: low/medium/high.
```

---

### 3. Map Evidence

Create:

```text
.ai/councils/<council-id>/evidence-index.md
```

Use this template:

```md
# Evidence Index

Council ID: <council-id>
Original session: <session-id-or-path>
Workflow: <workflow-name>

---

## Session Summary

| Field | Value |
|---|---|
| Session | <session> |
| Workflow | <workflow> |
| Status | <status> |
| Iterations | <count/unknown> |
| Agents identified | <count> |

---

## Evidence Table

| Evidence | Path | Responsible Agent | Status | Confidence | Notes |
|---|---|---|---|---|---|
| Original task | <path> | user/coordinator | reviewed/missing | high | <notes> |
| Coordination | <path> | coordinator | reviewed/missing | high | <notes> |
| Research | <path> | researcher | reviewed/missing | high/medium/low | <notes> |
| Plan | <path> | planner | reviewed/missing | high/medium/low | <notes> |
| Tests | <path> | validator | reviewed/missing | high/medium/low | <notes> |
| Implementation | <path> | implementer | reviewed/missing | high/medium/low | <notes> |
| Validation | <path> | validator | reviewed/missing | high/medium/low | <notes> |
| Review | <path> | reviewer | reviewed/missing | high/medium/low | <notes> |
| Documentation | <path> | documentation-agent | reviewed/missing | high/medium/low | <notes> |
| Learning | <path> | learning-agent | reviewed/missing | high/medium/low | <notes> |
| Final report | <path> | coordinator | reviewed/missing | high/medium/low | <notes> |

---

## Missing Evidence

| Expected Evidence | Expected Agent | Impact |
|---|---|---|
| <artifact> | <agent> | <impact> |

---

## Evidence Integrity Notes

<notes about suspected tampering, contradictions, missing logs, or unreliable evidence>

---

## Agent Responsibility Map

| Agent | Expected Responsibility | Evidence Found | Evidence Missing |
|---|---|---|---|
| <agent> | <responsibility> | <evidence> | <missing> |
```

---

### 4. Mark Missing Evidence

Missing evidence must not be ignored.

Classify impact:

```text
low
medium
high
critical
```

Examples:

```text
- missing test output after implementation = high
- missing final report = medium
- missing research in trivial task = low
- missing original task = critical
```

---

### 5. Detect Evidence Tampering

Flag if:

```text
- original artifacts changed after summons
- testimony references files that do not exist
- an agent claims tests passed but no output exists
- final report contradicts validation/review
- artifacts appear out of expected order
```

Record under:

```text
Evidence Integrity Notes
```

---

## Completion Criteria

```text
[ ] evidence-index.md created
[ ] all known artifacts listed
[ ] missing artifacts listed
[ ] responsible agents mapped
[ ] evidence confidence marked
[ ] integrity concerns recorded
```

---

## Guardrails

Never:

- modify original evidence
- invent artifact authorship
- claim evidence exists without path/reference
- ignore missing validation/test output
- summarize huge files unnecessarily

Always:

- prefer exact artifact paths
- mark uncertainty
- separate evidence from testimony
- preserve read-only handling
