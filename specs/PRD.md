# PRD - Vibe and Kalika Code

## 1. Product Name

Official product name:

```txt
Vibe and Kalika Code
```

Recommended binary name:

```txt
vkc
```

Reason: terminal commands with spaces are cursed. Do not do that.

---

## 2. Objective

Build an interactive CLI setup tool that installs AI coding agents and skills locally for supported coding-agent environments.

The CLI must help developers install reusable agent/skill packs safely, consistently and with clear visibility into what will be changed.

---

## 3. MVP Scope

The MVP installs only:

- Agents
- Skills

The MVP supports package sources from:

- Local directory
- Git repository

The MVP prioritizes:

| Platform | Priority |
|---|---|
| GitHub Copilot CLI | Must-have |
| Claude Code | Nice-to-have |
| OpenAI Codex CLI | Nice-to-have |

---

## 4. Non-goals

The MVP must not implement:

- Hooks
- Prompts
- Workflows
- MCP servers
- Memory systems
- Uninstall
- Update
- Custom registries
- Script execution
- Package marketplace
- Web UI

The architecture should allow these later without rewriting the core.

---

## 5. User Personas

### 5.1 Developer

A developer wants to install agents and skills into their local coding assistant environment without manually copying files into multiple hidden folders.

### 5.2 Tech Lead

A tech lead wants to distribute standardized agents and skills to a team with predictable installation behavior and minimal risk.

### 5.3 Agent Pack Maintainer

A maintainer wants to publish local or Git-based packages containing agents and skills.

---

## 6. Core User Stories

### 6.1 Install package interactively

As a developer, I want to run:

```bash
vkc init
```

Then select supported platforms, choose a local or Git package, review the install plan, and apply the setup safely.

### 6.2 Detect supported environments

As a developer, I want to run:

```bash
vkc detect
```

So I can see which AI coding environments are available on my machine.

### 6.3 Validate a package

As a maintainer, I want to run:

```bash
vkc validate ./my-pack
```

So I can verify the manifest, file layout and target compatibility before distribution.

### 6.4 Install from local path

As a developer, I want to run:

```bash
vkc install ./packs/kalika-reviewer
```

So I can install an agent pack from a local directory.

### 6.5 Install from Git

As a developer, I want to run:

```bash
vkc install https://github.com/org/agent-pack.git
```

So I can install a package from a Git repository.

---

## 7. Functional Requirements

### FR-001: Detect platforms

The CLI must detect:

- GitHub Copilot CLI home
- Claude Code project/user folders
- OpenAI Codex CLI user/project configuration

### FR-002: Interactive selection

The CLI must show checkboxes for detected and supported platforms.

The user must be able to install only into selected platforms.

### FR-003: Manifest validation

The CLI must parse and validate `manifest.yaml` before planning installation.

### FR-004: Dry-run

The CLI must always show a dry-run plan before applying changes.

### FR-005: Conflict handling

When a target already exists, the CLI must not overwrite without explicit confirmation.

Allowed actions:

- skip
- overwrite
- backup-and-overwrite
- write-as-new

Default action:

```txt
skip
```

### FR-006: Backups

The CLI must create backups in:

```txt
.ai-setup/backups/<timestamp>/
```

Before modifying or overwriting existing files.

### FR-007: Reports

The CLI must generate installation reports in:

```txt
.ai-setup/reports/<timestamp>-install-report.md
```

### FR-008: Installation state

The CLI must maintain state in:

```txt
.ai-setup/installed.yaml
```

### FR-009: Security

The CLI must:

- reject path traversal
- reject unknown targets
- avoid writing outside allowed target directories
- not execute package scripts
- warn about script files
- resolve symlinks safely
- require confirmation for suspicious files

---

## 8. Quality Requirements

The implementation must prioritize:

- Safety
- Idempotency
- Testability
- Clear terminal UX
- Minimal hidden behavior
- Idiomatic Go
- Clean adapter-based platform separation

---

## 9. Success Criteria

The MVP is complete when:

- `vkc init` works interactively
- `vkc detect` shows supported platform state
- `vkc install <local-path>` installs a valid package
- `vkc install <git-url>` installs a valid Git package
- Copilot CLI agents and skills install correctly
- Dry-run happens before changes
- Backups are created when needed
- Reports are generated
- `.ai-setup/installed.yaml` is updated
- Unit and integration tests cover critical behavior
