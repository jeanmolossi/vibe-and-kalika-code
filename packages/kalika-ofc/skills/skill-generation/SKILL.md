---
name: skill-generation
description: Use this skill when an agent must convert stable, reusable project learning into a new Copilot skill under .github/skills/<skill-name>/SKILL.md. Call it during learning workflows after knowledge extraction and filtering.
---

# Skill Generation Skill

## Purpose

Generate reusable Copilot skills from stable project knowledge.

This skill prevents noisy memory dumps by forcing each generated skill to have a clear trigger, purpose, inputs, outputs, and stop conditions.

## Trigger

Call this skill when:

- A learning workflow has extracted reusable knowledge.
- The learning is likely to be reused across future tasks.
- A recurring project-specific pattern should become an operational procedure.
- A mistake should be prevented in future workflows.
- The parent explicitly asks to create skills from project learnings.

## Do not use when

Do not call this skill when:

- The learning is a one-off task detail.
- The learning contains secrets or sensitive data.
- The learning is speculative.
- An existing skill already covers the same behavior.
- The skill would require loading huge context dumps.
- The task is not completed or validated.

## Required inputs

The agent must provide:

- Learning summary.
- Evidence source.
- Why the learning is reusable.
- Existing skills list.
- Proposed skill name.
- Trigger conditions.
- Required inputs.
- Expected outputs.
- Stop conditions.

## Skill file template

Generated skills must follow this structure:

```md
---
name: <skill-name>
description: Use this skill when <specific trigger>. Call it <phase/timing>. It helps <purpose>.
---

# <Skill Title>

## Purpose

<what this skill does>

## Trigger

Call this skill when:

- 

## Do not use when

Do not call this skill when:

- 

## Required inputs

The agent must provide:

- 

## Procedure

1. 
2. 
3. 

## Expected output

Return:

```md
# <Skill Result>

## Status
PASS | BLOCKED | NOT_APPLICABLE

## Findings
- 

## Next step
- 
```

## Stop conditions

Stop and return `BLOCKED` when:

- 
```

## Expected output

Return:

```md
# Skill Generation Result

## Status
CREATED | PROPOSED | REJECTED | BLOCKED

## Skill path
.github/skills/<skill-name>/SKILL.md

## Skill trigger
<when it should be called>

## Reason this is reusable
<why it deserves to exist>

## Duplicate check
<existing skills checked>

## Safety check
<what sensitive/noisy data was removed>

## Next step
<what the parent workflow should do>
```

## Stop conditions

Stop and return `REJECTED` when:

- The learning is too specific.
- The learning is not actionable.
- The same skill already exists.
- The source evidence is weak.
- The skill would store sensitive or noisy data.
