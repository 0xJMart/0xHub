# ADR-0001: Platform Foundation (Phase 0)

- **Status:** Accepted
- **Date:** 2025-11-11
- **Authors:** Platform Team
- **Context:** Establish repository structure, tooling, and deployment fundamentals for the Hub platform.

## Context

We are launching the Hub as a mono-repo hosting the React frontend, Go API, shared packages, and deployment assets. Early alignment on tooling and structure unlocks parallel feature work in later phases while keeping operational requirements in sight.

## Decision

1. **Mono-repo layout:** Adopt a top-level structure with dedicated directories for `frontend/`, `services/hub-api/`, `packages/shared/`, `deploy/`, and `docs/`. Shared tooling (npm workspaces, Husky, lint-staged) lives at the repository root.
2. **Backend stack:** Implement the Hub API in Go using a clean-architecture layout (app/domain/ports/adapters), managed via Go modules, `golangci-lint`, and a distroless multi-stage Docker build.
3. **Frontend stack:** Build the UI with Vite + React + TypeScript, Tailwind, Vitest, and Storybook (Vite builder). Enforce formatting via Prettier (Tailwind plugin) and linting via ESLint flat config.
4. **Local development:** Provide docker-compose stacks for dev and CI parity that orchestrate Postgres, Redis, MinIO, Keycloak, and service containers. Include a root `Makefile` for common workflows.
5. **Deployment scaffolding:** Introduce a Helm chart (`deploy/helm/hub`) with environment overlays, plus GitHub Actions CI covering lint/test/build, Docker image builds, vulnerability scanning, and Helm lint.
6. **Architecture governance:** Track decisions using ADRs, starting with this baseline and a reusable template.

## Consequences

- ✅ Teams can start feature development immediately with consistent linters, type-checkers, Storybook, testing, and containerization.
- ✅ CI pipelines ship with guardrails (linting, testing, security scanning) before code reaches mainline.
- ✅ Helm chart and docker-compose overlays provide a launch point for GitOps integration and local parity.
- ⚠️ Mono-repo introduces tooling expectations (Node 20.19+, Go 1.22+) that contributors must meet; vendored toolchains mitigate friction.
- ⚠️ Docker-based local stack has a heavier footprint; future iterations may introduce lightweight mocks for quick dev flows.

## Alternatives Considered

1. **Poly-repo layout** – Rejected due to duplicated tooling, higher coordination cost, and slower feature iteration.
2. **Docker-only stack without Helm scaffolding** – Rejected; we need Kubernetes deployment parity early to integrate with GitOps and platform policies.

## References

- Phase 0 plan (`.cursor/plans/hub-d9ed01db.plan.md`)
- `.github/workflows/ci.yml`
- `deploy/helm/hub`

