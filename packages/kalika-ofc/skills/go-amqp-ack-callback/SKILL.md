---
name: go-amqp-ack-callback
description: Use this skill when a Go service struct holds an `amqp.Delivery` and calls `Ack()`/`Nack()` inside business logic methods, making the service untestable without an AMQP broker. Call it before implementing any AMQP decoupling refactoring. It produces a clean callback-based design where business logic never imports `amqp` directly.
---

# Go AMQP Ack Callback Decoupling

## Purpose

Remove `amqp.Delivery` from service structs and business methods by replacing it with
a typed `ack func(requeue bool) error` callback. The caller (orchestrator/consumer)
creates the closure capturing the real delivery and injects it at construction time.
The service becomes testable without any AMQP broker.

## Trigger

Call this skill when:

- A Go service struct has a field like `Message amqp.Delivery`.
- A service method calls `msg.Ack(false)` or `msg.Nack(false, true)` inside business logic.
- A service file imports `github.com/rabbitmq/amqp091-go` or `github.com/streadway/amqp`.
- `shouldRequeue` / requeue logic using `amqp.Table` or delivery headers lives in a service package instead of the infrastructure/drivers package.
- A unit test for the service cannot run because it needs a real `amqp.Delivery`.

## Do not use when

- AMQP is already decoupled behind a callback or an interface.
- The service is the entry-point/consumer itself (e.g., `processorManager.go`) — that file is allowed to hold the concrete delivery.
- The refactoring scope would require changing more than 3 non-wiring callers of `NewXProcessor()`.

## Required inputs

The agent must provide:

- The service file that holds `amqp.Delivery` and calls `Ack`/`Nack`.
- The orchestrator file that creates the service (single caller preferred).
- The location of any `shouldRequeue`-style logic that reads delivery headers.
- The entry-point file where the AMQP consumer loop runs.

## Procedure

1. **Locate all Ack/Nack call sites** in the service:
   ```bash
   grep -rn "\.Ack\|\.Nack\|amqp\.Delivery" services/
   ```

2. **Move `shouldRequeue` logic to the drivers/infrastructure package** (e.g., `shared/drivers/queue.go`). It must not live in the service. Signature:
   ```go
   // shared/drivers/queue.go
   func ShouldRequeue(msg amqp.Delivery) bool {
       // reads msg.Headers["x-death"] or similar
   }
   ```

3. **Replace `amqp.Delivery` field with `ack func(requeue bool) error` in the service struct**:
   ```go
   // Before
   type JobProcessor struct {
       Message amqp.Delivery
       JobId   int
   }

   // After
   type JobProcessor struct {
       JobId int
       ack   func(requeue bool) error  // unexported — callers cannot bypass it
   }
   ```

4. **Update the constructor** to accept `ack func(requeue bool) error` instead of `amqp.Delivery`:
   ```go
   func NewJobProcessor(jobId int, ack func(requeue bool) error, …) *JobProcessor {
       return &JobProcessor{JobId: jobId, ack: ack, …}
   }
   ```

5. **Replace `j.Message.Ack()` / `j.Message.Nack()` calls in `RunProcess()`** with `j.ack(requeue)`:
   ```go
   // Before
   j.Message.Ack(false)

   // After
   if err := j.ack(false); err != nil {
       return err
   }
   ```

6. **Create the closure in the orchestrator** (the single caller of `NewJobProcessor`). The closure captures the real `amqp.Delivery` and calls `ShouldRequeue` there:
   ```go
   // services/processorManager.go
   func (pm *ProcessorManager) HeapUpConsumers(msg amqp.Delivery) {
       jobId := extractJobId(msg)
       ack := func(requeue bool) error {
           if requeue {
               return msg.Nack(false, drivers.ShouldRequeue(msg))
           }
           return msg.Ack(false)
       }
       processor := jobprocessor.NewJobProcessor(jobId, ack, …)
       go processor.RunProcess(ctx)
   }
   ```

7. **Verify the service no longer imports `amqp`**:
   ```bash
   grep -n "amqp" services/job-processor/job-processor.go
   # Should return zero results
   ```

8. **Make `messageChannel` / delivery-passing fields unexported** in the orchestrator if they were public only to support the old pattern.

9. **Run validation gates**:
   - `go build ./...` — must pass.
   - `go test ./...` — must pass (now tests can inject a fake `ack` func).
   - `golangci-lint run` — must pass.

## Test pattern enabled by this refactoring

After decoupling, unit tests can inject a simple ack:

```go
ackCalled := false
ack := func(requeue bool) error {
    ackCalled = true
    return nil
}
p := NewJobProcessor(42, ack, mockRepo, mockCache)
err := p.RunProcess(context.Background())
assert.NoError(t, err)
assert.True(t, ackCalled)
```

## Key Go idioms

- **Closures as first-class values**: Go closures capture variables by reference, making them ideal for encapsulating AMQP delivery in the caller without passing a concrete type downstream.
- **Unexported `ack` field**: prevents callers from bypassing the callback contract after construction.
- **`func(requeue bool) error` signature**: simple, broker-agnostic, mockable with a single anonymous function in tests.

## Pitfalls

| Pitfall | Mitigation |
|---------|-----------|
| `shouldRequeue` reads `amqp.Delivery.Headers` — needs `amqp` in the driver | `shared/drivers/queue.go` already imports `amqp`; no new import cycle |
| The closure captures `msg` by reference in a loop | Use `msg := msg` shadow variable before the goroutine if dispatch is in a range loop |
| Multiple callers of `NewJobProcessor` need updating | If more than one caller exists, check each before changing the signature |
| Nack with `multiple=false, requeue=bool` vs Ack — conflated in one `ack func` | Use a two-parameter signature if both multiple and requeue flags are needed: `ack func(success bool, requeue bool) error` |
| Test panics because `ack` is nil | Guard in the service: `if p.ack == nil { return errors.New("ack not set") }` |

## Expected output

Return:

```md
# AMQP Ack Callback Result

## Status
PASS | BLOCKED | NOT_APPLICABLE

## Changes made
- Service struct: <field removed, ack func added>
- Constructor: <signature change>
- Business method: <Ack/Nack replaced>
- shouldRequeue: <moved to shared/drivers/queue.go>
- Orchestrator: <closure created in HeapUpConsumers>

## Import verification
- `amqp` import removed from service: YES | NO

## Validation
- go build ./...: PASS | FAIL
- go test ./...: PASS | FAIL
- golangci-lint: PASS | FAIL

## Test pattern
- Fake ack injectable: YES | NO
```

## Stop conditions

Stop and return `BLOCKED` when:

- More than one orchestrator calls `NewJobProcessor` with a delivery — align on a single wiring point first.
- The `shouldRequeue` logic depends on service-layer business rules (not just header inspection) — do not move it to `drivers/`; create a strategy function or interface instead.
- The delivery carries business payload (beyond the routing key/job ID) that the service actually needs — extract that payload before decoupling, not inside the callback.
