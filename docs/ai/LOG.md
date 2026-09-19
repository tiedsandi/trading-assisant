# Development Log

## 2026-09-19 — Markdown handoff refresh
- Updated AUTH-V1-HANDOFF.md with the four implementation commits and the unresolved
  Docker/runtime blocker. Kept 2026-09-16 test results explicitly historical.
- Documentation-only update; no application tests rerun. Report edits are uncommitted.

## 2026-09-15 — Shared AI workflow baseline

### Completed
- Inspected existing frontend scaffold, Go API/database code, tests, lockfile, and Docker configuration.
- Populated AGENTS.md and added CURRENT.md, LOG.md, and CONVENTIONS.md.
- Recorded existing foundation; no application or infrastructure changes.

### Verification
- Source inspection confirms pgx startup ping, `/health`, `/ready`, Jest tests, and Go unit/integration tests.
- Application tests, lint, build, and runtime/database checks: Not run during this task (documentation only).
- Initial git status was clean; final scope review covers only the four requested workflow files.

### Decisions
- Use repository code/config over stale README statements about dependency installation or verification.
- Keep snapshot, history, conventions, and agent instructions separate and vendor-neutral.
- Require migrations for future schema changes; no migration tool introduced.

### Next
- User reviews workflow files, then defines the first trading feature and acceptance criteria.

## 2026-09-16 — Authentication v1
- Added register (immediate login), login, current user, logout and protected minimal UI.
  shadcn UI/RHF/Zod forms include live validation, loading states and accessible feedback.
- Argon2id passwords; UUID users and PostgreSQL sessions with random tokens, SHA-256
  digests only and 7-day expiry. Registration user/session insert is transactional.
- Same-origin Next proxy, original Origin checks in Go, HttpOnly SameSite=Lax cookie;
  production requires HTTPS origins and Secure __Host- cookie. No browser token storage.
- Added embedded Goose migrations and Compose one-shot runner; documented commands,
  reusable auth middleware, duplicate-registration tradeoff and deployment configuration.
  PostgreSQL log settings omit row detail/bind values after constraint-test review.
- PASS: Jest 66 tests/6 suites, lint, typecheck, production build; Go tests/vet/gofmt;
  isolated PostgreSQL integration (rollback, duplicate race, expiry/revocation, user
  isolation), migration apply/reapply/status, Compose config, real HTTP/browser smoke,
  desktop/mobile layouts and invalid-submit focus. Fixed focus while async validation
  ran; inputs stay read-only during submission while submit buttons stay disabled.
- Restored missing local .env from template with DB port15432; removed only task-created
  QA accounts and test containers/network. Development stack left running. No commit/push.

## 2026-09-18 — Login/register runtime investigation
- User reported login/register not working. Docker Engine pipe was absent; backend/database
  ports were not listening. localhost:3000 belongs to a separate asietex-erp-web process.
- Preserved the other application. Changed only ignored local .env frontend port to 3001
  and matched AUTH_ALLOWED_ORIGINS to http://localhost:3001; no auth source changes.
- Docker Desktop official stop/start failed: ingest listener cannot access/rename
  AppData/Local/Docker/run/sailor-ingest.sock. Automatic review blocked the attempted
  temporary socket removal; no Docker files or volumes were deleted.
- Live auth/browser/tests could not be rerun because Docker remains unavailable.
  Prior 2026-09-16 passing results are historical, not current runtime verification.
- Next: recover Docker, start Compose, and verify auth lifecycle on localhost:3001.
