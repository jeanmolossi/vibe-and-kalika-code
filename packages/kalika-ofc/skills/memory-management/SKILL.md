---
name: memory-management
description: Use this skill when the agent needs to retrieve, correlate, or append reusable project/session/user memories without loading the full memory dump into context. Token-efficient memory retrieval and registration skill. 
---

# memory-management

## Purpose

Manage reusable memories with minimal token cost.
This skill exists to prevent agents from loading the entire memory file into context.
The agent MUST use a local memory selector script to retrieve only relevant memory blocks.
The agent MUST NOT read, cat, print, summarize, or load the full memory file.

## Memory File

Default memory file:

```txt
~/.vkc/memory/MEMORY.md
```

Memory format:

```md
## <timestamp>
Contexto: <how and when this memory should be used>
Memória:
<reusable memory content>
Tags: <comma-separated tags useful for retrieval>
```

Example:

```md
## 2026-05-04T14:35:00Z
Contexto: Use when reviewing Go code that interacts with MySQL 5.7 and needs safe query behavior.
Memória:
The project uses MySQL 5.7. Avoid CTEs, window functions, and syntax only available in MySQL 8+. Prefer explicit JOINs, stable indexes, and parameterized queries through database/sql.
Tags: go, mysql, mysql57, database, sql, backend, review
```

## When To Use

Use this skill before:

- planning a task that may depend on previous project knowledge
- reviewing code with project-specific conventions
- implementing changes in a known repository
- generating agent instructions based on previous workflow decisions
- answering questions that may depend on stored project/user/session knowledge
- extracting reusable lessons after completing a workflow

## Hard Rules

- NEVER load the entire memory file into the agent context.
- NEVER run `cat ~/.vkc/memory/MEMORY.md`.
- NEVER ask another agent to read the full memory dump.
- ALWAYS retrieve memories through `./scripts/select-memory.sh`.
- Only use returned memory blocks.
- If no relevant memory is found, continue without memory instead of forcing irrelevant context.
- Do not treat memory as truth if the current task contradicts it.
- Current user input, repository files, tests, and command outputs take priority over memory.
- Memory is supporting context, not authority.

## Retrieval Workflow

### Step 1: Build Retrieval Query

Extract 3 to 10 strong keywords from the current task.

Prefer:

- project name
- language
- framework
- database
- cloud provider
- architectural pattern
- workflow name
- agent name
- domain concept
- recurring technical constraint

Avoid generic words:

- task
- code
- fix
- implement
- update
- check
- improve

### Step 2: Run Memory Selector

Use:

```bash
bash ./scripts/select-memory.sh \
  --file ~/.vkc/memory/MEMORY.md \
  --query "<task-specific keywords>" \
  --tags "<optional,tags,comma,separated>" \
  --limit 5 \
  --max-chars 6000
```

Example:

```bash
bash ./scripts/select-memory.sh \
  --file ~/.vkc/memory/MEMORY.md \
  --query "golang mysql 5.7 database/sql query review performance" \
  --tags "go,mysql,mysql57,backend,review" \
  --limit 5 \
  --max-chars 6000
```

### Step 3: Use Only Returned Blocks

The agent may use only the selected memory blocks returned by the script.

If the script returns:

```txt
NO_RELEVANT_MEMORY_FOUND
```

Then proceed without memory context.

### Step 4: Retry Once If Needed

If the result is empty but memory is likely relevant, retry once with broader terms.

Example:

```bash
bash ./scripts/select-memory.sh \
  --file ~/.vkc/memory/MEMORY.md \
  --query "golang backend review" \
  --tags "go,backend" \
  --limit 3 \
  --max-chars 4000
```

Do not retry endlessly.

## Token Budget

Default retrieval budget:

```txt
limit: 5 memory blocks
max-chars: 6000
```

For small tasks:

```txt
limit: 3
max-chars: 3000
```

For complex planning/review tasks:

```txt
limit: 8
max-chars: 9000
```

Never exceed:

```txt
limit: 10
max-chars: 12000
```

unless explicitly instructed by the coordinator.

## Memory Relevance Rules

A memory is relevant when it helps with at least one of:

- avoiding a known mistake
- following an established project convention
- preserving a previous architectural decision
- respecting a technology constraint
- understanding a recurring workflow
- applying a known user preference
- maintaining consistency across agents or sessions

A memory is not relevant when it is:

- purely historical
- unrelated to the current task
- too vague to act on
- contradicted by current repository evidence
- only emotionally/contextually interesting but not useful

## Appending New Memories

Append memory only when the information is reusable across future tasks.

Good memory candidates:

- project-specific constraints
- recurring bugs and fixes
- architecture decisions
- workflow decisions
- naming conventions
- command/tooling quirks
- compatibility constraints
- user preferences that affect future output
- lessons learned after validation failures

Bad memory candidates:

- one-off task status
- temporary observations
- raw logs
- huge code snippets
- generic programming advice
- secrets, credentials, tokens, private keys
- sensitive information not explicitly requested to be stored

## Append Format

Use UTC timestamp.

```bash
mkdir -p ~/.vkc/memory

cat >> ~/.vkc/memory/MEMORY.md <<'EOF_MEMORY'

## <timestamp>
Contexto: <when/how this memory should be used>
Memória:
<short reusable memory>
Tags: <tag1, tag2, tag3>
EOF_MEMORY
```

Example:

```bash
mkdir -p ~/.vkc/memory

cat >> ~/.vkc/memory/MEMORY.md <<EOF_MEMORY

## $(date -u '+%Y-%m-%dT%H:%M:%SZ')
Contexto: Use when generating or reviewing GitHub Copilot multi-agent workflow files.
Memória:
The coordinator is the only agent allowed to orchestrate the pipeline. Sub-agents must return their outputs to the coordinator and must never call each other or advance the workflow by themselves.
Tags: agents, coordinator, workflow, github-copilot, orchestration
EOF_MEMORY
```

## Memory Writing Rules

When writing memory:

- keep it short
- keep it reusable
- include clear usage context
- include searchable tags
- prefer technical specificity over vague summaries
- do not store full transcripts
- do not store full patches
- do not store sensitive secrets
- do not duplicate an existing memory unless updating or correcting it

## Memory Output Contract

When using selected memories in a task, the agent should briefly report:

```txt
Selected memories used:
- <timestamp> — <why it was relevant>
```

If no memory was used:

```txt
No relevant memory found.
```

Do not paste all retrieved memories into the final response unless explicitly requested.

## Failure Behavior

If the memory file does not exist:

- continue without memory
- optionally create the directory when appending a new memory

If the selector script does not exist:

- report that memory retrieval is unavailable
- do not read the memory file directly
- ask the coordinator or user to install the script

If selected memories are stale or contradicted:

- ignore the stale memory
- prefer current evidence
- optionally append a corrective memory after the task is complete

## Coordinator Interaction

This skill does not advance the pipeline.

After retrieving, applying, or appending memory, the agent must return control to the coordinator with:

- memories selected
- relevance summary
- decisions influenced by memory
- any stale/contradictory memory discovered
- suggested memory updates, if any

The coordinator decides the next workflow step.
