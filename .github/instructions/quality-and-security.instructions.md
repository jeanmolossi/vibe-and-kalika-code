---
applyTo: "**/*.go"
---

# Go Code Quality & Security Gates

You are working in a Go project (`github.com/jeanmolossi/vibe-and-kalika-code`).
After every code change you make, you **must** run the following checks in order and fix all failures before considering the task done.

## 1. After implementing any Go code

Run these commands **in order**. Stop and fix failures before continuing.

```bash
# 1. Auto-format (always run, non-blocking)
gofmt -w .
goimports -w -local github.com/jeanmolossi/vibe-and-kalika-code .

# 2. Build (must pass)
go build ./...

# 3. Vet (must pass)
go vet ./...

# 4. Lint (must pass — zero issues)
golangci-lint run ./...

# 5. Tests (must pass — all green)
go test -race -count=1 ./...

# 6. Vulnerability scan (must pass — zero vulnerabilities in code paths)
govulncheck ./...
```

## 2. After changing `go.mod` or `go.sum`

```bash
go mod tidy
govulncheck ./...
go build ./...
go test ./...
```

## 3. Lint rules

- Never use `//nolint` without an explanation comment: `//nolint:linter // reason here`
- Never suppress `gosec` or `errcheck` globally
- Config file: `.golangci.yml` at repo root

## 4. Security rules

- Do not commit secrets, tokens, or credentials — use environment variables
- `GITHUB_TOKEN` is the only token available in CI (no extra tokens configured)
- If `govulncheck` reports a vulnerability, upgrade the dependency: `go get <pkg>@latest && go mod tidy`
- Keep `toolchain` in `go.mod` pinned to the latest Go patch to avoid stdlib CVEs

## 5. Dependency rules

- Prefer stdlib over external dependencies
- Always run `go mod tidy` after adding or removing a dependency
- Always run `govulncheck ./...` after any dependency change

## 6. Deferred error handling

Deferred calls that return errors (e.g., `defer f.Close()`) must either handle the error or use:
```go
defer f.Close() //nolint:errcheck // close error on read-only file is non-actionable
```
