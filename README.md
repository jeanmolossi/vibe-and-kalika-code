# vkc — Vibe and Kalika Code

`vkc` is a CLI tool that installs AI coding agents and skills into supported AI coding environments.

Supported platforms: **GitHub Copilot CLI**, **Claude Code**, **OpenAI Codex CLI**.

---

## Installation

**Build from source (recommended):**

```bash
git clone https://github.com/jeanmolossi/vibe-and-kalika-code.git
cd vibe-and-kalika-code
go build -o vkc ./cmd/vkc
# optionally move to PATH
mv vkc /usr/local/bin/vkc
```

**Or install directly with `go install`:**

```bash
go install github.com/jeanmolossi/vibe-and-kalika-code/cmd/vkc@latest
```

Requires **Go 1.23+**.

---

## Quick Start

```bash
# 1. Check what platforms are detected in the current project
vkc detect

# 2. Validate a package before installing
vkc validate ./examples/basic-pack

# 3. Preview what would be installed (no filesystem changes)
vkc install ./examples/basic-pack --dry-run

# 4. Install
vkc install ./examples/basic-pack

# 5. Check environment health
vkc doctor
```

---

## Commands

### `vkc detect`

Detects supported AI coding platforms in the current directory and user home.

```bash
vkc detect
```

Output example:

```
copilot-cli detected=true  base=/home/user/.copilot
claude-code  detected=true  base=/home/user/project/.claude
codex-cli    detected=false base=
```

No filesystem writes.

---

### `vkc validate <source>`

Validates a package without installing it.

```bash
vkc validate ./my-pack
vkc validate https://github.com/org/agent-pack.git
```

Checks performed:
- `manifest.yaml` exists and all required fields are present
- All agent and skill source files exist
- Each skill directory contains a `SKILL.md`
- No path traversal in source paths
- All declared targets are supported platforms
- No duplicate target operations

Exits `0` on success, `2` on validation error.

---

### `vkc install <source>`

Installs a package from a local directory or Git repository.

```bash
vkc install ./my-pack
vkc install https://github.com/org/agent-pack.git
vkc install git@github.com:org/agent-pack.git
```

**Flags:**

| Flag | Default | Description |
|---|---|---|
| `--dry-run` | false | Show plan without writing any files |
| `--yes` | false | Skip confirmation prompts |
| `--targets` | (all) | Comma-separated platform targets, e.g. `copilot-cli,claude-code` |
| `--conflict-action` | `skip` | How to handle conflicts: `skip`, `overwrite`, `backup-and-overwrite` |

A dry-run plan is always shown before any changes are applied, even without `--dry-run`.

---

### `vkc init [source]`

Interactive setup wizard. Guides through platform detection, source selection, conflict handling, and installation.

```bash
vkc init
vkc init ./my-pack
vkc init --yes ./my-pack
```

> **Note:** The interactive wizard is currently a placeholder. For reliable non-interactive use, prefer `vkc install <source>`.

---

### `vkc doctor`

Reports environment health.

```bash
vkc doctor
```

Checks:
- Project root detection
- Detected platforms and their paths
- Git availability
- Write permissions to target directories
- `.ai-setup/installed.yaml` state validity

No filesystem writes.

---

## Package Format

A `vkc`-installable package is a directory containing a `manifest.yaml` and source files for agents and skills.

### Directory Layout

```
my-pack/
├── manifest.yaml
├── agents/
│   └── my-agent.md
└── skills/
    └── my-skill/
        ├── SKILL.md
        └── (any supporting files)
```

### `manifest.yaml`

```yaml
name: my-pack
version: 0.1.0
description: What this pack does.
author: Your Name

targets:
  - copilot-cli
  - claude-code
  - codex-cli

agents:
  - name: my-agent
    description: What this agent does.
    source: agents/my-agent.md
    targets:
      copilot-cli:
        scope: user        # user (~/.copilot) or project (./.github)
      claude-code:
        scope: project     # project (./.claude)
      codex-cli:
        mode: agents-md-section  # writes a managed block in AGENTS.md

skills:
  - name: my-skill
    description: What this skill does.
    source: skills/my-skill
    targets:
      copilot-cli:
        scope: user
      claude-code:
        scope: project
      codex-cli:
        scope: project
```

See `examples/basic-pack/` for a working reference package.

---

## Platform Support

| Platform | ID | Agents path | Skills path | Status |
|---|---|---|---|---|
| GitHub Copilot CLI | `copilot-cli` | `~/.copilot/agents/` (user) | `~/.copilot/skills/` (user) | ✅ Supported |
| Claude Code | `claude-code` | `./.claude/agents/` (project) | `./.claude/skills/` (project) | ✅ Supported |
| OpenAI Codex CLI | `codex-cli` | `AGENTS.md` (managed section) | `./.agents/skills/` (project) | ✅ Supported |

**`copilot-cli`** — Writes to `~/.copilot/` by default (user scope). Project scope (`.github/`) is planned but not active.

**`claude-code`** — Writes to `./.claude/` (project scope). User scope (`~/.claude/`) is detected but not installed to by default.

**`codex-cli`** — Agents are injected as a managed section in `AGENTS.md`. Skills are written to `./.agents/skills/`.

---

## Configuration

`vkc` uses the current working directory as the project root. No config file is needed.

**Environment variables:**

| Variable | Description |
|---|---|
| `COPILOT_HOME` | Override the Copilot CLI home directory (default: `~/.copilot`) |

---

## Exit Codes

| Code | Meaning |
|---:|---|
| 0 | Success |
| 1 | Generic error |
| 2 | Validation error |
| 3 | Security violation (path traversal, disallowed root) |
| 4 | User cancelled |
| 5 | Source fetch error (network, missing repo) |
| 6 | Conflict unresolved |

---

## Known Limitations

- `vkc init` interactive wizard is a placeholder — use `vkc install <source>` for non-interactive installs.
- Git sources require network connectivity.
- `vkc uninstall` is not implemented. Installation state is tracked in `.ai-setup/installed.yaml` for future use.
- `vkc update` is not implemented.
- Claude Code user scope (`~/.claude/`) is detected but not installed to by default.

---

## Development

```bash
# Run all tests
go test ./...

# Run tests with race detector
go test -race ./...

# Lint and vet
go vet ./...

# Build binary
go build -o vkc ./cmd/vkc

# Run integration tests only
go test ./internal/integration/...
```

---

## Spec Files

These files describe the design intent and are kept as-is for reference:

- `PRD.md` — Product requirements
- `ARCHITECTURE.md` — Technical architecture
- `IMPLEMENTATION_PLAN.md` — Execution plan
- `CLI_SPEC.md` — CLI commands and UX
- `MANIFEST_SPEC.md` — Package manifest contract
- `FILE_MAPPING_SPEC.md` — Target filesystem mapping per platform
- `VALIDATION_AND_ACCEPTANCE.md` — Validation and acceptance criteria
- `PROMPT_FOR_LLM_WORKER.md` — LLM implementation prompt
- `examples/basic-pack/` — Reference installable package
