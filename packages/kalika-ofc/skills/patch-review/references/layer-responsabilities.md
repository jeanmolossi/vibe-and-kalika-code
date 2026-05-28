# Layer Responsibilities Lens

Use this lens when the patch crosses boundaries between UI, application, domain, infrastructure, persistence, or external integrations.

## What to check

- Each layer owns the logic that belongs to it.
- Domain rules are not leaking into transport or persistence code.
- UI code is not making business decisions that should live deeper.
- Infrastructure code is not deciding policy.
- Helpers are not turning into a garbage chute for unrelated responsibilities.

## Questions to ask

- Which layer should own this decision?
- Is the patch preserving the dependency direction, or reversing it for convenience?
- Is this a reusable boundary, or just a shortcut around proper separation?
- Will future changes become harder because the responsibility moved to the wrong place?

## Red flags

- Validation split across layers with inconsistent rules
- Database-specific logic in view/controller code
- Business rules embedded in serializers, handlers, or repositories
- Service classes that become god-objects
- Presentation concerns contaminating domain models

## Review rule

When a layer boundary is crossed, require a clear justification. Convenience is not a justification. Sometimes it is just debt with better branding.
