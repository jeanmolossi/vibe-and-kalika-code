# Clean Code Lens

Use this lens when the patch changes naming, structure, control flow, or function size.

## What to check

- Names match the intent of the code.
- A function does one thing, or at least one coherent thing.
- Duplicate logic is not spreading unless the duplication is deliberate and temporary.
- Branching is not hiding the real behavior.
- Early returns improve clarity instead of scattering the logic.
- Comments explain why, not what the code already says.

## Questions to ask

- Does this code read like it was designed, or assembled under pressure?
- Does the function force the reader to keep too many things in working memory?
- Is the code simple because it is simple, or because complexity was moved somewhere worse?
- Would a future maintainer understand the intent without reverse-engineering it?

## Red flags

- Long functions with mixed responsibilities
- Nested conditionals that bury the main path
- Repeated logic that should be a helper
- Vague names like `data`, `item`, `handleStuff`, `process`
- Comments that describe the obvious instead of the risk
- Clever code that is harder to validate than the problem deserves

## Review rule

Do not block on cleanliness alone unless it creates a real maintenance or correctness problem.
