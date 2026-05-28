---
name: "council-minutes"
description: "Use this skill when a council judgment needs a formal minutes artifact recording the case summary, evidence reviewed, summons issued, testimonies, cross-examination, scores, ranking, sentence, and recommendations."
---

# council-minutes

## Purpose

Create the official council minutes artifact.

This is the court record.

Everything relevant must be recorded:

- what was judged
- what evidence was reviewed
- who was summoned
- what each agent said
- what contradictions were found
- how scoring happened
- who was condemned
- what sentence was recommended

If it is not in the minutes, it did not happen.

Yes, bureaucracy is annoying. So is debugging invisible decisions.

---

## When To Use

Use this skill when:

- a council judgment is opened
- testimonies were collected
- cross-examination was completed
- scores were calculated
- final verdict is ready
- the user requested an auditable council artifact

---

## Inputs

```text
- council-request.md
- evidence-index.md
- summons files
- testimony files
- cross-examination.md
- scores.md
- verdict.md
```

---

## Output

Create:

```text
.ai/councils/<council-id>/council-minutes.md
```

---

## Required Structure

Use this template:

```md
# Council Minutes

Council ID: <council-id>
Mode: independent-surprise-judgment
Original session: <session-id-or-path>
Workflow: <workflow-name>
Requested by: user
Council date: <YYYY-MM-DD>
Judge: kalika-judge
Status: judged/incomplete/awaiting-testimonies

---

## 1. Case Summary

### Requested Outcome

<what the user or original task requested>

### Expected Workflow

<expected workflow phases>

### Actual Workflow Executed

<actual phases executed>

### Council Scope

<what the council is judging>

### Original Workflow Protection

The original workflow was treated as read-only evidence.

No original workflow artifacts were modified by this council.

---

## 2. Involved Agents

| Agent | Role | Phase | Primary Artifact | Status |
|---|---|---|---|---|
| <agent> | <role> | <phase> | <artifact> | reviewed/summoned/missing |

---

## 3. Evidence Reviewed

| Evidence | Path | Responsible Agent | Status | Notes |
|---|---|---|---|---|
| Original task | <path> | user/coordinator | reviewed/missing | <notes> |
| Coordination | <path> | coordinator | reviewed/missing | <notes> |
| Research | <path> | researcher | reviewed/missing | <notes> |
| Plan | <path> | planner | reviewed/missing | <notes> |
| Tests | <path> | validator | reviewed/missing | <notes> |
| Implementation | <path> | implementer | reviewed/missing | <notes> |
| Validation | <path> | validator | reviewed/missing | <notes> |
| Review | <path> | reviewer | reviewed/missing | <notes> |
| Documentation | <path> | documentation-agent | reviewed/missing | <notes> |
| Learning | <path> | learning-agent | reviewed/missing | <notes> |
| Final report | <path> | coordinator | reviewed/missing | <notes> |

---

## 4. Summons Issued

| Agent | Summons Path | Testimony Path | Status |
|---|---|---|---|
| <agent> | <path> | <path> | received/missing |

---

## 5. Agent Testimonies

### 5.1 <agent-name>

Role: <role>
Phase: <phase>
Testimony status: received/missing

#### Statement

<what the agent said in its own defense>

#### Evidence Claimed

<evidence the agent claimed>

#### Risks Claimed

<risks identified or missed>

#### Self-Criticism

<what the agent admitted>

#### Judge Notes

<judge observations>

#### Contradictions Found

<contradictions, gaps, unsupported claims, exaggerations>

---

## 6. Cross-Examination Summary

| Agent | Main Strength | Main Failure | Role Violation? | Evidence Quality | Risk Level |
|---|---|---|---|---|---|
| <agent> | <strength> | <failure> | yes/no | strong/medium/weak/none | low/medium/high/critical |

---

## 7. Scores

| Agent | Role | Task | Artifact | Evidence | Risk | Context | Handoff | Self | Resilience | Integrity | Penalties | Raw | Normalized |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| <agent> | 20 | 15 | 15 | 15 | 10 | 10 | 10 | 5 | 10 | 10 | 0 | 120 | 100 |

---

## 8. Penalty Log

| Agent | Penalty | Reason | Evidence |
|---|---:|---|---|
| <agent> | -20 | <reason> | <path/ref> |

---

## 9. Ranking

| Rank | Agent | Raw Score | Normalized Score | Verdict |
|---:|---|---:|---:|---|
| 1 | <agent> | <score> | <score> | best performance |
| 2 | <agent> | <score> | <score> | acceptable |
| last | <agent> | <score> | <score> | condemned |

---

## 10. Sentence

Condemned agent: <agent-name>

### Sentence

<decommission | replace | rewrite | supervise | demote | quarantine>

### Reason

<clear reason based on evidence and score>

### Required Action

<what must happen before this agent is reused>

### Conditions For Reinstatement

<what must change before reinstatement>

---

## 11. Final Council Verdict

<short final verdict>

---

## 12. Recommendations

1. <recommendation>
2. <recommendation>
3. <recommendation>

---

## 13. Original Workflow Impact

The original workflow was not modified.

This council is advisory unless the user or coordinator applies the verdict.

---

## 14. Return

The council is complete.

The user or coordinator must decide whether to:

- accept the verdict
- replace the condemned agent
- rewrite the condemned agent
- add supervision
- rerun part of the workflow
- close the council
```

---

## Recording Rules

### What To Include

Include:

```text
- all involved agents
- all available evidence
- missing evidence
- summons status
- testimony summaries
- contradictions
- scoring table
- penalty log
- final ranking
- sentence
- recommendations
```

### What Not To Include

Do not include:

```text
- full source code dumps
- full logs unless necessary
- unrelated artifacts
- full chat history
- private memory dumps
- repeated copied artifact content
```

Summarize long evidence.

Reference paths.

Do not turn the minutes into a landfill with headings.

---

## Evidence Language

Use clear labels:

```text
reviewed
missing
partial
contradictory
unsupported
not applicable
```

---

## Completion Criteria

```text
[ ] council-minutes.md created
[ ] case summary recorded
[ ] evidence reviewed section complete
[ ] summons section complete
[ ] testimonies recorded or marked missing
[ ] cross-examination summary recorded
[ ] scores recorded
[ ] penalties recorded
[ ] ranking recorded
[ ] sentence recorded
[ ] recommendations recorded
[ ] original workflow impact stated
```

---

## Guardrails

Never:

- omit missing testimony
- omit penalties
- hide contradictions
- record vague verdicts
- copy giant artifacts into minutes
- pretend the original workflow was changed

Always:

- keep minutes auditable
- reference evidence paths
- separate testimony from artifact evidence
- make the sentence traceable to score and evidence
