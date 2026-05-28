# Architecture - Vibe and Kalika Code

## 1. Overview

Vibe and Kalika Code is a Go CLI that installs agent and skill packages into supported AI coding environments.

The architecture must be adapter-based:

```txt
CLI command
  -> Application service
    -> Source provider
    -> Manifest parser/validator
    -> Platform detector
    -> Platform adapters
    -> Planner
    -> Security validator
    -> Installer
    -> Backup manager
    -> State store
    -> Report generator
```

No platform-specific logic should leak into generic installation logic.

---

## 2. Suggested Project Structure

```txt
vibe-and-kalika-code/
  cmd/
    vkc/
      main.go

  internal/
    app/
      init.go
      install.go
      detect.go
      validate.go
      doctor.go

    cli/
      root.go
      init_cmd.go
      install_cmd.go
      detect_cmd.go
      validate_cmd.go
      doctor_cmd.go

    ui/
      prompts.go
      styles.go
      tables.go

    manifest/
      manifest.go
      parser.go
      validator.go

    source/
      source.go
      local.go
      git.go

    platform/
      platform.go
      detector.go
      adapter.go
      copilot/
        adapter.go
        detector.go
      claude/
        adapter.go
        detector.go
      codex/
        adapter.go
        detector.go

    planner/
      plan.go
      builder.go
      conflict.go

    installer/
      installer.go
      copy.go
      merge.go
      permissions.go

    backup/
      backup.go
      restore.go

    state/
      installed.go
      store.go

    report/
      report.go
      markdown.go

    security/
      paths.go
      scripts.go
      symlinks.go

  testdata/
    packages/
      valid-basic/
      with-conflicts/
      with-scripts/
      invalid-path-traversal/

  go.mod
  go.sum
  README.md
```

---

## 3. Core Domain Types

### 3.1 Platform

```go
type Platform string

const (
    PlatformCopilotCLI Platform = "copilot-cli"
    PlatformClaudeCode Platform = "claude-code"
    PlatformCodexCLI   Platform = "codex-cli"
)
```

### 3.2 PlatformAdapter

```go
type PlatformAdapter interface {
    Platform() Platform
    Detect(ctx context.Context, projectRoot string) (DetectionResult, error)
    Plan(ctx context.Context, input PlanInput) ([]PlannedOperation, error)
    Validate(ctx context.Context, input ValidateInput) error
}
```

### 3.3 SourceProvider

```go
type SourceProvider interface {
    Fetch(ctx context.Context, source string) (PackageSource, error)
}
```

### 3.4 PlannedOperation

```go
type OperationType string

const (
    OperationCreate OperationType = "create"
    OperationModify OperationType = "modify"
    OperationSkip   OperationType = "skip"
)

type PlannedOperation struct {
    Type        OperationType
    Platform    Platform
    SourcePath  string
    TargetPath  string
    Description string
    Conflict    *Conflict
    Warnings    []Warning
}
```

---

## 4. Core Components

## 4.1 CLI Layer

Responsibilities:

- Parse command arguments
- Call application services
- Render output
- Trigger interactive forms

Recommended library:

- `github.com/spf13/cobra`

Commands:

```bash
vkc init
vkc detect
vkc install <source>
vkc validate
vkc doctor
```

---

## 4.2 UI Layer

Responsibilities:

- Welcome screens
- Checkbox selection
- Confirmation screens
- Conflict prompts
- Dry-run plan rendering
- Final summaries

Recommended libraries:

- `github.com/charmbracelet/huh`
- `github.com/charmbracelet/lipgloss`

---

## 4.3 Manifest Module

Responsibilities:

- Parse `manifest.yaml`
- Validate required fields
- Validate target compatibility
- Validate source paths
- Validate skill folder shape
- Validate agent file shape

Must reject:

- missing manifest
- missing package name
- unknown target
- path traversal
- skill without `SKILL.md`
- duplicate final target path

---

## 4.4 Source Module

### Local source

Reads a local directory.

### Git source

Clones a Git repository into a temporary directory.

Rules:

- Do not execute scripts.
- Validate manifest after cloning.
- Remove temp directory after installation.
- Version pinning is nice-to-have.
- Checksums are nice-to-have.

---

## 4.5 Platform Module

Each platform has:

- Detector
- Adapter
- Mapping rules

Adapters must produce planned operations, not directly write files.

---

## 4.6 Planner Module

Responsibilities:

- Convert manifest + selected platforms into planned operations
- Detect conflicts
- Detect duplicates
- Attach security warnings
- Produce dry-run data

The planner must not mutate filesystem state.

---

## 4.7 Installer Module

Responsibilities:

- Apply planned operations
- Copy files/directories
- Merge managed blocks in `AGENTS.md`
- Respect conflict decisions
- Preserve file permissions when safe
- Refuse unsafe writes

---

## 4.8 Backup Module

Responsibilities:

- Create backup directory
- Copy pre-existing target files before modification
- Write `backup-report.yaml`

Backup path:

```txt
.ai-setup/backups/<timestamp>/
```

---

## 4.9 State Module

Responsibilities:

- Read `.ai-setup/installed.yaml`
- Append installation records
- Preserve existing state
- Support future update/uninstall

---

## 4.10 Report Module

Responsibilities:

- Generate Markdown installation report
- Include created, modified, skipped files
- Include conflicts
- Include warnings
- Include backup path
- Include selected platforms

Report path:

```txt
.ai-setup/reports/<timestamp>-install-report.md
```

---

## 4.11 Security Module

Responsibilities:

- Path traversal prevention
- Symlink safety
- Allowed-target validation
- Script detection
- Binary detection if practical

Rules:

- Never write outside approved target roots.
- Never execute package files.
- Warn about scripts.
- Require confirmation for scripts.
- Resolve symlinks before copying or writing.
