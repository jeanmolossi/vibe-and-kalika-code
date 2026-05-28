# File Mapping Spec - Vibe and Kalika Code

## 1. Platform identifiers

```txt
copilot-cli
claude-code
codex-cli
```

---

## 2. GitHub Copilot CLI

Priority:

```txt
must-have
```

### 2.1 Base directory

Use:

```txt
$COPILOT_HOME
```

If defined.

Fallback:

```txt
~/.copilot
```

### 2.2 Agents

Target:

```txt
<copilot_home>/agents/<agent-name>.md
```

Example:

```txt
~/.copilot/agents/kalika-reviewer.md
```

### 2.3 Skills

Target:

```txt
<copilot_home>/skills/<skill-name>/SKILL.md
```

Example:

```txt
~/.copilot/skills/code-review/SKILL.md
```

The whole skill directory must be copied, not only `SKILL.md`.

---

## 3. Claude Code

Priority:

```txt
nice-to-have
```

### 3.1 Default scope

Use project scope by default.

Base:

```txt
<project_root>/.claude
```

### 3.2 Agents

Target:

```txt
<project_root>/.claude/agents/<agent-name>.md
```

### 3.3 Skills

Target:

```txt
<project_root>/.claude/skills/<skill-name>/SKILL.md
```

The whole skill directory must be copied.

### 3.4 Future user scope

Future user scope may use:

```txt
~/.claude/agents/<agent-name>.md
~/.claude/skills/<skill-name>/SKILL.md
```

---

## 4. OpenAI Codex CLI

Priority:

```txt
nice-to-have
```

### 4.1 Skills

Project target:

```txt
<project_root>/.agents/skills/<skill-name>/SKILL.md
```

User target:

```txt
$HOME/.agents/skills/<skill-name>/SKILL.md
```

Default for MVP:

```txt
project
```

The whole skill directory must be copied.

### 4.2 Agents

Codex CLI uses `AGENTS.md` as instructions.

Project target:

```txt
<project_root>/AGENTS.md
```

User target:

```txt
$CODEX_HOME/AGENTS.md
```

Fallback user target:

```txt
~/.codex/AGENTS.md
```

Default for MVP:

```txt
project
```

### 4.3 Managed block format

Agent content must be inserted using managed markers:

```md
<!-- BEGIN VKC AGENT: kalika-reviewer -->
## Agent: kalika-reviewer

<agent markdown content>
<!-- END VKC AGENT: kalika-reviewer -->
```

Rules:

- If block does not exist, append it.
- If block exists, replace only that block.
- Preserve all content outside managed markers.
- Never delete the full `AGENTS.md`.
- Always backup before modifying.

---

## 5. `.ai-setup` files

The CLI writes internal setup artifacts into project root:

```txt
.ai-setup/
  installed.yaml
  backups/
  reports/
```

### 5.1 Backups

```txt
.ai-setup/backups/<timestamp>/
```

### 5.2 Reports

```txt
.ai-setup/reports/<timestamp>-install-report.md
```

### 5.3 State

```txt
.ai-setup/installed.yaml
```

---

## 6. Path Safety

Before writing:

1. Resolve target path.
2. Resolve symlink path.
3. Ensure target remains under allowed root.
4. Reject path traversal.
5. Reject absolute source paths from manifest.
6. Reject attempts to write outside platform target roots.

Filesystem is not a vibes-based democracy. Validate everything.
