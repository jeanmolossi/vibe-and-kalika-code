# Refactoring Lens

Use this lens when the patch reshapes existing code, splits responsibilities, or changes structure without changing behavior.

## What to check

- The refactor preserves behavior unless the task explicitly says otherwise.
- Changes are small enough to verify in isolation.
- Each step is reversible or at least understandable.
- Tests or evidence protect the behavior during the move.
- Names and boundaries improve after the change, not just the line count.

## Questions to ask

- Is this a safe refactor, or a hidden rewrite?
- Can I explain the behavior preservation step by step?
- Did the author change structure and behavior in the same move?
- Is there a rollback path if something turns out wrong?

## Red flags

- Large rewrites disguised as refactoring
- Moving code across layers with no clear reason
- Simultaneous behavior change plus structural change without tests
- Dead code left behind after extraction
- “Cleanup” that quietly alters edge-case behavior

## Review rule

If the patch claims to be a refactor, demand stronger evidence that behavior did not change.
