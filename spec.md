### IHB Transport API — Engineering Specification

### Document control
- **Project**: IHB Transport (Go + Gin + PostgreSQL/Supabase)
- **Purpose**: Define engineering requirements and standards with emphasis on **reliability**, **scalability**, **maintainability**, and **security**.
- **Audience**: Backend engineers, DevOps/SRE, security reviewers.
- **Status**: Living document; update alongside architectural changes.

### Scope
- **In scope**: HTTP API, authentication/authorization, database schema & migrations, operational concerns (logging, rate limiting, error handling), container deployment (Render/Railway/Fly.io), Supabase Postgres connectivity.
- **Out of scope**: Frontend/UI, mobile apps, payments, analytics, customer support tooling.

### System overview
- **Runtime**: Go API server using Gin.
- **Persistence**: PostgreSQL (Supabase recommended).
- **Auth**: JWT for protected routes; role-based access control.
- **Core domain**: Delivery requests lifecycle and notifications.

### Architecture
- **Layers**
  - **Handlers (`internal/handlers`)**: HTTP transport, request binding/validation, response formatting.
  - **Services (`internal/services`)**: Business logic, state transitions, domain validation.
  - **Repositories (`internal/repository`)**: Persistence operations (GORM).
  - **Auth (`internal/auth`)**: JWT generation/validation, auth/role middleware.
  - **Database (`internal/database`)**: Connection management.
  - **Migrations (`migrations`, `cmd/migrate`)**: Schema changes.
- **Design principles**
  - **Dependency inversion**: Handlers depend on service interfaces; services depend on repository interfaces.
  - **Side-effects centralized**: Emails, logs, and persistence are invoked from services, not scattered.
  - **Idempotence**: Migrations and state transitions should be safe to retry.

### API surface (high level)
- **Auth**
  - `POST /login` — obtains JWT (admin in current implementation)
- **Public**
  - `POST /api/public/deliveries` — create delivery request
  - `GET /api/public/deliveries/:id` — track by ID
  - `GET /api/public/deliveries/track?email=` — track by email
  - accept/decline endpoints for quotes (POST + GET variants)
- **Protected (JWT)**
  - Admin: `GET /api/admin/deliveries`, `GET /api/admin/deliveries/status`, `POST /api/admin/deliveries/:id/price`
  - Driver: pickup/complete status endpoints and listing

### Data model (summary)
- **Admin**
  - `email` unique; `password_hash` stored (bcrypt).
- **DeliveryRequest**
  - Client details, pickup/dropoff fields, item/service fields, `status`, `price`, timestamps.
- **StatusLog / EmailLog**
  - Audit trail for status changes and email events.

### Reliability requirements
- **Service availability**
  - **MVP target**: 99.5% monthly.
  - **Future target**: 99.9% monthly (requires redundant hosting + stronger DB/observability posture).
- **Startup behavior**
  - Fail fast on missing critical configuration (JWT secret, admin seed config when required, DB URL for production).
  - On DB connectivity failure, exit non-zero so the platform restarts the container.
- **Migrations**
  - Must be **repeatable** and **safe to run multiple times** (no destructive changes without a migration plan).
  - Prefer separate migration job/step in production; if migrations run at boot, they must be guarded to avoid concurrent execution across replicas.
- **Request handling**
  - Return consistent, machine-readable error shapes.
  - Enforce request size limits.
  - Apply rate limiting to protect availability.
- **Timeouts**
  - Set reasonable server timeouts (read, write, idle) and upstream timeouts for external services (email).
- **Health**
  - Provide a lightweight health endpoint (e.g., `/ping`) for platform health checks.
  - Future: add `/healthz` with dependency checks (DB ping with short timeout).
- **Error budget**
  - Track elevated 5xx rates, timeouts, and DB connection errors; treat sustained breaches as incidents.

### Scalability requirements
- **Horizontal scaling**
  - The API must support running multiple replicas behind a load balancer.
  - Server instances should be **stateless** with respect to user sessions (JWT is stateless).
- **Database connections**
  - Use connection pooling; set max open/idle connections per replica to avoid exceeding Supabase limits.
  - Prefer Supabase pooler for environments that create many short-lived connections (optional for traditional servers, useful for serverless).
- **Query performance**
  - Ensure indexes exist for common filters (delivery by ID, by email, by status, and logs by delivery ID).
  - Avoid N+1 query patterns when loading deliveries with logs; use preloading selectively.
- **Rate limiting & abuse resistance**
  - Public endpoints must be rate limited per IP / per key if introduced.
  - Future: add per-email throttling for tracking endpoints to reduce scraping/abuse.
- **Workload separation (future)**
  - Background tasks (email sending, notifications) should move to a queue/worker model when volume grows.

### Maintainability requirements
- **Code organization**
  - Keep handlers thin; put validation and state transitions in services.
  - Keep repository methods focused and testable (CRUD + query helpers).
- **Interfaces**
  - Maintain clear service/repository interfaces to enable mocking in tests.
- **Configuration**
  - All runtime configuration via environment variables; no secrets in git.
  - Provide documented examples in `QUICKSTART.md` / `DEPLOYMENT.md`.
- **Logging**
  - Use structured logs where possible; include request ID/correlation ID (future).
  - Never log secrets (JWT secret, DB URL, passwords).
- **Error handling**
  - Centralize error formatting; avoid leaking internals in responses.
  - Use typed errors in services to map to HTTP status codes deterministically.
- **Testing strategy**
  - Unit tests: services validation and state transitions.
  - Integration tests: repository layer with ephemeral Postgres (local Docker).
  - Contract tests: basic endpoint smoke tests (can re-use existing shell scripts).
- **Documentation**
  - Keep endpoint docs in `DEPLOYMENT.md` or an OpenAPI spec (future).
  - Keep operational runbooks for common failures (DB unreachable, SMTP failures).

### Security requirements
- **Authentication**
  - JWT tokens must be signed with a strong secret (>= 32 bytes random).
  - Token validation must enforce signature and expiration.
  - Use short-lived access tokens (future: refresh tokens if needed).
- **Authorization**
  - Enforce role checks on admin endpoints.
  - Driver endpoints should enforce driver role once driver auth is implemented (currently driver group is protected but role enforcement may be added).
- **Transport security**
  - Public deployment must use HTTPS at the edge (handled by platform).
  - DB connections to Supabase must use SSL/TLS.
- **Input validation**
  - Validate required fields, email format, max lengths.
  - Treat all user input as untrusted.
- **Injection defenses**
  - Use parameterized queries (GORM) to avoid SQL injection.
  - Avoid rendering unescaped user input in HTML responses; escape any interpolated values.
- **Secrets management**
  - Store secrets only in the hosting platform’s secret store / environment variables.
  - Rotate secrets on compromise; support rolling JWT secret rotation plan (future).
- **CORS**
  - Default CORS should be restrictive in production (allow only trusted origins).
- **Rate limiting**
  - Public endpoints must be protected from brute force and abuse.
  - Login endpoint should have stricter limits and potential lockout policy (future).
- **Auditability**
  - Persist status change logs with actor and timestamp.
  - Persist email event logs with delivery ID and status.

### Configuration (environment variables)
- **Required (production)**
  - `DATABASE_URL` — Postgres connection string (Supabase)
  - `JWT_SECRET` — secret for signing tokens
  - `GIN_MODE=release`
  - `PORT` — platform-provided or default 8080
- **Operational/admin**
  - `ADMIN_EMAIL`, `ADMIN_PASSWORD` — used for seeding initial admin (ensure strong password)
- **Optional**
  - SMTP provider variables (if email sending is enabled beyond logging)

### Deployment standards
- **Build**
  - Container image builds must be reproducible; pinned Go version recommended.
- **Runtime**
  - Run as non-root user in container (future hardening).
  - Configure health checks to hit `/ping`.
- **Migrations**
  - Prefer a one-off migration job in production pipelines; avoid concurrent auto-migration across replicas.
- **Rollouts**
  - Use rolling deploys; ensure backward compatible DB changes or staged migrations.

### Observability
- **Logs**
  - Request logs include method, path, status, latency.
  - Error logs include stable error codes and minimal safe context.
- **Metrics (future)**
  - Request count, latency percentiles, error rates.
  - DB connection pool metrics.
  - Rate limiter drops.
- **Tracing (future)**
  - Propagate trace/request IDs; integrate OpenTelemetry when needed.

### Non-functional acceptance criteria
- **Reliability**
  - Service returns 200 from `/ping` within 2 seconds under normal conditions.
  - DB connection failures surface clearly and terminate the process for restart.
- **Scalability**
  - Multiple replicas can run without shared local state.
  - DB connection limits are respected under expected concurrency.
- **Maintainability**
  - Business logic changes are isolated to service layer with unit tests.
  - New endpoints follow existing patterns and error conventions.
- **Security**
  - Protected endpoints reject requests without `Authorization: Bearer <token>`.
  - Secrets are not present in repository history and are not logged.

### Backlog (recommended next improvements)
- **Add OpenAPI spec** (`openapi.yaml`) for API contract and client generation.
- **Add CI**: lint, `go test`, static analysis (govulncheck), dependency scanning.
- **Add DB connection pool config** (max open/idle, lifetime) and document recommended values for Supabase.
- **Harden auth**: token expiry, refresh flow (if needed), login rate limiting, audit logs.
- **Separate migration execution** from server boot in production.

