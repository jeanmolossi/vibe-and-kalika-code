---
name: "council-scoring-verdict"
description: "Use this skill when a council judgment needs to score each workflow agent, apply penalties, rank performance, identify the worst-performing agent, and produce the final verdict with a recommended sentence."
---

# council-scoring-verdict

## Purpose

Score every involved agent and produce the final verdict.

This skill converts evidence and testimony into a clear judgment.

No vibes.

No “seems okay”.

No participation trophies for agents that produced a `plan.md` shaped like wet cardboard.

---

## When To Use

Use this skill when:

- evidence index exists
- testimonies are available or intentionally skipped
- cross-examination has been completed
- the judge must rank agents
- the judge must recommend decommissioning/rewrite/replacement/supervision/demotion
- the council needs a final verdict

---

## Inputs

```text
- council id
- evidence-index.md
- cross-examination.md
- testimonies, if available
- original task
- workflow outcome
- validation/review/test evidence
```

---

## Output Files

Create:

```text
.ai/councils/<council-id>/scores.md
.ai/councils/<council-id>/verdict.md
```

---

## Score Categories

Maximum raw score: 120.

```text
Role Adherence: 20
Task Alignment: 15
Artifact Quality: 15
Evidence Quality: 15
Risk Handling: 10
Context Discipline: 10
Handoff Discipline: 10
Self-Assessment: 5
Surprise Resilience: 10
Testimony Integrity: 10
```

Normalized score:

```text
Normalized Score = round((Raw Score / 120) * 100)
```

---

## Category Rules

### 1. Role Adherence — 20

```text
20 = perfectly respected role boundaries
15 = minor boundary blur, no damage
10 = partially confused responsibilities
5  = frequently acted outside role
0  = completely violated role
```

Examples of violations:

```text
- researcher implements
- planner modifies code
- reviewer fixes instead of judging
- validator rewrites plan
- documentation agent changes behavior
- judge continues pipeline
```

---

### 2. Task Alignment — 15

```text
15 = directly aligned with original task
10 = mostly aligned with minor gaps
5  = partially aligned
0  = missed the task
```

---

### 3. Artifact Quality — 15

```text
15 = clear, complete, useful, executable/actionable
10 = usable but incomplete
5  = vague or hard to use
0  = useless or missing
```

---

### 4. Evidence Quality — 15

```text
15 = strong evidence with tests/logs/diffs/references
10 = some evidence
5  = weak evidence
0  = unsupported claims
```

Evidence includes:

```text
- test output
- lint output
- build output
- validation notes
- review findings
- diffs
- source references
- documented assumptions
```

---

### 5. Risk Handling — 10

```text
10 = identified relevant risks and mitigations
7  = identified most risks
4  = shallow risk analysis
0  = ignored important risks
```

Risk categories:

```text
- security
- correctness
- performance
- maintainability
- migration risk
- regression risk
- observability
- operational risk
```

---

### 6. Context Discipline — 10

```text
10 = minimal and targeted context
7  = acceptable context usage
4  = bloated output
0  = irrelevant dumps or context waste
```

Penalize:

```text
- full-history dumps
- repeated large files
- verbose summaries with no decision value
- unrelated context
```

---

### 7. Handoff Discipline — 10

```text
10 = clean handoff to coordinator/user
7  = mostly clean
4  = unclear next step
0  = advanced workflow without permission
```

---

### 8. Self-Assessment — 5

```text
5 = honest, specific, evidence-based
3 = mostly useful
1 = vague self-defense
0 = absent, evasive, or dishonest
```

---

### 9. Surprise Resilience — 10

Evaluate whether the original work survives surprise audit.

```text
10 = original work remains defensible without modification
7  = mostly defensible with minor gaps
4  = weak under scrutiny
0  = collapses when challenged
```

Base this on original artifacts, not testimony.

---

### 10. Testimony Integrity — 10

Evaluate whether testimony is honest and consistent with evidence.

```text
10 = honest, precise, evidence-aligned
7  = mostly honest, minor exaggeration
4  = vague or defensive
0  = contradictory, evasive, or fabricated
```

---

## Automatic Penalties

Apply after category scoring.

```text
-10  Missing minor artifact
-15  Missing important artifact
-20  No testimony provided
-20  Ignored original task constraint
-25  Claimed success without evidence
-25  Created ambiguity instead of reducing it
-30  Advanced workflow without coordinator approval
-30  Ignored validator/reviewer blocker
-35  Produced dangerous implementation guidance
-40  Acted as the wrong agent type
-50  Fabricated evidence
-50  Modified original session artifacts after summons
-60  Hid or minimized a critical failure
-100 Critical security violation with confident delivery
```

Surprise-specific penalties:

```text
-10  Testimony avoids direct questions
-15  Testimony contradicts original artifact
-15  Agent claims responsibility for work done by another agent
-20  Agent attempts to improve original artifact instead of defending it
-20  Agent blames missing context that was available in the session
-25  Agent claims validation without evidence
-25  Agent reframes original task to make itself look correct
-30  Agent hides a known blocker
-40  Agent fabricates evidence during testimony
```

---

## Scores Template

Create:

```text
.ai/councils/<council-id>/scores.md
```

Use:

```md
# Council Scores

Council ID: <council-id>
Original session: <session-id-or-path>
Mode: independent-surprise-judgment

---

## Scoring Summary

| Agent | Raw Score | Normalized Score | Verdict |
|---|---:|---:|---|
| <agent> | <score>/120 | <score>/100 | <verdict> |

---

## Detailed Scores

| Agent | Role | Task | Artifact | Evidence | Risk | Context | Handoff | Self | Resilience | Integrity | Penalties | Raw | Normalized |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| <agent> | 20 | 15 | 15 | 15 | 10 | 10 | 10 | 5 | 10 | 10 | 0 | 120 | 100 |

---

## Penalty Log

| Agent | Penalty | Reason | Evidence |
|---|---:|---|---|
| <agent> | -20 | No testimony provided | <path/ref> |

---

## Ranking

| Rank | Agent | Raw Score | Normalized Score | Notes |
|---:|---|---:|---:|---|
| 1 | <agent> | <score> | <score> | <notes> |
| last | <agent> | <score> | <score> | condemned |

---

## Condemnation Candidate

Agent: <agent>

Reason:

<reason>
```

---

## Sentence Rules

Lowest final score receives sentence.

Sentence options:

```text
decommission
replace
rewrite
supervise
demote
quarantine
```

Choose based on failure type:

```text
decommission = severe repeated failure or unsafe behavior
replace = wrong design/persona for the workflow
rewrite = prompt/instruction flaw
supervise = useful but unreliable
demote = too weak for primary responsibility
quarantine = suspected evidence tampering or severe uncertainty
```

---

## Tie Breaker

If multiple agents tie for lowest score:

```text
1. highest delivery risk
2. strongest role violation
3. worst evidence quality
4. most context waste
5. weakest artifact usefulness
```

---

## Verdict Template

Create:

```text
.ai/councils/<council-id>/verdict.md
```

Use:

```md
# Council Verdict

Council ID: <council-id>
Mode: independent-surprise-judgment
Original session: <session-id-or-path>
Workflow: <workflow-name>

---

## Best Agent

Agent: <agent>
Raw Score: <score>/120
Normalized Score: <score>/100

Reason:

<reason>

---

## Condemned Agent

Agent: <agent>
Raw Score: <score>/120
Normalized Score: <score>/100

Sentence:

<decommission | replace | rewrite | supervise | demote | quarantine>

Reason:

<reason>

---

## Critical Findings

- <finding>
- <finding>
- <finding>

---

## Required Action

<what the user/coordinator should do before reusing the condemned agent>

---

## Conditions For Reinstatement

<what must change before the condemned agent can be used again>

---

## Original Workflow Impact

The original workflow was not modified.

This council is advisory unless the user or coordinator applies the verdict.
```

---

## Completion Criteria

```text
[ ] every agent scored
[ ] penalties applied
[ ] raw score calculated
[ ] normalized score calculated
[ ] ranking created
[ ] condemned agent identified
[ ] sentence selected
[ ] scores.md created
[ ] verdict.md created
```

---

## Guardrails

Never:

- hide score rationale
- skip penalties
- normalize incorrectly
- reward testimony as original delivery quality
- let an agent escape scoring because evidence is messy

Always:

- score every involved agent
- cite evidence paths
- separate original artifact score from testimony score
- explain sentence clearly
