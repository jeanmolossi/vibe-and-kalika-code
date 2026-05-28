# Implementation Plan: <title>

Status: READY | BLOCKED
Generated at: <YYYY-MM-DDTHH:mm:ssZ>
Source PRD: <path-or-inline-reference>
Session: <session-id>

## 1. Context summary

<short summary of the product goal and implementation target>

## 2. Source requirement map

| ID | Type | Requirement | Source section | Implementation impact |
|---|---|---|---|---|
| REQ-001 | Functional | <requirement> | <PRD section> | <impact> |
| AC-001 | Acceptance | <criterion> | <PRD section> | <impact> |
| NFR-001 | Non-functional | <requirement> | <PRD section> | <impact> |

## 3. Accepted assumptions

Only include assumptions explicitly accepted by the human.

| ID | Assumption | Accepted by | Risk | Mitigation |
|---|---|---|---|---|
| ASM-001 | <assumption> | <human/user> | <risk> | <mitigation> |

## 4. Out of scope

- <item explicitly out of scope>

## 5. Architecture impact

### Current architecture impact

- <component/module/system affected>

### Proposed implementation shape

- <high-level technical shape>

### Constraints

- <constraint>

## 6. Data model impact

| ID | Change | Required? | Notes |
|---|---|---|---|
| DATA-001 | <change or none> | YES/NO | <notes> |

## 7. API and integration impact

| ID | API/Integration | Change | Contract impact | Compatibility |
|---|---|---|---|---|
| API-001 | <endpoint/service> | <change> | <request/response/event impact> | <backward-compatible or not> |

## 8. Frontend/UI impact

| ID | UI area | Change | State behavior | Validation |
|---|---|---|---|---|
| UI-001 | <screen/component> | <change> | <state impact> | <validation> |

## 9. Backend/domain impact

| ID | Domain/service area | Change | Business rule impact |
|---|---|---|---|
| BE-001 | <domain/service> | <change> | <rule impact> |

## 10. Security and compliance

| ID | Concern | Impact | Required control |
|---|---|---|---|
| SEC-001 | <concern> | <impact> | <control> |

## 11. Observability

| ID | Signal | Where | Purpose |
|---|---|---|---|
| OBS-001 | <metric/log/trace/event> | <location> | <purpose> |

## 12. Rollout strategy

- Strategy: <feature flag / A-B test / staged rollout / direct release>
- Disabled behavior: <expected behavior when disabled>
- Rollback: <rollback path>
- Production safety checks:
  - <check>

## 13. Test strategy

### Unit tests

- TEST-001: <unit test target and cases>

### Integration tests

- TEST-002: <integration test target and cases>

### E2E/manual validation

- TEST-003: <manual or E2E scenario>

### Regression tests

- TEST-004: <existing behavior that must not break>

## 14. Implementation phases

## Phase 1: <name>

ID: PHASE-001
Goal: <one sentence>
Depends on: none
Inputs:
- <REQ/AC/NFR/API/DATA IDs>
Outputs:
- <expected artifact/code outcome>
Tasks:
- TASK-001: <atomic task>
- TASK-002: <atomic task>
Validation:
- <command or manual check>
Done when:
- <binary condition>
Rollback/mitigation:
- <note>

## Phase 2: <name>

ID: PHASE-002
Goal: <one sentence>
Depends on: PHASE-001
Inputs:
- <IDs>
Outputs:
- <expected artifact/code outcome>
Tasks:
- TASK-003: <atomic task>
- TASK-004: <atomic task>
Validation:
- <command or manual check>
Done when:
- <binary condition>
Rollback/mitigation:
- <note>

## 15. Atomic task breakdown

| Task ID | Owner role | Depends on | Traceability | Objective | Output | Validation |
|---|---|---|---|---|---|---|
| TASK-001 | researcher/planner/implementer/tester | none | REQ-001 | <objective> | <output> | <validation> |

## 16. Validation gates

| Gate ID | When | Check | Pass condition | Fail action |
|---|---|---|---|---|
| GATE-001 | Before implementation | Plan completeness | No blocking questions remain | Return to clarification |
| GATE-002 | Before merge/release | Tests and review | All listed checks pass | Return to implementation |

## 17. Risks and mitigations

| Risk ID | Risk | Severity | Mitigation | Owner role |
|---|---|---|---|---|
| RISK-001 | <risk> | LOW/MEDIUM/HIGH | <mitigation> | <role> |

## 18. Handoff contract

### For implementer

Input artifacts:
- This implementation plan.
- Source PRD: <path>

Allowed scope:
- <scope>

Forbidden actions:
- Do not change unrelated files.
- Do not expand product scope.
- Do not skip validation gates.
- Do not implement behavior not mapped to a requirement, accepted assumption, risk mitigation, or validation task.

Expected output:
- Code changes.
- Test changes.
- Command outputs.
- Summary mapped to task IDs.

### For validator/tester

Input artifacts:
- This implementation plan.
- Implementation summary.
- Test outputs.

Expected output:
- Validation report mapped to TEST and GATE IDs.

### For reviewer

Input artifacts:
- This implementation plan.
- Diff.
- Test outputs.

Expected output:
- Review report with blocking and non-blocking findings mapped to requirement/task IDs.

## 19. Final readiness check

- [ ] Every requirement has implementation coverage.
- [ ] Every acceptance criterion has validation coverage.
- [ ] Every task has an owner role.
- [ ] Every task has explicit input and output.
- [ ] Every implementation phase has validation.
- [ ] Every risk has mitigation.
- [ ] Every assumption is human-accepted.
- [ ] No blocking question remains open.
- [ ] Rollout and rollback are defined.
- [ ] The plan does not require hidden context.
