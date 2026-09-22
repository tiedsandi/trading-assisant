# Current State

## Status
Authentication v1 implemented; full verification on 2026-09-16 remains historical.
On 2026-09-22, the owner confirmed Docker recovery and successful manual registration,
login, authenticated dashboard access, logout and route guards (user-reported evidence).
The previous Docker blocker and register/login failure report are resolved.
RISK-001 linear perpetual risk domain foundation is implemented locally and awaiting review.
Stack: Next.js 16.3.5, React 19.2.8, TypeScript 5.9.3, Tailwind 4.3.3, Bun 1.4.2;
Go 1.27.1, pgx/v5 5.11.0, PostgreSQL 18; Goose 3.28.0, argon2id 1.0.0.

## Available
- Register with immediate authentication, login, current user, logout/revocation.
- Email normalized/unique; passwords use Argon2id, 8–128 Unicode codepoints.
- PostgreSQL opaque sessions: random 32-byte tokens, only SHA-256 stored, 7 days.
- Next server auth checks protect `/dashboard`; guests use `/login` and `/register`.
- shadcn UI + React Hook Form + Zod, live validation, accessible errors, visibility
  controls, submission states, desktop/mobile layouts. Dashboard is minimal.
- Browser uses same-origin `/api/auth/*` proxy; Go checks exact Origin on mutations.
  HttpOnly/SameSite=Lax cookies; production requires HTTPS origins and Secure host cookie.
- Embedded Goose migrations via `cmd/migrate up|status`; Compose migration service
  completes before backend starts. Production migration is a deployment step.
- Reusable auth middleware resolves UUID user identity for future user-owned data.
- `/health` and `/ready` behavior preserved. PostgreSQL logs omit auth row detail/binds.
- Pure USDT-margined linear perpetual risk calculation supports LONG/SHORT, exact rational
  arithmetic, adverse fee/slippage fills, step-down quantity sizing and structured errors.
  It has no API, UI, persistence, exchange data, leverage or liquidation behavior.

## Verification (2026-09-16; historical)
PASS: frontend Jest (6 suites, 66 tests), lint, typecheck, production build;
Go unit tests, vet, gofmt; isolated PostgreSQL integration workflow; migration
apply/reapply/status; Compose configuration; real HTTP lifecycle/proxy/CSRF/session
smoke; browser registration/login/logout/redirects and desktop/mobile layouts.
Browser also confirmed invalid-submit focus after the accessibility fix.
Test-only accounts were deleted, isolated test containers/network removed.
Implementation committed in four parts: docs, backend, chore and frontend (through c24b8a2).
No push performed in this task. See AUTH-V1-HANDOFF.md for acceptance checklist and handoff.

## RISK-001 Verification (2026-09-22)
PASS: focused risk package tests, full backend Go tests and `go vet ./...`.
PASS: focused risk package formatting and exact-arithmetic source review (no binary floats).
The repository-wide `gofmt -l cmd internal db tests` check reports pre-existing files outside
the new risk package; none were changed by RISK-001. Runtime configuration was not modified.

## Local Runtime
Docker recovery and the local auth flow were confirmed by the owner on 2026-09-22.
This documentation task did not independently inspect runtime or rerun application tests.
The 2026-09-18 investigation remains in LOG.md as history, not an active blocker.
Last recorded local configuration: frontend/origin localhost:3001, backend 8080,
PostgreSQL 15432; current ports were not independently checked in this task.
Missing local `.env` was recreated from `.env.example`; example DB port remains 5432.
Docker executable is under the user's `AppData/Local/Programs/DockerDesktop/resources/bin`;
add that directory to the command process PATH when Docker is absent from PATH.
No production deployment performed; Secure cookie behavior covered by backend tests.

## Scope / Next
Reusable planning/review skill and coding-brief/handoff templates are maintained in
`docs/ai/skills/planner-reviewer/`; a copy is installed in personal Codex skills.
This does not install the skill into ChatGPT web. Application runtime was not rechecked
during this workflow-only change (2026-09-19).
No exchange connections, OAuth, verification/reset flows, roles, or trading features.
RISK-001 is awaiting review; do not start its API/UI follow-up until it is accepted.
Technical follow-ups (not local auth blockers):
- Map backend registration validation response 422 to a specific frontend message.
- Add login rate limiting per account/IP or at ingress before public deployment.
- Clean up expired session rows for database maintenance; authentication already rejects
  expired sessions. Rate limiting and cleanup require separate deployment/operations tasks.
