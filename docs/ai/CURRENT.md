# Current State

## Status
Authentication v1 implemented; last successful verification 2026-09-16.
Local runtime blocked on 2026-09-18; login/register recheck is incomplete.
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

## Verification (2026-09-16; historical)
PASS: frontend Jest (6 suites, 66 tests), lint, typecheck, production build;
Go unit tests, vet, gofmt; isolated PostgreSQL integration workflow; migration
apply/reapply/status; Compose configuration; real HTTP lifecycle/proxy/CSRF/session
smoke; browser registration/login/logout/redirects and desktop/mobile layouts.
Browser also confirmed invalid-submit focus after the accessibility fix.
Test-only accounts were deleted, isolated test containers/network removed.
No commit or push. See AUTH-V1-HANDOFF.md for full acceptance checklist.

## Local Runtime
As of 2026-09-18, Docker Engine is unavailable and the development stack is not running.
Port 3000 belongs to the separate asietex-erp-web project; leave it running.
Local ignored .env now uses FRONTEND_PORT=3001 and AUTH_ALLOWED_ORIGINS=http://localhost:3001;
backend/PostgreSQL ports remain 8080/15432. Port 3001 is not yet serving this app.
Docker Desktop crashes initializing sailor-ingest.sock in its local runtime directory.
Official stop/start did not recover it; automatic review blocked temporary socket removal.
No Docker files or volumes were deleted. Restore Docker before rerunning auth verification.
Missing local `.env` was recreated from `.env.example`; example DB port remains 5432.
Docker executable is under the user's `AppData/Local/Programs/DockerDesktop/resources/bin`;
add that directory to the command process PATH when Docker is absent from PATH.
No production deployment performed; Secure cookie behavior covered by backend tests.

## Scope / Next
No exchange connections, OAuth, verification/reset flows, roles, or trading features.
Choose the next product feature; scope owned records using the authenticated user UUID.
Next: restore Docker Engine, start Compose, then recheck registration/login/logout at localhost:3001.
The user-reported auth failure is not yet confirmed resolved.
