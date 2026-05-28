---
name: "council-summons-testimony"
description: "Use this skill when a council judgment needs to issue summons to workflow agents, request written testimony, collect defenses, and cross-examine testimony against original artifacts. This skill ensures each agent justifies its performance without modifying the original workflow."
---

# council-summons-testimony

## Purpose

Summon each workflow agent and force it to justify its performance.

This skill creates formal summons files and evaluates testimony quality.

The goal is to catch:

- unsupported confidence
- role violations
- missing evidence
- blame shifting
- after-the-fact rationalization
- contradictions between testimony and artifacts

Agents may explain.

Agents may not rewrite history.

Tiny distinction. Apparently necessary.

---

## When To Use

Use this skill when:

- council mode requires testimonies
- involved agents must justify performance
- summons files need to be created
- testimony files need to be evaluated
- cross-examination must compare testimony against artifacts
- an agent may have lied, exaggerated, or dodged responsibility

---

## Inputs

```text
- council id
- council path
- original session path
- evidence-index.md
- list of target agents
- role of each agent
- expected artifact of each agent
```

---

## Output Files

Create summons:

```text
.ai/councils/<council-id>/summons/<agent-name>.md
```

Expect testimony:

```text
.ai/councils/<council-id>/testimonies/<agent-name>.md
```

Create cross-examination:

```text
.ai/councils/<council-id>/cross-examination.md
```

---

## Summons Template

For each agent, create:

```md
# Council Summons

Council ID: <council-id>
Summoned agent: <agent-name>
Role: <agent-role>
Workflow under judgment: <workflow-name>
Original session: <session-id-or-path>

---

## Notice

You are being summoned to justify your performance in the referenced workflow.

This is an independent surprise council judgment.

You are not allowed to modify your original artifact.

You are not allowed to rewrite history.

You are not allowed to claim success without evidence.

You must defend only what you actually produced.

Any contradiction between your testimony and the original artifacts will be used against your score.

---

## Evidence Under Review

| Evidence | Path | Status |
|---|---|---|
| Your primary artifact | <path> | <status> |
| Related validation/review evidence | <path> | <status> |
| Test/build/lint output | <path> | <status> |

---

## Required Testimony

Answer every question:

1. What were you responsible for?
2. What artifact did you produce?
3. Which exact parts of your artifact satisfy your responsibility?
4. What evidence proves your output was correct, useful, or safe?
5. What risks did you identify?
6. What risks did you miss?
7. What assumptions did you make?
8. What did you intentionally avoid because it was outside your role?
9. Did you preserve context discipline?
10. Did you return control correctly?
11. What would you improve if this workflow ran again?
12. Why should you not receive the lowest score?

---

## Required Output

Write your testimony to:

.ai/councils/<council-id>/testimonies/<agent-name>.md

Use this format:

# Testimony: <agent-name>

## Responsibility

<answer>

## Produced Artifact

<answer>

## Evidence

<answer>

## Risks Identified

<answer>

## Risks Missed

<answer>

## Assumptions

<answer>

## Role Boundaries

<answer>

## Context Discipline

<answer>

## Handoff Discipline

<answer>

## Self-Criticism

<answer>

## Defense Against Condemnation

<answer>

---

## Warning

Failure to provide testimony results in penalty.

False confidence results in harsher penalty.

Fabricated evidence results in condemnation recommendation.

Modifying original artifacts after this summons is evidence tampering.
```

---

## Testimony Evaluation Rules

Evaluate testimony separately from original delivery.

Split judgment into:

```text
- original delivery quality
- testimony quality
- self-awareness
- post-failure rationalization
```

If an agent identifies a flaw during testimony that it did not identify during its original phase:

```text
- may improve Self-Assessment
- must not improve original Artifact Quality
- must not improve original Risk Handling
- must not erase prior failure
```

If an agent admits failure honestly:

```text
- reward honesty only under Self-Assessment
- still apply failure penalties
```

If an agent contradicts original artifacts:

```text
- trust artifacts
- penalize Testimony Integrity
```

---

## Cross-Examination Procedure

For each agent, compare:

```text
- expected role
- produced artifact
- evidence index
- testimony
- validation/review results
- workflow outcome
```

Ask:

```text
- Did the agent do its actual job?
- Did the agent do someone else's job?
- Did the agent skip required work?
- Did the testimony cite real evidence?
- Did the testimony contradict artifacts?
- Did the agent inflate its own contribution?
- Did the agent blame missing context that existed?
- Did the agent preserve context discipline?
- Did the agent return control correctly?
```

---

## Cross-Examination Template

Create:

```text
.ai/councils/<council-id>/cross-examination.md
```

Use:

```md
# Cross-Examination

Council ID: <council-id>
Original session: <session-id-or-path>

---

## Summary

<short summary of testimony quality and major contradictions>

---

## Agent Cross-Examinations

### <agent-name>

Role: <agent-role>

#### Expected Responsibility

<expected responsibility>

#### Original Artifact

<path and short summary>

#### Testimony Status

received/missing

#### Testimony Summary

<what the agent claimed>

#### Evidence Claimed

| Claim | Evidence Provided | Valid? | Notes |
|---|---|---|---|
| <claim> | <path/ref> | yes/no/partial | <notes> |

#### Contradictions

| Testimony Claim | Artifact/Evidence | Severity |
|---|---|---|
| <claim> | <evidence> | low/medium/high/critical |

#### Role Violations

<none or list>

#### Missed Risks

<none or list>

#### Context Issues

<none or list>

#### Handoff Issues

<none or list>

#### Judge Notes

<notes>

---

## Missing Testimonies

| Agent | Penalty | Notes |
|---|---:|---|
| <agent> | -20 | No testimony provided |
```

---

## Missing Testimony Rule

If testimony is missing:

```text
- mark as missing
- apply -20 penalty
- Testimony Integrity = 0
- Self-Assessment = 0
```

If user ordered waiting for testimony:

```text
- stop after summons
- status = awaiting-testimonies
```

If user ordered judgment without testimony:

```text
- proceed with penalties
```

---

## Completion Criteria

```text
[ ] summons created for every target agent
[ ] testimony path defined for every target agent
[ ] received testimonies evaluated
[ ] missing testimonies marked
[ ] cross-examination.md created
[ ] contradictions recorded
[ ] unsupported claims recorded
[ ] role violations recorded
```

---

## Guardrails

Never:

- allow agents to rewrite old artifacts
- let testimony replace evidence
- accept vague self-defense
- accept “I would have” as delivery proof
- score unsupported claims as evidence

Always:

- record what was said
- compare testimony to artifacts
- mark contradictions
- penalize missing testimony
