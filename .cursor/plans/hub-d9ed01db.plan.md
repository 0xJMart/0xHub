<!-- d9ed01db-af8a-484a-8667-da9148e1a808 1325bd88-41c5-4577-abfd-b6ed404fba4a -->
# Hub Implementation Plan (Option B)

## Phase 0 – Foundation & Tooling ✅ Completed

- Outcomes: repository structure, coding standards, CI guardrails, and baseline architecture governance established.
- Summary:
- Mono-repo skeleton in place with workspace documentation.
- Go hub-api service scaffolded with clean architecture layout, Makefile targets, and Go tooling.
- Frontend workspace bootstrapped (Vite + React + TypeScript) with Tailwind, ESLint flat config, Prettier (Tailwind plugin), Storybook (Vite builder), Vitest, and absolute import aliases.
- Shared configs delivered: `.editorconfig`, `.gitignore`, `.gitattributes`, Husky + lint-staged hooks, Dependabot updates.
- Dockerfiles and docker-compose stacks (dev/ci) for API, frontend, Postgres, Redis, MinIO, Keycloak; Makefile shortcuts; Keycloak realm import.
- Helm chart `deploy/helm/hub` with helper templates, deployments/services, ingress, ServiceMonitor, and environment-specific values.
- CI workflows covering Go lint/test, frontend lint/typecheck/test/build, container builds, Trivy scans, Helm lint.
- ADR template plus ADR-0001 capturing Phase 0 foundation decisions.
- Validation: `go test ./...`, `npm run lint`, `npm run typecheck`, `npm run test`, and `helm lint deploy/helm/hub` all pass locally (Playwright browsers installed).

- Testing checklist (achieved):
- `go test ./...`
- `npm run lint`
- `npm run typecheck`
- `npm run test` (Vitest + Storybook + Playwright)
- `helm lint deploy/helm/hub`

## Phase 1 – Core Domain & Schema ✅ Completed

- Outcomes: lightweight domain model and database shape ready for a single homelab deployment.
- Summary:
- ADR-0002 records the ERD decisions and trade-offs tailored to homelab constraints.
- Goose-backed migrations enable `pgcrypto`, create baseline tables, and seed demo projects/tags/media.
- Hub API Makefile exposes `migrate-*` targets that proxy the Goose CLI via `go run`.
- OpenAPI spec (`deploy/api/openapi.yaml`) documents project, tag, media upload, and health endpoints.
- Domain structs gained validation helpers with unit tests covering nested entities.

- Testing checklist (achieved):
- `cd services/hub-api && go test ./...`
- Goose migrations apply cleanly against a local Postgres instance with `pgcrypto`.

## Phase 2 – Service Layer & HTTP API ✅ Completed

- Outcomes: functional Go service exposing REST endpoints backed by the domain logic.
- Summary:
  - Repository ports defined with a Postgres adapter that hydrates categories, tags, links, and media.
  - Application services encapsulate project CRUD, tag listing, and readiness checks with structured logging.
  - Gin HTTP server delivers `/healthz`, `/readyz`, `/projects`, `/projects/{slug}`, and `/tags` with consistent error payloads.
  - Shared TypeScript client published at `packages/shared/api` (`@0xhub/api`) mirroring the OpenAPI contract.
  - Make targets `make test` and `make test-api` run Vitest Storybook checks plus Go handler/service smoke tests.

- Validation:
  - `npm run lint`, `npm run typecheck`, `npm run test`
  - `go test ./...`, `make test-api`

## Phase 3 – UI Shell & Core Views

- Outcomes: responsive but minimal React UI suited for a homelab dashboard.
- Actionable Steps:

1. Finalize base Tailwind theme tokens and a small shared component set (`AppShell`, `PageHeader`, `ProjectCard`) inside `packages/shared/ui`.
2. Wire the frontend to the API client with `react-query` (or simple fetch hooks) for project list/detail screens; support filtering and search with URL params.
3. Implement a project detail page with markdown description rendering and lightweight media gallery support using presigned URLs (no advanced image processing).
4. Add simple UX helpers: toast notifications, skeleton loading states, and a global error boundary.
5. Configure Storybook for component previews with MSW mocks, but keep the setup lean (no Chromatic or visual diff tooling).

- Testing: rely on Vitest + React Testing Library for component/unit tests and one or two Playwright smoke specs to verify basic flows.

## Phase 4 – Storage & Background Tasks

- Outcomes: reliable persistence and optional background work without enterprise observability overhead.
- Actionable Steps:

1. Finish the Postgres adapter with straightforward SQL (using `sqlc` only if it keeps code cleaner); ensure connection pooling with sensible defaults.
2. Integrate MinIO (or plain file-system storage) for media assets, generating presigned URLs directly from the API.
3. Provide a simple background worker (Go goroutine or cron-like job) for periodic clean-up of orphaned media files.
4. Configure structured logging using the standard library or `zap` in development-friendly mode; capture request IDs if easy.
5. Offer a basic `/metrics` endpoint only if Prometheus scraping is already present in the homelab; otherwise document manual health checks.

- Testing: integration tests using `testcontainers-go` to verify DB/media round-trips plus manual verification of worker jobs in Docker Compose.

## Phase 5 – Lightweight Auth & Hardening

- Outcomes: sensible security for a private homelab without enterprise IAM complexity.
- Actionable Steps:

1. Start with a simple local auth strategy (shared admin password or Keycloak if already running) and guard routes accordingly.
2. Add middleware for JWT/session verification, rate limiting (middleware-level), and secure headers via Gin.
3. Implement basic audit logging to Postgres for project mutations, tagged with user id and timestamp.
4. Manage secrets via `.env` files checked into `deploy/docker-compose/` templates; document how to override them securely on the host.
5. Run `gosec`/`govulncheck` and `npm audit` during CI to catch obvious issues without layering additional scanners.

- Testing: auth middleware unit tests, a Playwright spec to confirm protected routes, and manual threat review focused on homelab exposure.

## Phase 6 – Deployment & Documentation

- Outcomes: easy-to-operate homelab deployment with clear notes for future you.
- Actionable Steps:

1. Maintain Docker Compose definitions for API, frontend, Postgres, and MinIO; ensure `make up` brings the stack online with seeded data.
2. Write concise operations docs in `docs/` covering setup, backups (pg_dump + MinIO sync), upgrade steps, and restore drills.
3. Capture ADRs only for significant choices (e.g., auth approach, storage layout) to avoid churn.
4. Provide a changelog template and release checklist focused on testing and data backups before upgrades.
5. Schedule a quarterly self-review: update dependencies, rerun tests, and archive stale projects.

- Testing: run `go test ./...`, `npm run lint`, `npm run test`, and `docker compose up` in CI or before manual releases to confirm the stack stays healthy.

## Continuous Quality (Homelab Scale)

- Keep CI limited to linting, unit tests, and type checks; run integration tests before major changes.
- Dependabot/Renovate may be optional—document how to update dependencies manually if automation is overkill.
- Use existing homelab monitoring (if any) for uptime; otherwise rely on health endpoints and manual checks.
- Revisit the plan yearly to see if new services or integrations warrant expanding beyond the homelab scope.