# Go Fiber Production Template

A production-ready **Go + Fiber** backend template designed as a
**feature-based modular monolith**. It is a starting point for new projects:
infrastructure, bootstrapping, middleware, observability, configuration,
database, Redis, logging, testing, Docker, Swagger and CI/CD are all prepared.
You add the business features.

> **Infrastructure is ready before business code exists.**
>
> **Observability observes the application, it does not reshape it.**

---

## 1. What this is

A reusable skeleton, **not** a complete application. The only feature module
present is `internal/modules/auth`, and it is a **minimal commented
architectural example**, not a real authentication implementation. You will
implement real features yourself.

## 2. Architecture philosophy

Modular monolith, single binary, DDD-like feature boundaries.

```
                    Go Application
                         │
              ┌──────────┴──────────┐
              │                     │
        Infrastructure          Business
              │                     │
       ┌──────┼──────┐        ┌─────┴─────┐
       │      │      │        │           │
      DB    Redis   OTel     Auth       Future
       │      │      │       Module     Modules
```

- **Feature-oriented**, not technical-layer-oriented. Each feature owns its
  handler/service/repository/model/dto/errors.
- Business modules depend on **abstractions**, not on Fiber internals,
  PostgreSQL/Redis driver details, OpenTelemetry exporters, or environment
  variables.
- Infrastructure is **injected** through explicit constructors. No DI framework.

## 3. Directory structure

```
.
├── cmd/api/main.go              Bootstrapping entry point (small)
├── internal/
│   ├── app/                     Composition: deps, routes, app
│   ├── config/                  Env-driven configuration
│   ├── infra/                   Infrastructure
│   │   ├── cache/redis.go       Redis
│   │   ├── database/            PostgreSQL + migrations
│   │   ├── email/               SMTP sender + no-op
│   │   ├── gateways/            Third-party API boundaries
│   │   ├── logger/              slog-based structured logging
│   │   ├── observability/       OpenTelemetry (optional)
│   │   ├── sms/                 SMS sender + no-op
│   │   ├── storage/             Object storage boundary
│   │   ├── swagger/             OpenAPI UI wiring
│   │   └── token/               JWT manager
│   ├── modules/auth/            Example feature module
│   ├── server/                  Response/request/error/middleware/validator
│   └── shared/                  Cross-module concepts (errors, pagination…)
├── migrations/                  golang-migrate SQL files
├── docs/                        Generated OpenAPI spec (swag)
├── pkg/                         Reserved for external libraries
├── test/                        Integration test space
└── Dockerfile, docker-compose.yml, Makefile, .github/workflows/
```

## 4. Run locally

Prerequisites: Go 1.25+, PostgreSQL, Redis.

```bash
cp .env.example .env       # adjust values
make run                   # or: go run ./cmd/api
```

Endpoints:
- `GET /health` — liveness (process alive)
- `GET /ready` — readiness (DB + Redis reachable)
- `GET /swagger/` — OpenAPI UI (if `SWAGGER_ENABLED=true`)

## 5. Environment configuration

Configuration is read from environment variables (see `.env.example`). A local
`.env` is loaded automatically. Sections: Application, Server, Database,
Redis, Logging, OpenTelemetry, Security, Rate limiting, External services,
Swagger. **Never commit real secrets.**

## 6. Docker

```bash
docker compose up --build
```

Starts the API, PostgreSQL, Redis, and optionally the OpenTelemetry Collector
and Jaeger (uncomment in `docker-compose.yml`). The app runs with `OTEL_ENABLED=false`
and requires no collector to start.

Build a single image:

```bash
docker build -t go-fiber-template .
```

## 7. Database migrations

Uses [golang-migrate](https://github.com/golang-migrate/migrate) with the file
source pointing at `migrations/`. Format:

```
migrations/
  000002_my_feature.up.sql
  000002_my_feature.down.sql
```

Migrations run at startup when `DB_RUN_MIGRATIONS=true`. Add feature tables in
your own migrations as needed.

## 8. Redis

Initialized in `internal/infra/cache` with pooling, health check and graceful
close. Exposed via `Deps.Redis` for modules to use behind their own
abstractions.

## 9. OpenTelemetry

Optional. When `OTEL_ENABLED=false` the app runs with no-op providers and no
collector is required. When enabled, it exports OTLP traces and metrics:

```env
OTEL_ENABLED=true
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
OTEL_SERVICE_NAME=my-service
OTEL_SERVICE_VERSION=1.0.0
OTEL_ENVIRONMENT=development
OTEL_INSECURE=true
OTEL_SAMPLING_RATIO=1.0
```

Instrumentation lives at the **infrastructure layer** (HTTP middleware, DB/Redis
drivers). **Business code does not call the tracer.** This keeps handlers clean.

## 10. Health checks

- `GET /health` — liveness: the process is alive.
- `GET /ready` — readiness: checks PostgreSQL and Redis.

No business logic lives in these endpoints.

## 11. Create a new module

For a new feature, e.g. `payment`, create `internal/modules/payment/`:

```
internal/modules/payment/
├── handler.go       HTTP/Fiber layer
├── service.go       business/application logic
├── repository.go    persistence interface
├── model.go         domain models
├── dto.go           transport DTOs
└── errors.go        module errors
```

**Steps:**

1. Define the domain model in `model.go` and DTOs in `dto.go`.
2. Define the persistence **interface** in `repository.go`.
3. Implement business rules in `service.go` (inject dependencies via the
   constructor; no Fiber/DB/OTel specifics).
4. Implement the HTTP layer in `handler.go` (the only place using Fiber).
5. Add a concrete repository adapter (e.g. in `internal/app/payment_repo.go`)
   that implements the interface against PostgreSQL.
6. Wire it in `internal/app/modules.go` (`newModules`) and mount routes in
   `internal/app/routes.go`.
7. Add Swagger annotations and run `make swagger`.
8. Add unit tests next to `service.go` and `handler.go`.

## 12. Dependency rules

```
HTTP
 ↓
Handler        (Fiber only here)
 ↓
Service        (business rules, no framework)
 ↓
Repository interface
 ↓
Infrastructure implementation
 ↓
PostgreSQL
```

- Modules never import Fiber, pgx, go-redis, OTel exporters, or env config.
- Infra is passed in through constructors.
- `internal/shared` holds only genuinely cross-module concepts — don't turn it
  into a dumping ground.

## 13. Testing

- Unit tests live next to the code (`_test.go`).
- Integration tests go in `test/` against real infra (`docker compose up`).
- See `test/README.md`.

```bash
make test    # go test ./... -race -cover
```

## 14. Production considerations

- Set a strong `JWT_SECRET` (required in production).
- Use `LOG_JSON=true`, `APP_ENV=production`.
- Enable `OTEL_ENABLED=true` and point at a collector/backend.
- Configure a real email/storage provider (the no-op senders are placeholders).
- Set a real `CORS_ALLOWED_ORIGINS`.
- Adjust rate-limit values.
- CI/CD is preconfigured in `.github/workflows` (CI lint/build/test + CD image
  push to GHCR on tags).

## License

Template for your own use. Modify freely.