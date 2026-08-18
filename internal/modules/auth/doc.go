// Package auth is a minimal, commented ARCHITECTURAL EXAMPLE.
//
// It is NOT a real authentication implementation. It exists to demonstrate
// the module structure future features should follow:
//
//	internal/modules/<feature>/
//	├── handler.go    HTTP/Fiber layer
//	├── service.go    business/application logic
//	├── repository.go persistence abstraction
//	├── model.go      domain models
//	├── dto.go        transport DTOs
//	└── errors.go     module errors
//
// Dependency direction:
//
//	HTTP -> Handler -> Service -> Repository interface -> implementation -> DB
//
// Business code here must NOT depend on Fiber, PostgreSQL driver details,
// Redis implementation details, OpenTelemetry exporters, or environment
// variables. Those are injected.
package auth
