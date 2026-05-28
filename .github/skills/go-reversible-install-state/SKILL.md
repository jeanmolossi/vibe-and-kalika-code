---
name: go-reversible-install-state
description: Use this skill when a Go CLI must record install-time state so a future uninstall command can fully reverse the installation. Call it before implementing the state struct, the install recorder, or the uninstall executor. It prevents the three most common uninstall bugs: wrong block-removal key, accidental shared-file deletion, and non-retryable cleanup ordering.
---

# Go Reversible Install State

## Purpose

When a `vkc install` (or any Go CLI install command) writes files and injects
managed blocks into shared config files, the install record must capture enough
information so a later `vkc uninstall` can reverse every side-effect exactly,
idempotently, and without corrupting other packages' content.

Three design rules prevent the most common uninstall bugs:

1. **Block identity ≠ artifact identity** — the key used to locate and remove a
   managed block (the *agent name*) is not the same as the package name.
2. **Shared files are never in `CreatedFiles`** — files that use managed-block
   merge must be removed via `RemoveManagedBlock`, not `os.Remove`.
3. **State is written last** — cleanup side effects run before the state record
   is updated so a partial failure leaves the record intact for retry.

## Trigger

Call this skill when:

- A Go CLI command installs packages that write files and inject managed blocks
  into a shared config file (e.g., `AGENTS.md`).
- You are designing or extending the install-state struct (`Installation`).
- You are implementing an uninstall command that must reverse a prior install.
- A code review finds that uninstall uses the package name as the block-removal
  key (this is always a bug).

## Do not use when

Do not call this skill when:

- The tool is the sole writer of all files it manages (no shared config files).
- Uninstall is not a requirement.
- The file format is binary or structured (JSON/YAML with schema); managed-block
  merge is not appropriate.

## Required inputs

The agent must provide:

- The `Installation` state struct (or a description of it).
- The list of file types the installer writes: owned files vs. shared
  managed-block files.
- The agent name(s) used as block keys in `ManagedBlock()` / `MergeManagedBlock()`.
- The package name(s) as the user-facing install identifier.

## Procedure

### 1 — Separate `AgentBlock` from `CreatedFiles` in the state struct

```go
// AgentBlock records a managed block that was merged into a shared config file.
// AgentName is the key used in <!-- BEGIN VKC AGENT: <AgentName> -->.
// It may differ from the package name.
type AgentBlock struct {
    Path      string `yaml:"path"`
    AgentName string `yaml:"agent_name"`
}

type Installation struct {
    Package      string       `yaml:"package"`       // user-facing identifier
    CreatedFiles []string     `yaml:"created_files"` // files to os.Remove/os.RemoveAll
    AgentBlocks  []AgentBlock `yaml:"agent_blocks"`  // blocks to RemoveManagedBlock
    BackupPath   string       `yaml:"backup_path"`
    // ...
}
```

**Rule:** Never add a shared managed-block file (e.g., `AGENTS.md`) to
`CreatedFiles`. It belongs in `AgentBlocks` only.

### 2 — Record both fields at install time

```go
// inside the install function, after MergeAgentFile() succeeds:
rec.AgentBlocks = append(rec.AgentBlocks, state.AgentBlock{
    Path:      agentsFilePath,   // absolute path to AGENTS.md
    AgentName: agentSpec.Name,   // e.g. "kalika-reviewer" — NOT the package name
})
// do NOT append agentsFilePath to rec.CreatedFiles
```

### 3 — Uninstall in this exact order

```go
func Uninstall(opts UninstallOptions) (*UninstallResult, int, error) {
    // 1. Validate: package exists in state
    rec := findInstallationRecord(opts.Package) // error if not found

    // 2. Remove managed blocks (idempotent — missing block/file = no-op)
    for _, block := range rec.AgentBlocks {
        security.EnsureWithinRoot(opts.ProjectRoot, block.Path) // path-traversal guard
        installer.RemoveManagedBlock(block.Path, block.AgentName)
    }

    // 3. Restore backups (if any)
    if rec.BackupPath != "" {
        backup.Restore(rec.BackupPath, opts.ProjectRoot)
    }

    // 4. Delete owned files (os.Stat → os.RemoveAll for dirs, os.Remove for files)
    for _, f := range rec.CreatedFiles {
        info, err := os.Stat(f)
        if os.IsNotExist(err) { continue } // idempotent
        if info.IsDir() {
            os.RemoveAll(f)
        } else {
            os.Remove(f)
        }
    }

    // 5. Write updated state LAST — ensures retryability on partial failure
    removeRecordFromState(opts.Package)
    state.Write(opts.ProjectRoot, updatedState)

    return result, ExitSuccess, nil
}
```

**Why state is last:** if steps 2–4 partially fail, the original state record
is still intact on the next invocation, so the user can retry without data loss.

### 4 — Test the agent-name ≠ package-name case explicitly

This is the highest-value test — it catches the most common bug:

```go
t.Run("agent_name_used_not_package_name", func(t *testing.T) {
    // block key is "kalika-reviewer"; package name is "basic-kalika-pack"
    content := "<!-- BEGIN VKC AGENT: kalika-reviewer -->...\n<!-- END VKC AGENT: kalika-reviewer -->\n"
    os.WriteFile(agentsPath, []byte(content), 0o644)
    writeState(t, root, &state.Store{Installations: []state.Installation{
        {
            Package:     "basic-kalika-pack",
            AgentBlocks: []state.AgentBlock{{Path: agentsPath, AgentName: "kalika-reviewer"}},
        },
    }})
    app.Uninstall(app.UninstallOptions{Package: "basic-kalika-pack", ProjectRoot: root})
    got, _ := os.ReadFile(agentsPath)
    // block gone, surrounding content preserved
    assert.NotContains(t, string(got), "BEGIN VKC AGENT: kalika-reviewer")
    assert.Contains(t, string(got), "# Header") // surrounding content preserved
})
```

### 5 — Dry-run must branch output on the flag

```go
if dryRun {
    fmt.Fprintf(out, "Would uninstall %s (dry-run)\n", res.Package)
} else {
    fmt.Fprintf(out, "Uninstalled %s\n", res.Package)
}
```

**Common bug:** printing `"Uninstalled"` unconditionally even when `--dry-run`
is set. The `dryRun` bool is always in closure scope — use it for both the
operation skip AND the output message.

## Expected output

```md
# Reversible Install State Result

## Status
PASS | BLOCKED | NOT_APPLICABLE

## Findings
- AgentBlock struct separates path from agent name: yes/no
- Shared config file absent from CreatedFiles: yes/no
- State write is the last step in uninstall: yes/no
- agent_name != package_name test exists: yes/no
- Dry-run output branches on flag: yes/no

## Next step
- Wire AgentBlock recording into the install recorder
- Add the agent_name_used_not_package_name sub-test to uninstall_test.go
```

## Stop conditions

Stop and return `BLOCKED` when:

- The agent name is user-supplied and not sanitized before use as a block key
  (validate / reject before recording).
- The install command does not use `MergeManagedBlock` — this skill assumes the
  managed-block merge pattern is already in use (see `go-managed-block-merge`).
- The project has no shared config files; all installed files are owned
  exclusively by the package — `CreatedFiles` alone is sufficient.
