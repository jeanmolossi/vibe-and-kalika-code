#!/usr/bin/env bash
# Quality gate — runs on agentStop.
# Blocks the agent if Go quality checks fail after any Go file was modified.
set -euo pipefail

ROOT=$(git rev-parse --show-toplevel 2>/dev/null || pwd)
cd "$ROOT"

# Only run when Go files were modified.
CHANGED=$(git status --porcelain 2>/dev/null | awk '{print $2}' | grep '\.go$' || true)
if [ -z "$CHANGED" ]; then
  exit 0
fi

ERRORS=""

# 1. Format check
FMT=$(gofmt -l . 2>/dev/null || true)
if [ -n "$FMT" ]; then
  ERRORS="${ERRORS}\n• gofmt: arquivos sem formatação (execute: gofmt -w .): ${FMT}"
fi

# 2. Build
if ! go build ./... >/dev/null 2>&1; then
  ERRORS="${ERRORS}\n• go build: falhou (execute: go build ./... para ver os erros)"
fi

# 3. Vet
if ! go vet ./... >/dev/null 2>&1; then
  ERRORS="${ERRORS}\n• go vet: falhou (execute: go vet ./... para ver os erros)"
fi

# 4. Lint (optional — only if installed)
if command -v golangci-lint >/dev/null 2>&1; then
  if ! golangci-lint run ./... >/dev/null 2>&1; then
    ERRORS="${ERRORS}\n• golangci-lint: issues encontrados (execute: golangci-lint run ./... para ver detalhes)"
  fi
fi

# 5. Tests
if ! go test ./... >/dev/null 2>&1; then
  ERRORS="${ERRORS}\n• go test: testes falharam (execute: go test ./... para ver detalhes)"
fi

if [ -n "$ERRORS" ]; then
  printf '{"decision":"block","reason":"Quality gate falhou. Corrija antes de finalizar:%s"}' "$ERRORS"
fi
