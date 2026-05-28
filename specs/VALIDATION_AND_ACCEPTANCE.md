# Validation and Acceptance - Vibe and Kalika Code

## 1. Validation Rules

## 1.1 Package validation

A package is valid only if:

- root contains `manifest.yaml`
- manifest has required fields
- all declared source paths exist
- all agent sources are Markdown files
- all skill sources are directories containing `SKILL.md`
- targets are supported
- no source path escapes package root
- no duplicate target path exists in plan

---

## 1.2 Security validation

The CLI must reject:

- `../`
- absolute malicious source paths
- symlinks escaping package root
- target paths outside allowed roots
- unsupported target names

The CLI must warn about:

- `.sh`
- `.ps1`
- `.bat`
- `.cmd`
- executable files
- likely binary files

The CLI must never execute package files.

---

## 1.3 Platform validation

### Copilot CLI

Must validate:

- base path resolved from `$COPILOT_HOME` or `~/.copilot`
- agents directory path
- skills directory path
- write permission where needed

### Claude Code

Must validate:

- project root
- `.claude/agents`
- `.claude/skills`

### Codex CLI

Must validate:

- project root
- `AGENTS.md` merge target
- `.agents/skills`

---

## 2. Acceptance Criteria

## 2.1 CLI commands

The following must work:

```bash
vkc --help
vkc init
vkc detect
vkc install <source>
vkc validate <source>
vkc doctor
```

---

## 2.2 Copilot CLI support

The CLI must:

- install agents into `~/.copilot/agents`
- install skills into `~/.copilot/skills`
- respect `$COPILOT_HOME`
- create directories if missing
- not overwrite without confirmation
- register state
- generate report

---

## 2.3 Local source install

Given a valid package directory:

```bash
vkc install ./testdata/packages/valid-basic
```

The CLI must:

- parse manifest
- validate package
- build dry-run plan
- ask confirmation
- copy agent and skill files
- write state
- generate report

---

## 2.4 Git source install

Given a valid Git repository URL:

```bash
vkc install https://github.com/org/agent-pack.git
```

The CLI must:

- clone into temp directory
- validate package
- install from cloned content
- cleanup temp directory
- not execute files from repo

---

## 2.5 Conflict handling

When target exists:

- default action is skip
- overwrite requires explicit confirmation
- backup-and-overwrite creates backup first
- write-as-new creates a non-conflicting filename

---

## 2.6 Dry-run

Dry-run must:

- happen before apply
- list creates
- list modifies
- list conflicts
- list warnings
- not mutate filesystem

---

## 2.7 Reports

After successful installation, report must exist:

```txt
.ai-setup/reports/<timestamp>-install-report.md
```

It must include:

- package name
- version
- source
- selected platforms
- created files
- modified files
- skipped files
- conflicts
- warnings
- backup path

---

## 2.8 State

After successful installation, state must exist:

```txt
.ai-setup/installed.yaml
```

It must include:

- package
- version
- source
- installed_at
- platforms
- files
- managed markers for partial files like `AGENTS.md`

---

## 3. Required Tests

Unit tests:

- manifest parser
- manifest validator
- source path validation
- path traversal rejection
- symlink rejection
- script detection
- Copilot mapping
- Claude mapping
- Codex mapping
- conflict detection
- managed block merge
- backup report generation
- install report generation
- installed state write/read

Integration tests:

- install local package into fake HOME
- install Git package from local Git repo
- install only Copilot
- install Copilot + Claude
- install Codex agent into existing `AGENTS.md`
- conflict skip
- conflict backup-and-overwrite

All tests must use `t.TempDir()` and must not write into the real user home.

---

## 4. Final Quality Gate

Before final delivery:

```bash
go test ./...
go vet ./...
go build ./cmd/vkc
```

No fake success. No "it should work". Evidence or shame.
