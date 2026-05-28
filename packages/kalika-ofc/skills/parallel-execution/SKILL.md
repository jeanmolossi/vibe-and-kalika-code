---
name: parallel-execution
description: Use this skill when deciding whether work can be split across multiple agents. Call it before parallel implementation, research, validation, or documentation to avoid file conflicts, shared contract changes, sequencing bugs, and merge chaos.
---

# Parallel Execution Skill

## Purpose

Determine whether a task can be safely split across multiple agents and define conflict-free work packages.

This skill exists because “just parallelize it” is how humans create distributed garbage fires with extra ceremony.

## Trigger

Call this skill when:

- The plan contains multiple independent tasks.
- The orchestrator considers launching multiple implementers.
- Research can be split by domain, module, or technology.
- Validation can be split by test type or risk area.
- Documentation can be split by artifact type.
- A large task needs faster execution without corrupting shared state.

## Do not use when

Do not call this skill when:

- The task touches one tightly coupled change.
- The plan is still ambiguous.
- Shared interfaces are not stable.
- The next step is only routing or phase-gate evaluation.
- There is no need for multiple agents.

## Required inputs

The agent must provide:

- Approved plan or task list.
- Known files/modules likely to change.
- Dependency map, if available.
- Shared contracts/interfaces.
- Validation strategy.
- Risk areas.
- Merge constraints.

## Safe to parallelize when

Parallelization is allowed when:

- Tasks touch different modules.
- Tasks touch different files.
- Interfaces/contracts are already stable.
- No shared migration is involved.
- No shared generated file is involved.
- Merge order does not affect behavior.
- Each package has independent validation.

## Do not parallelize when

Parallelization is blocked when:

- Tasks modify the same files.
- Tasks change shared interfaces.
- Tasks change database schema.
- Tasks change authentication or authorization behavior.
- Tasks modify global config.
- Tasks require strict sequence.
- Tests are likely to conflict.
- The rollback plan is unclear.

## Expected output

Return this exact structure:

```md
# Parallelization Decision

## Mode
sequential | parallel | partially_parallel

## Reason
<why this mode was selected>

## Work packages
| Package | Agent type | Files/modules | Dependencies | Conflict risk | Decision |
|---|---|---|---|---|---|

## Required ordering
<if sequential or partially parallel>

## Merge strategy
<how outputs should be combined>

## Validation per package
<what must be validated for each package>

## Blockers
<list, or none>
```

## Stop conditions

Stop and return `sequential` when:

- File boundaries are unknown.
- Shared contracts are unstable.
- Parallel packages cannot be validated independently.
- Merge strategy is unclear.
