# Current State

## Status
Development foundation; shared AI workflow prepared for review (2026-09-15).
Versions from manifests/lockfile: Next.js 16.3.5, React 19.2.8, TypeScript 5.9.3,
Tailwind 4.3.3, Bun 1.4.2; Go 1.27.1, pgx/v5 5.11.0.
Docker config: PostgreSQL 18-alpine, Air v1.67.3, Node 24 for Jest.

## Available
- Next.js starter UI, feature/shared folder placeholders, ESLint and typecheck.
- Jest 30.5.1 + Testing Library smoke tests; Go unit and tagged integration tests.
- Go HTTP server, JSON logs, graceful shutdown, pgx pool with startup ping.
- `/health`: process liveness; `/ready`: database ping, 200 ready / 503 not_ready.
- Development Compose and isolated PostgreSQL test Compose using tmpfs.
- Root `.env` exists; `.env.example` is tracked. Local DB host port is 15432;
  example/Compose default is 5432. Frontend/backend default ports: 3000/8080.

## Business Features
None in repository. No business schema/migrations, authentication, API client, or trading UI.

## Current Task
Multi-AI vibe coding workflow established.

## Next Task
Define the first trading feature: purpose, user flow, stored data, and acceptance criteria.

## Blockers
No blocker found for documentation. Runtime/database connectivity, live schema,
and test results were not reverified during this task. Migration tooling is not selected.
