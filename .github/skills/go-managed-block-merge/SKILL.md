---
name: go-managed-block-merge
description: Use this skill when a Go tool must idempotently insert or update named sections inside a shared Markdown config file (e.g., AGENTS.md, CLAUDE.md) without destroying pre-existing content. Call it before implementing any file-merge or config-injection logic. It produces delimited managed blocks that can be inserted, replaced, or preserved in a single regex pass.
---

# Go Managed-Block Merge for Shared Config Files

## Purpose

When multiple agents or tools write to the same Markdown config file (e.g.,
`AGENTS.md` for Codex CLI), a naive overwrite destroys previous content.
This skill defines an HTML-comment-delimited block format and a merge function
that inserts a new block on first write and replaces it on subsequent writes,
leaving all other content untouched.

## Trigger

Call this skill when:

- A tool appends to a shared file that other tools also modify.
- The same tool may run multiple times and must not duplicate its section.
- A config file must preserve human-written content outside managed sections.
- A rollback or uninstall operation must remove exactly one named section.

## Do not use when

Do not call this skill when:

- The tool is the sole writer of the file (plain overwrite is simpler).
- The file format is not line-oriented text (binary, JSON, YAML with schema).
- Sections are ordered and the order must be enforced (use a different mechanism).

## Required inputs

The agent must provide:

- The agent/section name (unique per managed section).
- The source content to inject.
- The target file path.
- The allowed root for security boundary enforcement.

## Block format

```
<!-- BEGIN VKC AGENT: <name> -->
## Agent: <name>

<content>
<!-- END VKC AGENT: <name> -->
```

Delimiters are HTML comments — invisible when rendered, unambiguous in source.

## Procedure

### 1 — `ManagedBlock()` — build the delimited block string

```go
func ManagedBlock(agentName, content string) string {
    trimmed := strings.TrimSpace(content)
    return fmt.Sprintf(
        "<!-- BEGIN VKC AGENT: %s -->\n## Agent: %s\n\n%s\n<!-- END VKC AGENT: %s -->\n",
        agentName, agentName, trimmed, agentName,
    )
}
```

### 2 — `MergeManagedBlock()` — insert or replace, never duplicate

```go
func MergeManagedBlock(existing, agentName, content string) string {
    block := ManagedBlock(agentName, content)
    pattern := regexp.MustCompile(
        `(?s)<!-- BEGIN VKC AGENT: ` + regexp.QuoteMeta(agentName) +
        ` -->.*?<!-- END VKC AGENT: ` + regexp.QuoteMeta(agentName) + ` -->\n?`,
    )
    if pattern.MatchString(existing) {
        // Replace existing block in-place.
        return strings.TrimRight(pattern.ReplaceAllString(existing, block), "\n") + "\n"
    }
    if strings.TrimSpace(existing) == "" {
        return block
    }
    // Append after trimming trailing newlines.
    return strings.TrimRight(existing, "\n") + "\n\n" + block
}
```

### 3 — `MergeAgentFile()` — read, merge, validate boundary, write

```go
func MergeAgentFile(targetPath, agentName, sourcePath, allowedRoot string) error {
    contentBytes, err := os.ReadFile(sourcePath)
    if err != nil {
        return err
    }
    existing, _ := os.ReadFile(targetPath) // empty slice if file doesn't exist
    merged := MergeManagedBlock(string(existing), agentName, string(contentBytes))
    if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
        return err
    }
    if allowedRoot != "" {
        if err := security.EnsureResolvedWithinRoot(allowedRoot, targetPath); err != nil {
            return fmt.Errorf("merge destination escapes allowed root: %w", err)
        }
    }
    return os.WriteFile(targetPath, []byte(merged), 0o644)
}
```

### 4 — Test all three cases

```go
func TestMergeManagedBlock(t *testing.T) {
    // Case 1: empty file — block is the full content
    got := MergeManagedBlock("", "my-agent", "hello")
    assert.Contains(t, got, "<!-- BEGIN VKC AGENT: my-agent -->")

    // Case 2: file with existing block — block is replaced, no duplicate
    existing := ManagedBlock("my-agent", "old content")
    got = MergeManagedBlock(existing, "my-agent", "new content")
    assert.Equal(t, 1, strings.Count(got, "<!-- BEGIN VKC AGENT: my-agent -->"))
    assert.Contains(t, got, "new content")
    assert.NotContains(t, got, "old content")

    // Case 3: file with other content — block is appended, other content preserved
    got = MergeManagedBlock("# My file\n\nsome text\n", "my-agent", "hello")
    assert.Contains(t, got, "# My file")
    assert.Contains(t, got, "<!-- BEGIN VKC AGENT: my-agent -->")
}
```

## Key rules

- Always use `regexp.QuoteMeta(agentName)` — agent names may contain hyphens or dots.
- The `(?s)` flag is required so `.` matches newlines inside the block.
- `existing, _ := os.ReadFile(targetPath)` — ignore the error; a missing file is treated as empty, which is correct.
- Always enforce `allowedRoot` before writing to prevent path traversal.
- Strip trailing newlines before appending to avoid blank-line accumulation on repeated runs.

## Expected output

```md
# Managed-Block Merge Result

## Status
PASS | BLOCKED | NOT_APPLICABLE

## Findings
- ManagedBlock() produces correct delimiters: yes/no
- Replace path tested (no duplicates): yes/no
- Append path tested (preserves other content): yes/no
- Empty-file path tested: yes/no
- Security boundary enforced: yes/no

## Next step
- Wire MergeAgentFile() into the platform adapter's execute step
```

## Stop conditions

Stop and return `BLOCKED` when:

- The target file is not line-oriented text (use a format-specific merger instead).
- The agent name is user-supplied and not sanitized (validate / reject unsafe names first).
- The file is under version control and the merge must produce a clean diff (consider a structured format instead).
