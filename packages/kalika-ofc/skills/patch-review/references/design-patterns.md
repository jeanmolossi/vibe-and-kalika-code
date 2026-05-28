# Design Patterns Lens

Use this lens when the patch introduces abstractions, new interfaces, factories, strategies, adapters, or other pattern-heavy code.

## What to check

- The pattern solves an actual problem in this codebase.
- The abstraction reduces complexity instead of just moving it around.
- The chosen pattern matches the scale and volatility of the problem.
- The pattern does not hide simple behavior behind ceremony.

## Questions to ask

- Is this pattern earning its keep?
- Would plain functions or direct calls be clearer here?
- Is the abstraction stable, or is it trying to predict future needs?
- Does the code get easier to test and reason about, or just more fashionable?

## Red flags

- Factories for objects that could be constructed directly
- Strategy/adapter layers that only wrap one implementation
- Interfaces added before there is a second implementation
- Pattern cargo culting from unrelated code
- Over-abstracting a small feature into a mini-framework

## Review rule

Prefer the simplest design that preserves correctness and future evolution. A pattern is justified when it reduces real change cost, not when it sounds architectural.
