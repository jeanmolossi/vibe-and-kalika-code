---
name: go-cobra-exit-code-propagation
description: Use this skill when a Go CLI built with Cobra must return distinct exit codes (e.g., 0 success, 1 generic error, 2 validation failure, 3 conflict) instead of always exiting 1. Call it before implementing error handling in any Cobra command. It produces a cliError sentinel type, an ExitCodeFrom() helper, and a main() that calls os.Exit() only once.
---

# Go Cobra Exit Code Propagation

## Purpose

Cobra's `Execute()` returns an `error`. By default, any non-nil error causes
`os.Exit(1)`. This skill wraps errors with a typed `cliError` struct that
carries an exit code, and provides `ExitCodeFrom()` to extract it cleanly.
`main()` remains the only place that calls `os.Exit()`.

## Trigger

Call this skill when:

- A Cobra CLI must return different exit codes per failure class.
- CI/CD pipelines or shell scripts branch on the CLI exit code.
- A command must propagate exit codes from sub-processes or domain errors.
- A code review flags multiple `os.Exit()` calls scattered across command files.

## Do not use when

Do not call this skill when:

- The CLI only needs success (0) / failure (1) — standard Cobra is sufficient.
- The CLI is a short script with no downstream exit-code consumers.

## Required inputs

The agent must provide:

- The exit code table (code → meaning) for the specific CLI.
- The Go module path.

## Procedure

### 1 — Define `cliError` in `internal/cli/error.go`

```go
package cli

import "fmt"

type cliError struct {
    code int
    err  error
}

func (e cliError) Error() string { return e.err.Error() }
func (e cliError) Unwrap() error { return e.err }
func (e cliError) ExitCode() int { return e.code }

// exitError wraps err with the given exit code.
func exitError(code int, err error) error {
    return cliError{code: code, err: fmt.Errorf("%w", err)}
}

// ExitCodeFrom extracts the exit code from a cliError; returns 1 for any
// other error type.
func ExitCodeFrom(err error) int {
    type coded interface{ ExitCode() int }
    if as, ok := err.(coded); ok {
        return as.ExitCode()
    }
    return 1
}
```

### 2 — Define exit code constants (same package or a sibling file)

```go
// internal/cli/codes.go  (or internal/app/codes.go)
const (
    ExitSuccess          = 0
    ExitGenericError     = 1
    ExitValidationError  = 2
    ExitConflict         = 3
    ExitNotFound         = 4
    ExitPermission       = 5
    ExitUnsupportedPlatform = 6
)
```

### 3 — Return `exitError` from Cobra `RunE` functions

```go
// internal/cli/install_cmd.go
func runInstall(cmd *cobra.Command, args []string) error {
    _, code, err := app.Install(opts)
    if err != nil {
        return exitError(code, err)
    }
    return nil
}
```

### 4 — Keep `main()` minimal — one call, one `os.Exit`

```go
// cmd/<name>/main.go
package main

import (
    "fmt"
    "os"

    "github.com/org/repo/internal/cli"
)

func main() {
    if err := cli.NewRootCmd().Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(cli.ExitCodeFrom(err))
    }
}
```

### 5 — Suppress Cobra's own error printing to avoid double output

```go
// internal/cli/root.go
func NewRootCmd() *cobra.Command {
    root := &cobra.Command{
        Use:           "vkc",
        SilenceErrors: true, // main() prints the error
        SilenceUsage:  true, // suppress usage on runtime errors
    }
    return root
}
```

## Key rules

- `os.Exit()` must appear **only** in `main()`.
- Every non-zero exit must go through `exitError(code, err)`.
- `ExitCodeFrom` must return `1` as the default fallback, never `0`.
- `SilenceErrors: true` on the root command prevents Cobra printing errors twice.
- Tests exercise `ExitCodeFrom` directly — never call `os.Exit` from tests.

## Expected output

```md
# Cobra Exit Code Propagation Result

## Status
PASS | BLOCKED | NOT_APPLICABLE

## Findings
- cliError type defined: yes/no
- ExitCodeFrom helper present: yes/no
- os.Exit only in main(): yes/no
- SilenceErrors set on root command: yes/no

## Next step
- Map domain errors to exit codes and wire through RunE
```

## Stop conditions

Stop and return `BLOCKED` when:

- A third-party library calls `os.Exit()` internally and cannot be replaced.
- Exit codes must cross process boundaries via IPC (different pattern needed).
