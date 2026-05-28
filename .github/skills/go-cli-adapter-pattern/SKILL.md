---
name: go-cli-adapter-pattern
description: Use this skill when building a Go CLI that must install, configure, or operate across multiple target platforms (e.g., Copilot CLI, Claude Code, Codex CLI). Call it before implementing any platform-specific dispatch logic. It produces a clean interface-per-platform design where generic orchestration code never imports platform packages directly.
---

# Go CLI Multi-Platform Adapter Pattern

## Purpose

Define a `PlatformAdapter` interface in a neutral `platform` package and implement
each platform in its own sub-package (`platform/copilot`, `platform/claude`,
`platform/codex`). Generic orchestration code depends only on the interface.
Each adapter file ends with a compile-time interface check.

## Trigger

Call this skill when:

- A Go CLI must perform the same logical operation (detect, plan, install) on ≥ 2
  different target platforms.
- You need to add a new platform without touching existing platform code.
- A code review asks why the orchestrator imports a platform-specific package.
- You want to prevent switch-statements scattered across business logic.

## Do not use when

Do not call this skill when:

- The project targets exactly one platform and there is no roadmap to add more.
- The platforms differ so fundamentally that a shared interface would be empty.
- The "platforms" are actually environment tiers (dev/staging/prod), not tooling targets.

## Required inputs

The agent must provide:

- List of platforms to support (names + their config paths / env vars).
- The shared operations each platform must implement.
- The project Go module path.

## Procedure

### 1 — Define the neutral `platform` package

```go
// internal/platform/platform.go
package platform

type Platform string

const (
    PlatformCopilotCLI Platform = "copilot-cli"
    PlatformClaudeCode Platform = "claude-code"
    PlatformCodexCLI   Platform = "codex-cli"
)

// DetectionResult is returned by Detect() for every adapter.
type DetectionResult struct {
    Detected   bool
    Platform   Platform
    BasePath   string
    AgentsPath string
    SkillsPath string
    Notes      []string
}
```

### 2 — Define the `PlatformAdapter` interface

```go
// internal/platform/adapter.go
package platform

type PlatformAdapter interface {
    Platform() Platform
    Detect(projectRoot string) DetectionResult
    Plan(input PlanInput) ([]PlannedOperation, error)
    Validate(input ValidateInput) error
    AllowedRoots(projectRoot string) []string
}
```

### 3 — Implement each platform in its own sub-package

```go
// internal/platform/copilot/adapter.go
package copilot

import "github.com/org/repo/internal/platform"

type Adapter struct{}

func NewAdapter() *Adapter { return &Adapter{} }

func (a *Adapter) Platform() platform.Platform { return platform.PlatformCopilotCLI }

func BasePath() string {
    if v := os.Getenv("COPILOT_HOME"); v != "" {
        return v
    }
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".copilot")
}

func (a *Adapter) Detect(projectRoot string) platform.DetectionResult { ... }
func (a *Adapter) Plan(input platform.PlanInput) ([]platform.PlannedOperation, error) { ... }
func (a *Adapter) Validate(input platform.ValidateInput) error { return nil }
func (a *Adapter) AllowedRoots(projectRoot string) []string { return []string{BasePath()} }

// Compile-time interface check — fails at build if Adapter drifts from interface.
var _ platform.PlatformAdapter = (*Adapter)(nil)
```

Repeat the same structure for `platform/claude` and `platform/codex`.

### 4 — Register adapters in the app layer, never in generic code

```go
// internal/app/install.go
func buildRegistry() []platform.PlatformAdapter {
    return []platform.PlatformAdapter{
        copilot.NewAdapter(),
        claude.NewAdapter(),
        codex.NewAdapter(),
    }
}
```

### 5 — Override install paths via env vars in tests

```go
func TestInstallCopilotAndClaude(t *testing.T) {
    copilotHome := filepath.Join(t.TempDir(), ".copilot")
    t.Setenv("COPILOT_HOME", copilotHome)
    // ... run app.Install(...) and assert on copilotHome
}
```

Each adapter reads its env override in `BasePath()`, so `t.Setenv` redirects
all writes to `t.TempDir()` without patching any global state.

## Key rules

- The `platform` package must **never** import any `platform/<name>` sub-package.
- Each `adapter.go` must end with `var _ platform.PlatformAdapter = (*Adapter)(nil)`.
- `BasePath()` must check an env var first so tests can redirect writes.
- `AllowedRoots()` must return the exact roots the security layer will enforce — never the full home dir.
- Generic orchestration code (app, installer, planner) imports only `internal/platform`, never sub-packages.

## Expected output

```md
# Go CLI Adapter Pattern Result

## Status
PASS | BLOCKED | NOT_APPLICABLE

## Findings
- Interface defined in neutral package: yes/no
- Each platform in its own sub-package: yes/no
- Compile-time checks in place: yes/no
- Env-var overrides for tests: yes/no

## Next step
- Proceed to implement platform operations
```

## Stop conditions

Stop and return `BLOCKED` when:

- The shared interface would require more than ~6 methods; consider splitting.
- A platform-specific type must leak into the shared interface (redesign instead).
- Import cycles form between `platform` and any sub-package (fix package boundaries first).
