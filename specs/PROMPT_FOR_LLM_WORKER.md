## Status

> **Implementation complete.** All commands (`vkc init`, `vkc detect`, `vkc install`, `vkc validate`, `vkc doctor`) are implemented, tested, and passing CI. This document is kept as a historical reference for the implementation prompt. For current usage, see `README.md`.

---

# Prompt for LLM Worker - Implement Vibe and Kalika Code

You are a coding worker responsible for implementing the **Vibe and Kalika Code** CLI.

The binary name is:

```txt
vkc
```

The CLI must be written in Go.

Read and follow these project artifacts:

- `PRD.md`
- `ARCHITECTURE.md`
- `IMPLEMENTATION_PLAN.md`
- `CLI_SPEC.md`
- `MANIFEST_SPEC.md`
- `FILE_MAPPING_SPEC.md`
- `VALIDATION_AND_ACCEPTANCE.md`

## Mission

Implement an interactive CLI setup tool that installs AI coding agents and skills locally.

The MVP supports:

- agents
- skills
- local directory packages
- Git repository packages
- GitHub Copilot CLI as must-have
- Claude Code as nice-to-have
- OpenAI Codex CLI as nice-to-have

## Hard Rules

1. Use Go.
2. Use mature helper libraries when justified.
3. Keep architecture adapter-based.
4. Do not hardcode all platforms in one giant function.
5. Do not implement out-of-scope features.
6. Do not execute package scripts.
7. Do not overwrite files without explicit confirmation.
8. Always generate a dry-run plan before applying changes.
9. Always generate backups before overwriting.
10. Always write an installation report.
11. Always maintain `.ai-setup/installed.yaml`.
12. Reject path traversal.
13. Reject unsafe symlinks.
14. Reject writes outside allowed target roots.
15. Write tests for core behavior.
16. Keep implementation idiomatic, simple and readable.

## MVP Commands

Implement:

```bash
vkc init
vkc detect
vkc install <source>
vkc validate <source>
vkc doctor
```

## Suggested Libraries

You may use:

- `github.com/spf13/cobra`
- `github.com/charmbracelet/huh`
- `github.com/charmbracelet/lipgloss`
- `gopkg.in/yaml.v3`
- `github.com/go-git/go-git/v5`

Use only what is necessary.

## Implementation Order

1. Create Go module and CLI skeleton.
2. Add Cobra root command.
3. Implement manifest structs.
4. Implement manifest parser.
5. Implement manifest validator.
6. Implement local source provider.
7. Implement Git source provider.
8. Implement platform detector registry.
9. Implement Copilot detector and adapter.
10. Implement planner.
11. Implement dry-run renderer.
12. Implement conflict detector.
13. Implement security path checks.
14. Implement installer.
15. Implement backup manager.
16. Implement report generator.
17. Implement state store.
18. Implement interactive `init`.
19. Implement `detect`.
20. Implement `validate`.
21. Implement `doctor`.
22. Implement Claude adapter.
23. Implement Codex adapter.
24. Add tests.
25. Add README.

## Critical Platform Mapping

### Copilot CLI

Use:

```txt
$COPILOT_HOME
```

or fallback:

```txt
~/.copilot
```

Agents:

```txt
<copilot_home>/agents/<agent-name>.md
```

Skills:

```txt
<copilot_home>/skills/<skill-name>/SKILL.md
```

### Claude Code

Project agents:

```txt
<project_root>/.claude/agents/<agent-name>.md
```

Project skills:

```txt
<project_root>/.claude/skills/<skill-name>/SKILL.md
```

### Codex CLI

Project skills:

```txt
<project_root>/.agents/skills/<skill-name>/SKILL.md
```

Agents must be merged into:

```txt
<project_root>/AGENTS.md
```

Using markers:

```md
<!-- BEGIN VKC AGENT: agent-name -->
...
<!-- END VKC AGENT: agent-name -->
```

## Completion Criteria

> **All criteria met.** `go test ./...`, `go vet ./...`, `go build ./cmd/vkc`, and `golangci-lint run ./...` all pass. CI: Quality ✅ Security ✅.

Original criteria:

```bash
go test ./...
go vet ./...
go build ./cmd/vkc
```

All pass.

Also provide:

- summary of implemented features
- list of files changed
- test output
- known limitations
- next recommended steps

Do not claim success without command output.
