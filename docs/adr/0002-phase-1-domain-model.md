# ADR-0002: Phase 1 Domain Model

- **Status:** Accepted
- **Date:** 2025-11-11
- **Authors:** Hub Team
- **Context:** Define the initial homelab-focused domain model and supporting storage layout for Projects, Tags, Links, and Media Assets.

## Context

We need a lightweight schema, matching single-tenant homelab needs, that keeps projects and their related metadata easy to seed and query. The service must expose everything required by the Phase 1 API (project listing, detail views, tag filtering, and presigned media uploads) while remaining simple to operate locally. The model should map cleanly to Go structs and SQL migrations, avoiding premature abstractions.

Key constraints:

- Prioritize readability and manual SQL debugging over generics or heavy ORMs.
- Keep optional relationships nullable to support incremental enrichment.
- Prefer UUID primary keys to avoid ID collisions between seed data and future imports.
- Record timestamps and soft-delete markers to support future automation without schema churn.

## Decision

We will model the hub around four core tables with supporting join tables:

```
Categories (1) ────────┐
                       │
Projects (many) ───┐   │
                   │   │
                   └─< ProjectTags >─┐
                       ▲             │
                       │             ▼
                     Tags           Links
                       │
                       ▼
                   MediaAssets
```

- `projects`: core entity with title, slug, status, summary, rich description, optional category, hero media reference, and audit timestamps.
- `categories`: optional grouping with display name and ordering. Projects reference categories via `category_id`.
- `tags`: free-form labels normalized via a join table `project_tags` keyed by `project_id` and `tag_id`.
- `project_links`: curated outbound references (e.g., repo, demo) with type and URL.
- `media_assets`: uploaded artifacts tied to a project (or standalone, in preparation for drafts), storing path, MIME type, size, and optional descriptions.

Supporting rules:

- All tables share UUID primary keys generated in SQL (`gen_random_uuid()`) and maintain `created_at`/`updated_at` timestamps.
- `projects.slug` is unique to back the `/projects/{slug}` lookup.
- Tags are unique on `name` (case-insensitive) with friendly `display_name`.
- `media_assets` may reference a project, but can be associated post-upload via nullable `project_id`.
- We will seed a single category, a handful of tags, and two projects to support frontend demos.

## Consequences

- Positive impacts:
  - Clear separation between taxonomy (`categories`, `tags`) and project content keeps queries simple.
  - UUID keys avoid sequence management in Docker Compose and future distributed setups.
  - Seed data accelerates frontend smoke testing and Storybook examples.
- Negative impacts:
  - Requires `pgcrypto` extension (`gen_random_uuid`) to be enabled in Postgres.
  - Additional join table introduces minor complexity for tag queries.
- Follow-up tasks:
  - Add repository implementations in Phase 2 that mirror this schema.
  - Revisit media storage metadata once MinIO integration lands in Phase 4.

## Alternatives Considered

1. **Serial integer primary keys** – easier to inspect manually but risk collisions when seeding or importing data; would require additional care during backups.
2. **Single `project_assets` table combining links and media** – simplifies structure but couples heterogeneous data (URLs vs. files), complicating validation and API responses.

## References

- Phase 1 plan: `@hub-d9ed01db.plan.md`
- `deploy/api/openapi.yaml`


