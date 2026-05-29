---
name: go-http-doer-injection
description: Use this skill when a Go function or service makes HTTP calls that must be unit-tested without a real network. Call it before implementing any HTTP client logic that needs test coverage. It produces an httpDoer interface, a package-level default client, and an internal WithClient variant so tests inject a mock without external libraries.
---

# Go HTTP Doer Injection

## Purpose

Make HTTP-calling code unit-testable without a real network by defining a
minimal `httpDoer` interface and exposing an internal `funcWithClient(opts, client httpDoer)`
variant. Tests inject a `mockHTTPClient` that returns pre-canned responses.
No external test libraries are required.

## Trigger

Call this skill when:

- A Go function calls `http.DefaultClient.Do(req)` or `(&http.Client{}).Do(req)` and you want to unit test it.
- A code review flags that tests require a live network or a running HTTP server.
- You are implementing a feature that downloads files, calls an API, or fetches metadata over HTTP and needs offline test coverage.
- You need to simulate HTTP error responses (4xx, 5xx, timeouts) in tests.

## Do not use when

Do not call this skill when:

- The HTTP call is already covered by an integration test suite and offline coverage is not required.
- The codebase already uses a full HTTP mock library (e.g., `httptest.NewServer`) and that is the project convention.
- The function under test is a thin wrapper with no logic beyond a single `client.Do(req)` call.

## Required inputs

The agent must provide:

- The package where the HTTP-calling logic lives.
- The function(s) that make HTTP calls.
- The test scenarios needed (success, 4xx, 5xx, timeout, body content).

## Procedure

### 1. Define the `httpDoer` interface in the same package

```go
// httpDoer abstracts *http.Client so tests can inject a mock.
type httpDoer interface {
    Do(req *http.Request) (*http.Response, error)
}

var defaultHTTPClient httpDoer = &http.Client{Timeout: 30 * time.Second}
```

Place this near the top of the file that makes HTTP calls. Do **not** export
the interface unless other packages need it.

### 2. Split the public function into a public shell + internal WithClient variant

```go
// MyFunc is the public entry point; it uses the shared default client.
func MyFunc(opts MyOptions) (*MyResult, error) {
    return myFuncWithClient(opts, defaultHTTPClient)
}

// myFuncWithClient is the testable implementation; tests inject their mock here.
func myFuncWithClient(opts MyOptions, client httpDoer) (*MyResult, error) {
    req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, opts.URL, nil)
    if err != nil {
        return nil, err
    }
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // IMPORTANT: always check the status code before reading the body.
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, opts.URL)
    }

    // ... process resp.Body ...
    return result, nil
}
```

> **Rule**: Check `resp.StatusCode != http.StatusOK` immediately after `client.Do(req)`,
> before reading or scanning the body. Error bodies (403 rate limit, 404 not found)
> will otherwise be scanned/written, producing misleading downstream errors.

### 3. Write the mock in the test file

```go
type mockHTTPClient struct {
    responses []*http.Response
    idx       int
}

func (m *mockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
    if m.idx >= len(m.responses) {
        return &http.Response{
            StatusCode: http.StatusNotFound,
            Body:       io.NopCloser(strings.NewReader("")),
        }, nil
    }
    resp := m.responses[m.idx]
    m.idx++
    return resp, nil
}

// helpers
func jsonBody(v any) io.ReadCloser {
    b, _ := json.Marshal(v)
    return io.NopCloser(strings.NewReader(string(b)))
}

func stringBody(s string) io.ReadCloser {
    return io.NopCloser(strings.NewReader(s))
}
```

### 4. Use the mock in table-driven tests

```go
func TestMyFunc_Success(t *testing.T) {
    client := &mockHTTPClient{
        responses: []*http.Response{
            {StatusCode: http.StatusOK, Body: jsonBody(myPayload)},
        },
    }
    result, err := myFuncWithClient(MyOptions{URL: "http://example.com"}, client)
    // assert result ...
}

func TestMyFunc_HTTPError(t *testing.T) {
    client := &mockHTTPClient{
        responses: []*http.Response{
            {StatusCode: http.StatusForbidden, Body: stringBody("rate limited\n")},
        },
    }
    _, err := myFuncWithClient(MyOptions{URL: "http://example.com"}, client)
    if err == nil || !strings.Contains(err.Error(), "403") {
        t.Errorf("expected 403 error, got %v", err)
    }
}
```

### 5. Build-var tests: do not hardcode the default version string

When testing version-dependent logic (e.g., "is current version already latest?"),
never hardcode the default value of a ldflags-injected variable (e.g., `"vdev"`).
CI builds with real ldflags will break those tests. Use the variable itself:

```go
// BAD — breaks when built with real ldflags
currentTag := "vdev"

// GOOD — always matches whatever the current default is
currentTag := "v" + version.Version
```

## Expected output

```md
# HTTP Doer Injection Result

## Status
PASS | BLOCKED | NOT_APPLICABLE

## Findings
- httpDoer interface defined: yes/no
- defaultHTTPClient wired: yes/no
- WithClient internal variant present: yes/no
- Status code checked before body read: yes/no
- Mock implementation present in test file: yes/no

## Next step
- If PASS: run `go test -race ./...` to confirm all HTTP paths are covered offline.
- If BLOCKED: document why injection is not possible and propose alternative (httptest.NewServer).
```

## Stop conditions

Stop and return `BLOCKED` when:

- The codebase already uses `httptest.NewServer` as the established convention for HTTP mocking and changing it would break existing test patterns.
- The HTTP call cannot be wrapped (e.g., it is inside a generated client with no extension point).
- The function has no meaningful logic beyond a single passthrough call.

## Reference implementation

Evidence: `internal/app/update.go` and `internal/app/update_test.go` in
`jeanmolossi/vibe-and-kalika-code` (branch `feat/vkc-update`, commit `f6b1e96`).
