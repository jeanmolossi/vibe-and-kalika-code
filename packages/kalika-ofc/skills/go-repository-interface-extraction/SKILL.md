---
name: go-repository-interface-extraction
description: Use this skill when a Go service depends directly on a concrete repository struct (e.g., `*repositories.SomeEntity`) instead of an interface, making the service untestable and violating DIP. Call it before planning or implementing any DIP fix in a Go layered service. It helps extract the correct interface, place it correctly, and wire everything at the entry point without breaking existing tests.
---

# Go Repository Interface Extraction

## Purpose

Extract a Go repository interface from a concrete struct so that services depend on
abstractions rather than on GORM/DB implementations. Enables mocking, respects the
Dependency Inversion Principle, and keeps the service layer free of `*gorm.DB`.

## Trigger

Call this skill when:

- A Go service struct holds a concrete repository type (e.g., `ConlineDB *gorm.DB`, `repo *repositories.SomeEntity`).
- A repository has no public interface, blocking dependency injection and mocking.
- A service constructor receives `*gorm.DB` and instantiates a repository inline.
- An existing interface is defined but unused — the service still takes the concrete type.
- A test cannot compile because it cannot swap out the real DB.

## Do not use when

- The pattern does not involve a Go layered architecture.
- The repository already has a public interface and the service already uses it.
- The change would require touching ORM-scan logic or creating ORM mapping structs — that is a separate, larger refactoring.
- The codebase has no tests at all and the goal is purely structural, with no immediate DIP benefit.

## Required inputs

The agent must provide:

- The concrete repository struct name and file path.
- The list of public methods on that struct that services actually call.
- The list of service files that import the concrete type.
- Whether existing tests create the processor/service directly or via a constructor.
- The entry point file where concrete wiring happens (e.g., `main.go`, `processorManager.go`).

## Procedure

1. **Audit callers**: `grep -r "SomeEntity" services/` to find every service that uses the concrete type. List method calls made on it.

2. **Define the interface in the same file as the implementation** (Go idiom — keeps interface and implementation co-located; avoids import cycles):
   ```go
   // repositories/someRepository.go
   type SomeRepository interface {
       Select(ctx context.Context, id int) (domain.SomeEntity, error)
       // … only methods actually called by services
   }
   ```

3. **Verify implicit satisfaction**: run `go build ./...`. If the concrete struct does not satisfy the interface, the compiler will report the missing method — fix the struct signature.

4. **Update each service constructor**: replace `*repositories.SomeEntity` parameter with `repositories.SomeRepository` (the interface). Remove any `*gorm.DB` fields that existed solely to instantiate the repo.
   ```go
   // Before
   func NewSomeProcessor(db *gorm.DB) *SomeProcessor { … }
   // After
   func NewSomeProcessor(repo repositories.SomeRepository) *SomeProcessor { … }
   ```

5. **Move concrete instantiation to the entry point**: in `main.go` or `processorManager.go`, create the concrete repo and pass it:
   ```go
   repo := &repositories.SomeEntity{Db: conlineDB}
   processor := services.NewSomeProcessor(repo)
   ```

6. **Handle sub-processor injection**: if the service creates sub-processors inline inside a method, move them to constructor parameters using a narrow interface for each. Factory-inside-method is the same anti-pattern one level deeper.

7. **Run validation gates**:
   - `go build ./...` — must pass.
   - `go test ./...` — must pass (no new failures).
   - `golangci-lint run` — must pass.

8. **Check for leaked internal types**: if the repository's methods return internal DTOs (e.g., `companyAreaRepository.UsrIDPerfilID`), replace them with `domain.*` types or create result types in the `repositories/` package. Services must not import internal repository entities.

## Key Go idioms

- **Interfaces are satisfied implicitly** — no `implements` keyword. The concrete struct satisfies the interface if it has all the methods.
- **Define interfaces where they are consumed** or co-locate with the implementation — both are valid Go patterns; prefer co-location when the interface is used by only one consumer layer.
- **Narrow interfaces**: each service should depend only on the methods it calls. Do not copy the entire struct API into the interface.
- **`*gorm.DB` belongs only in `repositories/` implementations and the entry point**. Never in `services/`.

## Pitfalls

| Pitfall | Mitigation |
|---------|-----------|
| Interface method signature does not match struct — build error | Check pointer vs value receiver on the struct |
| Test creates concrete struct directly, bypasses interface | Verify whether tests use the constructor or create the struct directly; if direct, the change is transparent |
| Internal repo DTO leaks into service via interface return type | Return `domain.*` types from the interface; map inside the repo implementation |
| `*gorm.DB` still used in a nested sub-processor | Trace the full call graph before declaring the batch done |
| Breaking a test that imports the concrete type for mocking | Define a minimal mock struct in `_test.go` implementing the interface instead |

## Expected output

Return:

```md
# Repository Interface Extraction Result

## Status
PASS | BLOCKED | NOT_APPLICABLE

## Interface defined
- File: <path>
- Interface name: <name>
- Methods: <list>

## Services updated
- <service file> — <field/param changed>

## Entry-point wiring
- <wiring file> — <concrete instantiation>

## Validation
- go build ./...: PASS | FAIL
- go test ./...: PASS | FAIL
- golangci-lint: PASS | FAIL

## Residual risks
- <any remaining concrete type in a service>
```

## Stop conditions

Stop and return `BLOCKED` when:

- Extracting the interface would require changing ORM scan targets (struct tags, GORM models) — that is a separate migration.
- The repository has 20+ methods and the service uses fewer than 5 — scope the interface to the minimal surface first and confirm with the team before proceeding.
- Removing `*gorm.DB` from a constructor would cascade to more than 10 files — schedule as a dedicated batch, not a side-effect.
