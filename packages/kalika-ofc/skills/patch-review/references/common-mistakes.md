# Common Review Mistakes

Use this lens to avoid noisy or low-value review comments.

## Common mistakes to avoid

- Blocking on formatting when behavior is fine
- Repeating a style preference as if it were a correctness issue
- Missing the actual bug because the code looks clean
- Focusing on lines changed instead of runtime impact
- Calling out hypothetical problems with no path to reproduce
- Turning a non-blocking concern into a blocker by habit
- Ignoring tests because the code looks “obvious”
- Letting a larger architectural issue hide inside a small diff

## Good review behavior

- Be specific about impact.
- Separate blockers from suggestions.
- Tie comments to a concrete risk or requirement.
- Prefer the smallest fix that addresses the problem.
- Ask whether the issue changes behavior, safety, or maintainability in a meaningful way.

## Review rule

If a comment cannot explain why the issue matters, it probably does not belong in the blocking section.
