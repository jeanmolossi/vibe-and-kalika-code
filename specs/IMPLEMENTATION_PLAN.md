# Implementation Plan - Vibe and Kalika Code

## Phase 1 - Project foundation

1. Create Go module.
2. Add Cobra root command.
3. Add command stubs:
   - `init`
   - `detect`
   - `install`
   - `validate`
   - `doctor`
4. Add basic terminal styling.
5. Add test structure and `testdata`.

Acceptance:

- `go test ./...` passes.
- `vkc --help` works.
- Commands exist.

---

## Phase 2 - Manifest parser and validator

1. Create manifest structs.
2. Parse `manifest.yaml`.
3. Validate package fields.
4. Validate agents.
5. Validate skills.
6. Validate target names.
7. Validate source paths.
8. Reject path traversal.
9. Reject skill without `SKILL.md`.

Acceptance:

- Valid example package passes.
- Invalid examples fail with clear errors.

---

## Phase 3 - Source providers

1. Implement local source provider.
2. Implement Git source provider.
3. Clone Git repository to temp directory.
4. Ensure temp cleanup.
5. Return package root path.

Acceptance:

- Local source works.
- Git source works with local/fake repository in tests.
- No package script is executed.

---

## Phase 4 - Platform detection

1. Implement platform detector registry.
2. Implement Copilot CLI detection.
3. Implement Claude Code detection.
4. Implement Codex CLI detection.

Detection rules:

- Copilot:
  - `$COPILOT_HOME`, fallback `~/.copilot`
- Claude:
  - project `.claude`
  - user `~/.claude`
- Codex:
  - `$CODEX_HOME`, fallback `~/.codex`
  - project `AGENTS.md`
  - project `.agents/skills`

Acceptance:

- Detection works with fake HOME in tests.
- No real user HOME is touched during tests.

---

## Phase 5 - Copilot adapter

1. Implement target mapping for Copilot agents.
2. Implement target mapping for Copilot skills.
3. Respect `$COPILOT_HOME`.
4. Generate planned operations.

Targets:

```txt
~/.copilot/agents/<agent-name>.md
~/.copilot/skills/<skill-name>/SKILL.md
```

Acceptance:

- Agent target resolves correctly.
- Skill target resolves correctly.
- Custom `COPILOT_HOME` works.

---

## Phase 6 - Planner

1. Build plan from manifest.
2. Filter only selected platforms.
3. Detect existing target files.
4. Detect duplicate target paths.
5. Detect script warnings.
6. Mark conflicts.
7. Produce dry-run output.

Acceptance:

- No filesystem writes happen during planning.
- Existing files become conflicts.
- Script files produce warnings.

---

## Phase 7 - Interactive UI

1. Add welcome flow.
2. Add platform checkbox selection.
3. Add source selection.
4. Show package summary.
5. Show dry-run plan.
6. Show conflict decisions.
7. Confirm final apply.

Acceptance:

- `vkc init` can complete a full interactive path.
- User can select only Copilot.
- User can cancel before apply.

---

## Phase 8 - Installer and security

1. Implement safe file copy.
2. Implement safe directory copy.
3. Implement symlink checks.
4. Block writes outside target roots.
5. Implement conflict actions:
   - skip
   - overwrite
   - backup-and-overwrite
   - write-as-new
6. Do not execute any copied file.

Acceptance:

- Path traversal is blocked.
- Unsafe symlink is blocked.
- Existing file is not overwritten without decision.

---

## Phase 9 - Backup manager

1. Create backup directory.
2. Copy original files before modification.
3. Generate `backup-report.yaml`.

Acceptance:

- Backup is created before overwrite.
- Backup report is valid YAML.

---

## Phase 10 - Codex adapter

1. Map skills to project `.agents/skills`.
2. Merge agents into `AGENTS.md`.
3. Use managed markers:

```md
<!-- BEGIN VKC AGENT: agent-name -->
<!-- END VKC AGENT: agent-name -->
```

4. Preserve manual content outside markers.

Acceptance:

- Creates `AGENTS.md` if missing.
- Adds managed section.
- Updates existing managed section.
- Preserves external content.

---

## Phase 11 - Claude adapter

1. Map agents to project `.claude/agents`.
2. Map skills to project `.claude/skills`.
3. Use project scope by default.

Acceptance:

- Creates project `.claude` folders.
- Installs agents and skills.

---

## Phase 12 - State and report

1. Write `.ai-setup/installed.yaml`.
2. Preserve existing state.
3. Generate Markdown report.

Acceptance:

- State includes package, version, source, selected platforms and files.
- Report includes created, modified, skipped, warnings and backup location.

---

## Phase 13 - Doctor and validation commands

1. Implement `vkc validate`.
2. Implement `vkc doctor`.
3. Add clear terminal output.
4. Add JSON-ready internal structures.

Acceptance:

- `vkc validate ./pack` reports package validity.
- `vkc doctor` reports environment health.

---

## Phase 14 - Tests and hardening

Required tests:

- Manifest parser
- Manifest validator
- Path traversal detection
- Script detection
- Source providers
- Copilot mapping
- Claude mapping
- Codex mapping
- Planner
- Conflicts
- Backup
- Report
- State store
- Git install
- Managed block merge

Final acceptance:

```bash
go test ./...
go vet ./...
go build ./cmd/vkc
```
