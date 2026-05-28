# Plan Quality Gate

Use this checklist before marking an implementation plan as final.

## Blocking checks

The plan is BLOCKED if any answer is `NO`.

- [ ] Does the plan reference a concrete PRD source?
- [ ] Are all blocking questions answered or explicitly accepted as assumptions by the human?
- [ ] Are requirements mapped with stable IDs?
- [ ] Are acceptance criteria mapped to validation steps?
- [ ] Are implementation tasks atomic?
- [ ] Does every task have traceability to a requirement, acceptance criterion, non-functional requirement, accepted assumption, or risk mitigation?
- [ ] Does every phase have validation?
- [ ] Does every risk have mitigation?
- [ ] Are rollout and rollback described?
- [ ] Is out-of-scope explicitly listed?
- [ ] Is the plan free from invented requirements?
- [ ] Is the plan clear enough for a low-context implementer?

## Warning checks

The plan can continue, but should be improved if any answer is `NO`.

- [ ] Are likely affected files/modules listed when known?
- [ ] Are test commands listed when known?
- [ ] Are observability needs listed?
- [ ] Are security concerns listed?
- [ ] Are compatibility concerns listed?
- [ ] Are discovery tasks bounded with explicit outputs?

## Final status rule

- If any blocking check fails: `Status: BLOCKED`.
- If all blocking checks pass: `Status: READY`.
