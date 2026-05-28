#!/usr/bin/env bash
# Security gate — runs on agentStop.
# Blocks the agent if vulnerabilities are found in Go code paths after dependency changes.
set -euo pipefail

ROOT=$(git rev-parse --show-toplevel 2>/dev/null || pwd)
cd "$ROOT"

# Only run when go.mod, go.sum, or Go source files were modified.
CHANGED=$(git status --porcelain 2>/dev/null | awk '{print $2}' | grep -E '\.(go|mod|sum)$' || true)
if [ -z "$CHANGED" ]; then
  exit 0
fi

ERRORS=""

# 1. govulncheck (vulnerabilities in code paths)
if command -v govulncheck >/dev/null 2>&1; then
  VULN_OUT=$(govulncheck ./... 2>&1 || true)
  if echo "$VULN_OUT" | grep -q "Your code is affected by"; then
    ERRORS="${ERRORS}\n• govulncheck: vulnerabilidades encontradas nos code paths (execute: govulncheck ./... para ver detalhes; corrija com go get <pkg>@latest && go mod tidy)"
  fi
else
  ERRORS="${ERRORS}\n• govulncheck: não instalado. Instale com: go install golang.org/x/vuln/cmd/govulncheck@latest"
fi

# 2. go mod tidy check (ensures go.sum is clean after changes)
if echo "$CHANGED" | grep -qE '\.(mod|sum)$'; then
  cp go.sum go.sum.bak 2>/dev/null || true
  go mod tidy >/dev/null 2>&1 || true
  if ! diff -q go.sum go.sum.bak >/dev/null 2>&1; then
    ERRORS="${ERRORS}\n• go mod tidy: go.sum diverge após tidy (execute: go mod tidy)"
  fi
  mv go.sum.bak go.sum 2>/dev/null || true
fi

if [ -n "$ERRORS" ]; then
  printf '{"decision":"block","reason":"Security gate falhou. Corrija antes de finalizar:%s"}' "$ERRORS"
fi
