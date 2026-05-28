#!/usr/bin/env bash
set -euo pipefail

PLAN_FILE="${1:-}"

if [[ -z "$PLAN_FILE" ]]; then
  echo "Usage: $0 <plan.md>" >&2
  exit 2
fi

if [[ ! -f "$PLAN_FILE" ]]; then
  echo "ERROR: plan file not found: $PLAN_FILE" >&2
  exit 2
fi

required_sections=(
  "# Implementation Plan:"
  "## 1. Context summary"
  "## 2. Source requirement map"
  "## 3. Accepted assumptions"
  "## 4. Out of scope"
  "## 5. Architecture impact"
  "## 6. Data model impact"
  "## 7. API and integration impact"
  "## 8. Frontend/UI impact"
  "## 9. Backend/domain impact"
  "## 10. Security and compliance"
  "## 11. Observability"
  "## 12. Rollout strategy"
  "## 13. Test strategy"
  "## 14. Implementation phases"
  "## 15. Atomic task breakdown"
  "## 16. Validation gates"
  "## 17. Risks and mitigations"
  "## 18. Handoff contract"
  "## 19. Final readiness check"
)

failures=0

for section in "${required_sections[@]}"; do
  if ! grep -Fq "$section" "$PLAN_FILE"; then
    echo "ERROR: missing section: $section" >&2
    failures=$((failures + 1))
  fi
done

required_patterns=(
  "REQ-[0-9][0-9][0-9]"
  "AC-[0-9][0-9][0-9]"
  "TASK-[0-9][0-9][0-9]"
  "PHASE-[0-9][0-9][0-9]"
  "GATE-[0-9][0-9][0-9]"
  "RISK-[0-9][0-9][0-9]"
)

for pattern in "${required_patterns[@]}"; do
  if ! grep -Eq "$pattern" "$PLAN_FILE"; then
    echo "ERROR: missing required ID pattern: $pattern" >&2
    failures=$((failures + 1))
  fi
done

if grep -Eiq "implement backend|adjust frontend|fix tests|review everything|handle edge cases|improve performance" "$PLAN_FILE"; then
  echo "WARNING: plan may contain vague task wording. Replace vague tasks with atomic tasks." >&2
fi

if [[ "$failures" -gt 0 ]]; then
  echo "Plan validation failed with $failures error(s)." >&2
  exit 1
fi

echo "Plan validation passed. Humanity avoided one tiny procedural disaster."
