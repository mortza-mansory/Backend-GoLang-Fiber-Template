# Testing

This directory is reserved for integration and end-to-end tests that span
multiple packages (e.g. HTTP API tests against a live database and Redis).

## Running tests

```bash
# All tests
go test ./...

# With race detector and coverage
go test ./... -race -cover

# A single package
go test ./internal/server/...
```

## Test structure

- **Unit tests** live next to the code they test (`_test.go` in the same
  package). Examples: `internal/server/errors_test.go`.
- **Integration tests** that need real infrastructure (PostgreSQL, Redis,
  HTTP server) can live here under `test/`, using `testing.T` and real
  connections (see `docker compose up` to start dependencies).

## No business tests

Because this is a template with no real business implementation, there are no
business tests. When you add a feature module (e.g. `internal/modules/payment`),
add unit tests for its `service.go` and `handler.go` using the same patterns.

## Example: adding a service unit test

```go
func TestPaymentService_CalculateTotal(t *testing.T) {
    svc := payment.NewService(fakeRepo, slog.Default())
    total, err := svc.CalculateTotal(10, 2)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if total != 20 {
        t.Fatalf("expected 20, got %d", total)
    }
}
```