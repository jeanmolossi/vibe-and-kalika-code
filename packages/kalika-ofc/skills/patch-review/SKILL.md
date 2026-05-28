---
name: patch-review
description: Use when reviewing a current patch against the task, plan, and validation evidence. Focus on correctness, security, architecture boundaries, and only raise important issues that affect shipping.
---

# Patch Review Skill

## Overview

Review the patch as a delivery artifact, not as a pile of changed files.

The goal is simple: decide whether this patch is safe, correct, and aligned with the task. Do not waste the review on cosmetic noise when the real risk is elsewhere. The reviewer should spend attention where failures would hurt: behavior, data, security, boundaries, regressions, and missing evidence.

Use the linked references as lenses, not as a checkbox ceremony.

## When to Use

Use this skill when:

- A patch, diff, or changed file set exists and needs review.
- You need to compare implementation against task, plan, or acceptance criteria.
- Validation evidence exists and must be judged before approval.
- You need a final verdict for a review workflow.
- The change set is complete enough to evaluate shipping risk.

Do not use this skill when:

- There is no patch or no changed-file list.
- The work is still exploratory or in planning.
- The task is to implement fixes, not review them.
- Only documentation is changing and there is no behavior risk.
- You are being asked only to design a validation strategy.

## Review Goal

Answer these questions in order:

1. Does the patch satisfy the actual task?
2. Could it break something users depend on?
3. Does it introduce security or data risks?
4. Does it respect architecture and layer boundaries?
5. Is the validation evidence enough to trust the change?
6. Are the remaining issues important enough to block shipping?

If the answer to any of the first four is “no” or “maybe with real risk”, stop treating the patch like a style exercise. Focus there.

## Review Strategy

### 1. Read the contract first

Start with the task, acceptance criteria, and any plan.

Look for:

- what was explicitly requested
- what must not change
- constraints from the repo or environment
- known edge cases
- behavior expected by downstream users

If the patch solves the wrong problem, the review is already negative.

### 2. Triage by impact

Review in this order:

1. Security, auth, secrets, data loss, privacy
2. Correctness and regressions
3. Architecture and layer boundaries
4. Tests and validation evidence
5. Maintainability and clarity
6. Cosmetic issues only if they hide risk

Do not give equal weight to everything. A mislabeled variable is not the same as a broken access control path.

### 3. Inspect behavior, not just syntax

Ask:

- What changed in runtime behavior?
- What inputs are new or now handled differently?
- What failure modes were added or removed?
- Are edge cases covered, or just the happy path?
- Is the new behavior consistent with existing conventions?

### 4. Check the boundaries

This matters more than people admit.

Make sure the patch does not:

- leak UI concerns into domain logic
- push persistence details into presentation code
- mix orchestration with business rules
- let shared utilities become dumping grounds
- bypass existing abstractions without reason

If the patch crosses a layer boundary, ask whether it is justified or just convenient.

### 5. Judge evidence, not hope

Validation should match the risk.

Good signs:

- tests cover the changed behavior
- evidence exercises the risky path, not only the happy path
- failures were considered and handled
- outputs or logs show the important branch actually ran

Bad signs:

- “looks fine” without proof
- tests that do not touch the changed path
- screenshots or logs that prove nothing
- vague claims that “it works locally”

### 6. Keep the review high-signal

Do not invent problems just to sound thorough.

Raise an issue only if it is:

- a real bug or likely regression
- a security risk
- a data integrity risk
- a boundary violation that will hurt maintenance
- a missing validation item that matters for this patch

Style comments are fine only when they improve correctness, readability, or future safety. Otherwise, let them live.

## Review Lenses

Use the reference files below when the patch touches those areas:

- `references/clean-code.md` — naming, duplication, complexity, readability
- `references/security.md` — auth, secrets, injection, unsafe input handling, sensitive data
- `references/design-patterns.md` — when a pattern is justified vs when it is cargo cult
- `references/refactoring.md` — safe incremental change vs risky rewrite
- `references/layer-responsabilities.md` — ownership boundaries between layers
- `references/common-mistakes.md` — recurring review traps and noise to ignore

## Required Review Questions

Before writing the verdict, confirm:

- What is the exact behavior change?
- Which users or flows are affected?
- What could break in production?
- What security risk is introduced, if any?
- Is the architecture cleaner or just rearranged?
- Does the validation actually cover the changed path?
- What is the smallest fix, if the patch is not ready?

## Output Expectations

Return a concise review with a real verdict.

Use this shape:

```md
# Patch Review Result

## Verdict

PASS | PASS_WITH_WARNINGS | BLOCKED | REJECTED

## Summary

<one short paragraph>

## What changed

- <bullets with the important behavior changes>

## Blocking issues

| File | Issue | Risk | Required fix |
| ---- | ----- | ---- | ------------ |

## Non-blocking issues

| File | Issue | Recommendation |
| ---- | ----- | -------------- |

## Missing evidence

- <only if relevant>

## Required next step

<what the parent workflow should do>
```

If there are no blocking issues, do not manufacture some. Say PASS or PASS_WITH_WARNINGS and be precise about why.

## Stop Conditions

Stop and return `BLOCKED` when:

- the diff or changed-file list is unavailable
- the original task or acceptance criteria are missing
- validation evidence is required but absent
- the patch changes unrelated areas without justification
- a serious issue cannot be verified with the available context
- the review cannot distinguish intended behavior from accidental side effect

## Common Pitfalls

1. Reviewing style before safety.

   Fix: check risk first. Clean prose does not save broken behavior.

2. Treating every lint warning as a blocker.

   Fix: only block on warnings that point to real bugs, real debt, or real policy violations.

3. Ignoring the task and reviewing “good code” instead.

   Fix: compare against the contract, not your preferences.

4. Calling something “fine” without evidence.

   Fix: tie verdicts to test results, diffs, or explicit reasoning.

5. Over-rotating on abstractions.

   Fix: simpler code is usually better unless the complexity buys something real.

## Verification Checklist

- [ ] Original task or acceptance criteria were read
- [ ] Changed files or diff were inspected
- [ ] High-risk areas were reviewed first
- [ ] Security and data risks were checked
- [ ] Layer boundaries were checked
- [ ] Validation evidence matched the change risk
- [ ] The verdict reflects impact, not vibes
- [ ] Non-issues were not promoted into blockers
