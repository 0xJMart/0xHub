## 0xHub Monorepo

This repository houses the core services and UI for the Hub platform. The structure follows a modular mono-repo layout to keep shared contracts, tooling, and deployment assets aligned.

### Directory Map

- `frontend/` – Vite + React client application and shared UI tooling.
- `services/hub-api/` – Go service implementing Hub APIs using clean architecture.
- `packages/shared/` – Shared libraries (schemas, API clients, utilities) published across workspaces.
- `deploy/` – Helm charts, Kubernetes manifests, and Docker-related assets.
- `docs/` – Documentation, ADRs, and onboarding guides.

Each workspace has its own README with ownership details, local development commands, and module boundaries. See `docs/` for architecture decision records and operational runbooks as they are added.

### Getting Started

1. Install local toolchains (already vendored for convenience):
   - Go 1.22.5 at `/home/admin/go/bin`
   - Node.js 20.19.0 at `/home/admin/node-v20.19.0-linux-x64/bin`
2. Install npm dependencies for all workspaces:

   ```bash
   npm install
   ```

3. Copy `config/env.example` to `.env` (or another file) and tweak values as needed for local development.

4. Bring up the development stack (API, frontend, and backing services):

   ```bash
   make dev-up
   ```

   Tear the stack down with `make dev-down` or stream logs with `make dev-logs`.

### Tooling Shortcuts

- `make lint`, `make format`, `make typecheck`, `make test` proxy the respective npm workspace scripts.
- `make ci-up` / `make ci-down` provide a deterministic docker-compose environment for CI runs.

### Additional Resources

- `config/env.example` – reference environment variables for local/CI.
- `docker-compose.dev.yml` – full local stack (Postgres, Redis, MinIO, Keycloak, services).
- `docker-compose.ci.yml` – minimal stack optimised for CI pipelines.

