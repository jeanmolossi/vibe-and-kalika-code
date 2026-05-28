# Manifest Spec - Vibe and Kalika Code

## 1. Required File

Every package must contain:

```txt
manifest.yaml
```

At the package root.

---

## 2. Example

```yaml
name: kalika-reviewer-pack
version: 1.0.0
description: Agents and skills for adversarial code review
author: Jean Molossi
license: MIT

targets:
  - copilot-cli
  - claude-code
  - codex-cli

agents:
  - name: kalika-reviewer
    description: Adversarial code reviewer focused on bugs, security, performance, maintainability and architectural drift.
    source: agents/kalika-reviewer.md
    targets:
      copilot-cli:
        scope: user
      claude-code:
        scope: project
      codex-cli:
        mode: agents-md-section

skills:
  - name: code-review
    description: Review code changes for correctness, security, performance, maintainability and task alignment.
    source: skills/code-review
    targets:
      copilot-cli:
        scope: user
      claude-code:
        scope: project
      codex-cli:
        scope: project
```

---

## 3. Package Fields

Required:

```yaml
name: string
version: string
description: string
targets: []string
```

Optional:

```yaml
author: string
license: string
homepage: string
repository: string
```

---

## 4. Agent Fields

Required:

```yaml
name: string
source: string
targets: map
```

Optional:

```yaml
description: string
```

Rules:

- `source` must point to a Markdown file.
- `source` must be inside package root.
- `name` should use kebab-case.
- target platforms must be known.

---

## 5. Skill Fields

Required:

```yaml
name: string
source: string
targets: map
```

Optional:

```yaml
description: string
```

Rules:

- `source` must point to a directory.
- source directory must contain `SKILL.md`.
- source directory may contain:
  - `scripts/`
  - `references/`
  - `assets/`
  - `templates/`
- source must be inside package root.
- target platforms must be known.

---

## 6. Supported Targets

```txt
copilot-cli
claude-code
codex-cli
```

Unknown target names must be rejected.

---

## 7. Scope Semantics

### Copilot CLI

Default scope:

```txt
user
```

Meaning:

```txt
~/.copilot
```

or:

```txt
$COPILOT_HOME
```

### Claude Code

Default scope:

```txt
project
```

Meaning:

```txt
<project_root>/.claude
```

### Codex CLI

Default skill scope:

```txt
project
```

Meaning:

```txt
<project_root>/.agents/skills
```

Default agent mode:

```txt
agents-md-section
```

Meaning: merge into `AGENTS.md` using managed markers.

---

## 8. Validation Rules

The manifest is invalid if:

- `name` is empty
- `version` is empty
- `description` is empty
- `targets` is empty
- any target is unknown
- an agent source does not exist
- an agent source is not Markdown
- a skill source does not exist
- a skill source is not a directory
- a skill source has no `SKILL.md`
- any source path escapes package root
- any source contains path traversal
- two records generate the same target path
- a record declares a platform not present in package-level targets

---

## 9. Future Extensions

Not required in MVP:

```yaml
integrity:
  checksum: sha256:<hash>

dependencies:
  - name: another-pack
    version: ">=1.0.0"

compatibility:
  vkc: ">=0.1.0"

tags:
  - review
  - golang
```
