# CLI Spec - Vibe and Kalika Code

## 1. Binary

```bash
vkc
```

Official product name:

```txt
Vibe and Kalika Code
```

---

## 2. Commands

MVP commands:

```bash
vkc init
vkc detect
vkc install <source>
vkc validate <source>
vkc doctor
```

---

## 3. `vkc init`

Interactive setup command.

### Flow

```txt
1. Show welcome
2. Detect project root
3. Detect supported platforms
4. Show platform checkboxes
5. Ask package source:
   - local directory
   - Git repository
6. Fetch package
7. Parse manifest
8. Validate manifest
9. Show package summary
10. Build dry-run plan
11. Show conflicts and warnings
12. Ask conflict actions
13. Ask final confirmation
14. Create backups
15. Apply install
16. Write state
17. Generate report
18. Show final summary
```

### Platform checkbox example

```txt
Select installation targets:

[x] GitHub Copilot CLI   detected at ~/.copilot
[ ] Claude Code          detected in project ./.claude
[ ] OpenAI Codex CLI     not detected, can still install project files
```

---

## 4. `vkc detect`

Detects known platform directories.

Example output:

```txt
Detected platforms:

[✓] GitHub Copilot CLI
    home: ~/.copilot
    agents: ~/.copilot/agents
    skills: ~/.copilot/skills

[✓] Claude Code
    project: ./.claude
    agents: ./.claude/agents
    skills: ./.claude/skills

[✓] OpenAI Codex CLI
    home: ~/.codex
    project instructions: ./AGENTS.md
    skills: ./.agents/skills
```

No filesystem mutation allowed.

---

## 5. `vkc install <source>`

Installs a package from source.

Supported sources:

```bash
vkc install ./packs/kalika-reviewer
vkc install https://github.com/org/agent-pack.git
vkc install git@github.com:org/agent-pack.git
```

The command may be interactive by default.

Future flags:

```bash
vkc install ./pack --targets copilot-cli,claude-code
vkc install ./pack --yes
vkc install ./pack --dry-run
vkc install ./pack --json
```

For MVP, `--dry-run` may exist, but dry-run must happen anyway before apply.

---

## 6. `vkc validate <source>`

Validates package source.

Checks:

- manifest exists
- manifest fields are valid
- sources exist
- skills contain `SKILL.md`
- agents are Markdown files
- no path traversal
- no unknown target
- no duplicate operation targets

No filesystem mutation allowed.

---

## 7. `vkc doctor`

Checks environment health.

Should report:

- project root
- detected platforms
- relevant environment variables
- write permissions
- Git availability
- `.ai-setup/installed.yaml` validity
- potential conflicts

No filesystem mutation allowed, except maybe reading state.

---

## 8. Terminal UX Rules

The CLI should be:

- clear
- colorful but not childish
- concise
- explicit before dangerous changes
- calm on errors
- strict on security

No hidden writes before final confirmation.

---

## 9. Exit Codes

Suggested:

| Code | Meaning |
|---:|---|
| 0 | success |
| 1 | generic error |
| 2 | validation error |
| 3 | security violation |
| 4 | user cancelled |
| 5 | source fetch error |
| 6 | install conflict unresolved |

---

## 10. Error Message Style

Bad:

```txt
failed
```

Good:

```txt
Manifest validation failed:
- skill "code-review" source exists, but SKILL.md was not found
- target "alien-cli" is not supported
```

Baka-level vague errors are forbidden.
