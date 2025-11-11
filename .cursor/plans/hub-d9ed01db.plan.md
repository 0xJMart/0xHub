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

## Phase 1 – Domain & API Contracts

- Outcomes: stable Hub domain model, REST contracts, and shared schemas enabling modular integrations.
- Actionable Steps:

1. Draft ERD/domain model for entities (Project, Category, Tag, Link, MediaAsset, AuditLog, IntegrationCapability); align with ADR documenting domain boundaries and bounded contexts.
2. Select migration tool (Goose/Atlas); scaffold migration directory, baseline schema, and automated migration runner integrated with Makefile and CI; define `sqlc.yaml` pointing at the migration sources, wire `make generate`/`go generate ./...` to rebuild query code, and add CI drift checks to ensure generated files stay in sync.
3. Define versioned OpenAPI spec (`deploy/api/openapi.yaml`) covering registry endpoints (v1/projects, filters, media upload URLs, health) with tags, pagination, error model, and API versioning strategy; adopt API style guide.
4. Implement domain structs, validators, use cases in Go under `internal/domain` and `internal/app` using ports/adapters pattern to support future modularization.
5. Create repository interfaces for Postgres, Redis, MinIO within `internal/ports` and stub adapters in `internal/adapters` (e.g., `postgres`, `redis`, `minio` packages) with dependency injection wiring.
6. Implement Gin HTTP handlers with request validation (`go-playground/validator`), service-layer calls, consistent response envelopes, and correlation IDs; ensure error handling middleware for standardized responses.
7. Add input sanitization utilities, slug generation, filtering logic, and query builders to support modular search scopes; define feature flag hooks for future extensions.
8. Generate TypeScript API client/types using `openapi-typescript` (or `orval`) into `packages/shared/api`, publish via local npm workspace, and ensure versioning matches API releases.
9. Seed sample project data and fixtures via SQL migrations or Go seed scripts; include synthetic data for integration testing and demos.

- Testing: Go unit tests for domain services and repositories with in-memory adapters, `schemathesis` or `dredd` contract validation against running API, migration tests via `testcontainers` Postgres, OpenAPI snapshot diff tests in CI, TypeScript client build and smoke tests.

## Phase 2 – UI Shell & Navigation

- Outcomes: responsive, modular Hub UI supporting future micro-frontends and service discovery.
- Actionable Steps:

1. Configure design tokens (Tailwind theme JSON) and shared UI component library under `packages/shared/ui`; integrate Radix UI primitives and theming support (dark/light, custom brand overrides).
2. Build foundational layout components (AppShell, Header, Sidebar, NavigationRail) with responsive breakpoints and keyboard accessibility; support module-level navigation configuration via JSON manifest.
3. Implement project catalog views with filters, tags, search, sort, and pagination using `react-query` and server-driven query params; ensure offline caching strategies.
4. Create project detail view including metadata, markdown description rendering, media galleries (carousel/lightbox), integration badges, and quick launch buttons.
5. Implement micro-frontend placeholder module: iframe container with postMessage bridge, explicit `sandbox` attributes/allowlisted permissions, origin allowlisting for postMessage, dynamic manifest loader for downstream services, and feature flag gating.
6. Integrate MinIO media fetch/display using signed URLs, image optimization (e.g., `sharp` optional service), lazy loading, fallback placeholders, and upload progress UI.
7. Add global UX utilities: toast/notification center, skeleton loaders, error boundaries, global loading indicators, and offline mode banners.
8. Configure React Router with nested layouts, data loaders, and route-based code splitting; persist filter state via URL params and local storage.
9. Set up Storybook with MSW for API mocking, design token documentation, accessibility panel, and Chromatic/Loki pipeline for visual regression; document UI component usage.

- Testing: Vitest/RTL tests for components, hooks, and navigation flows; Storybook visual regression baselines; accessibility testing via `axe` and `storybook-addon-a11y`; Playwright smoke tests against mocked API; responsive layout checks via Chromatic per viewport.

## Phase 3 – Service Integration & Observability

- Outcomes: production-ready backend with persistence, caching, telemetry, and modular configuration controls.
- Actionable Steps:

1. Implement Postgres repositories with `sqlc` or `squirrel`, ensuring connection pooling (`pgxpool`), query tracing, and repository unit tests; support soft delete/archive operations.
2. Build Redis caching layer for project list/detail with TTL strategy, cache versioning keys, and invalidation hooks triggered by domain events; consider pub/sub for cache busting.
3. Finish MinIO integration: presigned URL generation, upload validation (size, MIME), media metadata storage, background cleanup jobs (cron or task queue) for orphaned assets.
4. Introduce background worker mechanism (Go workers or lightweight queue) for async tasks (e.g., screenshot refresh, analytics aggregation); encapsulate via `internal/app/tasks`.
5. Add structured logging using `zap` with request/trace IDs; configure log sampling and redaction policies for sensitive fields.
6. Integrate OpenTelemetry with `otelgin`, `otelhttp`, `otelsql`; export metrics/traces to Prometheus/Tempo via OTLP; define dashboards and tracing exemplars.
7. Expose `/healthz`, `/readyz`, `/metrics` endpoints; wire liveness/readiness probes in Helm values; add self-diagnostics endpoint summarizing config state.
8. Externalize configuration: environment variables, Helm values, and config maps with sane defaults; document configuration matrix and layering strategy (base + env overrides + secrets).
9. Implement feature flag system (e.g., open-source `flagd` integration or simple config-based toggles) to enable modular rollout of features.

- Testing: integration tests via `testcontainers-go` covering DB/cache/media flows; fuzz tests for API handlers; load/perf tests with `k6` automated nightly; observability smoke tests verifying metrics/traces/log formats; feature flag unit tests ensuring defaults.

## Phase 4 – Security & Access Control

- Outcomes: secure authn/z, hardened services, and governance around secrets and compliance.
- Actionable Steps:

1. Configure Keycloak realm/client; document roles, scopes, and OIDC discovery URLs; provide automation script for realm import/export.
2. Implement OIDC login with PKCE and silent token refresh in React; secure token storage (memory + refresh rotation); add `AuthProvider` context and guard components.
3. Develop Gin middleware verifying JWTs, extracting claims/roles, enforcing scopes; support API key/service account path for automation.
4. Define authorization policies in centralized policy module (e.g., `internal/app/policy`) with support for future ABAC rules; enforce per-endpoint and per-service logic.
5. Add audit logging to Postgres for privileged operations (project CRUD, media changes); include correlation IDs and user context for compliance.
6. Harden HTTP responses: set security headers (CSP, HSTS, X-Frame-Options), rate limiting (redis backed), request size limits, and body parsers.
7. Manage secrets via Kubernetes Sealed Secrets or External Secrets; ensure Helm templates reference sealed secret resources; document rotation procedure.
8. Integrate container scanning (Trivy), Go security scanning (`gosec`, `govulncheck`), JS dependency scanning (`npm audit`, `yarn audit`), and SAST tooling; schedule regular reports.
9. Establish incident response playbooks and access review cadence; integrate with centralized IAM if available.

- Testing: auth integration tests with valid/invalid tokens, RBAC policy tests, static analysis pipeline, OWASP ZAP baseline scan, dependency scan gating, manual security review checklist, automated headless browser tests verifying CSP enforcement.

## Phase 5 – GitOps Deployment & Cluster Validation

- Outcomes: Flux-managed continuous delivery with modular configuration per environment.
- Actionable Steps:

1. Enhance Helm chart with modular values for resources, HPA, pod disruption budgets, service monitors, ingress, RBAC, feature flags, and config map templating; add `values.dev.yaml`, `values.staging.yaml`, `values.prod.yaml` overlays.
2. Author Flux `HelmRelease`, `Kustomization`, and `ImagePolicy/ImageRepository` manifests in existing GitOps repo; align namespaces and modular structure with other services.
3. Configure CI/CD pipeline to build and push signed images (Cosign) with SBOM (Syft); update Flux image automation for semantic version tags and promotion workflows.
4. Create environment overlays via `kustomize` or `helmfile`, ensuring separate secrets, ingress domains, and resource profiles; document environment promotion process.
5. Deploy in-cluster Postgres/Redis/MinIO using Helm (Bitnami/operator charts) with persistent volume claims, replication, monitoring, and backup solutions (pgBackRest, Velero, MinIO lifecycle policies).
6. Automate preview environments per PR by templating Flux `Kustomization` with unique namespaces and ephemeral secrets; include teardown automation post-merge.
7. Document GitOps workflow, branch strategy, Flux reconciliation cadence, and troubleshooting guide; add guardrails for manual overrides (e.g., pause annotations).
8. Implement disaster recovery simulations: backup restore drills, Flux drift correction tests, network partition scenarios; capture runbooks.

- Testing: Helm unit tests (`helm unittest`), chart install tests (`ct install`) against kind cluster, Flux `kustomize build` validation in CI, automated smoke deployment tests (Argo Rollouts or custom scripts), Playwright E2E suite against staging namespace, canary rollout verification with synthetic monitoring.

## Phase 6 – Launch Readiness & Documentation

- Outcomes: operational handbook, governance, and release candidate ready for production adoption.
- Actionable Steps:

1. Finalize ADRs for architecture decisions (auth, storage, observability, feature flags, micro-frontend approach) and publish to `docs/adr`.
2. Produce API docs (Redoc/Stoplight) and SDK usage guides; integrate into frontend `/docs` route or dedicated docs site; include changelog/versioning strategy.
3. Create developer onboarding guide detailing setup, workflows, code standards, module ownership, and debugging tips; include checklists for new contributors.
4. Build operations runbooks: incident response, alert catalog, on-call rotation, SLO/SLA definitions, dashboard references (Grafana panels) and log search patterns.
5. Integrate uptime/status monitoring (Better Uptime/Pingdom) hitting `/healthz` and synthetic transaction monitors; configure alerting to messaging channels.
6. Conduct chaos engineering exercises (pod kill, node drain, dependency outage simulations) and document remediation steps; ensure alerting works end-to-end.
7. Prepare go-live materials: marketing assets, walkthrough videos, release notes, stakeholder demo script; align with downstream service owners.
8. Define onboarding checklist and contract for downstream services registering with Hub (service manifest, auth scopes, health check requirements) and publish to `docs/consumers`.
9. Establish post-launch roadmap and governance cadence (architecture review board, quarterly tech debt review).

- Testing: full regression suite (unit/integration/e2e) prior to launch, chaos drill verification reports, documentation peer review, release readiness checklist sign-off with stakeholders.

## Continuous Quality & Maintenance

- Maintain layered CI pipelines covering lint/unit/integration/e2e, nightly performance runs, weekly security scans, monthly resilience drills.
- Track code coverage thresholds (`go test -cover`, `vitest --coverage`) and enforce via CI gates; monitor mutation testing options for critical modules.
- Automate dependency updates (Renovate/Dependabot) with required test pipeline execution and changelog review; schedule quarterly library review.
- Monitor observability dashboards and error budgets; integrate analytics to inform backlog prioritization.
- Review ADRs and architecture quarterly, ensuring modularity remains intact and future services can integrate without architectural rewrites.