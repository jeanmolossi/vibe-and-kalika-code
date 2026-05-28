---
name: "surprise-council-judgment"
description: "Use this skill when the user manually requests an independent surprise judgment of a completed or ongoing multi-agent workflow. This skill creates a separate council directory, defines scope, preserves the original workflow as read-only evidence, and coordinates evidence indexing, summons, testimony, scoring, council minutes, and verdict."
---

# surprise-council-judgment

## Purpose

Run an independent surprise council judgment against an existing workflow.

This skill is used when the user wants to audit agent performance without changing the original workflow.

The council acts as an external court:

- open the council
- define scope
- preserve original workflow
- create independent artifacts
- summon agents
- collect testimony
- judge original delivery
- produce verdict

The original workflow is evidence, not a workspace.

---

## When To Use

Use this skill when:

- the user manually asks for a council
- the user asks for a surprise judgment
- the user wants to evaluate agents after a workflow
- the user wants agents to justify themselves
- the workflow did not include a judgment phase
- the user wants a separate council artifact
- the user suspects poor agent performance
- the user wants accountability without interfering with current workflows

---

## Do Not Use When

Do not use this skill when:

- the user wants normal code review
- the user wants implementation
- the user wants validation
- the user wants documentation
- the user wants planning
- the user wants the judge to continue the workflow

This skill judges.

It does not produce the original work.

---

## Inputs

Minimum:

```text
- original session path or session id
- original task or task artifact
```

Recommended:

```text
- workflow name
- involved agents
- coordination.md
- produced artifacts
- validation outputs
- review outputs
- test/lint/build outputs
- diffs
```

---

## Output Directory

Create:

```text
.ai/councils/<council-id>/
```

Required files:

```text
council-request.md
evidence-index.md
summons/<agent>.md
testimonies/<agent>.md
cross-examination.md
scores.md
council-minutes.md
verdict.md
```

If testimonies are missing and the user did not request immediate judgment, stop with:

```text
status = awaiting-testimonies
```

---

## Procedure

### 1. Generate Council ID

Use:

```text
YYYY-MM-DD--HH-mm_<short-workflow-name>_council
```

Rules:

```text
- no spaces
- no colons
- lowercase preferred
- keep name short
```

---

### 2. Create Council Directory

Create:

```text
.ai/councils/<council-id>/
.ai/councils/<council-id>/summons/
.ai/councils/<council-id>/testimonies/
```

---

### 3. Create Council Request

Create:

```text
.ai/councils/<council-id>/council-request.md
```

Use this template:

```md
# Council Request

Council ID: <council-id>
Requested by: user
Mode: independent-surprise-judgment
Original session: <session-id-or-path>
Workflow under judgment: <workflow-name>
Requested at: <YYYY-MM-DD HH:mm>
Status: opened

---

## Judgment Scope

This council will evaluate agent performance in the referenced workflow.

The original workflow must not be modified.

The council must create independent artifacts under:

.ai/councils/<council-id>/

---

## Original Workflow Protection

The original workflow is read-only evidence.

The judge and summoned agents must not alter:

- task.md
- coordination.md
- research.md
- plan.md
- tests.md
- validation.md
- implementation artifacts
- review.md
- documentation.md
- learning.md
- final-report.md

---

## Target Agents

| Agent | Role | Expected Artifact | Status |
|---|---|---|---|
| <agent> | <role> | <artifact> | pending evidence review |

---

## Evidence Sources

| Source | Path | Status |
|---|---|---|
| Original task | <path> | pending |
| Coordination state | <path> | pending |
| Research | <path> | pending |
| Plan | <path> | pending |
| Tests | <path> | pending |
| Implementation | <path> | pending |
| Validation | <path> | pending |
| Review | <path> | pending |
| Documentation | <path> | pending |
| Learning | <path> | pending |

---

## Council Status

Current status: opened

Next step:

Issue summons to all target agents.
```

---

### 4. Decide Testimony Mode

If the user wants testimonies:

```text
create summons
wait for testimonies
do not score yet
```

If the user says to judge immediately:

```text
mark all missing testimonies
apply penalties
continue scoring
```

---

## Status Values

Use only:

```text
opened
awaiting-testimonies
ready-for-judgment
judged
incomplete
```

---

## Completion Criteria

This skill is complete when:

```text
[ ] council directory exists
[ ] council-request.md exists
[ ] original workflow protection is recorded
[ ] target agents are listed or marked unknown
[ ] next step is clear
```

---

## Guardrails

Never:

- write inside original session
- modify original workflow artifacts
- continue workflow
- let testimony replace original evidence
- skip council-request.md

Always:

- separate council artifacts from workflow artifacts
- mark original workflow as read-only
- make the next council step explicit
