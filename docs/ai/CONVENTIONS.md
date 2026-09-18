# Code Conventions

## Frontend
- Next.js App Router routes/layout/global CSS live in `src/app`; routes compose
  feature code in `src/features/<feature>`. Shared primitives live in `src/components/ui`.
- TypeScript strict; `@/*` maps to `src/*`; kebab-case files, PascalCase components,
  camelCase functions, `use`-prefixed hooks. Preserve framework filenames.
- Tailwind 4 with shadcn/ui `components.json`, shared CSS tokens, and `cn` helper.
  Auth forms use React Hook Form + Zod; backend remains the security authority.
- Bun + `bun.lock`; commit manifest/lock together. Jest runs on Node through
  `bun run test`, not `bun test`. ESLint Next core-web-vitals/TypeScript; no formatter.
- Browser auth uses same-origin `/api/auth/*`. A narrow Next route proxy forwards
  original Origin, session cookies and backend Set-Cookie. No token in browser storage.
- Server-rendered pages resolve `/auth/me` through `INTERNAL_API_URL`, with no-store
  and a timeout. Only unauthorized sessions become guests; outages surface as errors.
- Resolve auth at each protected page/data access, not solely in a shared layout or
  client navigation. Future protected backend endpoints must independently check auth.

## Backend and API
- `cmd/api` wires lifecycle; `internal/config` validates process environment;
  `internal/app` owns net/http ServeMux; `internal/modules/auth` owns authentication;
  `internal/platform/database` owns pgxpool. Keep module logic out of main.
- Auth routes are `/auth/register`, `/auth/login`, `/auth/me`, `/auth/logout`.
  JSON success uses `{ "user": ... }`; errors use `{ "error": { "code", "message" } }`.
  Logout succeeds with 204. Return safe user fields only, never password/session hashes.
- Reuse auth middleware and obtain the authenticated user from request context.
  User-owned records should reference the user's UUID and queries must use this ID.
- Mutation endpoints require an exact allowlisted Origin and JSON content type;
  reject absent/untrusted Origin. Do not infer trust from client-supplied forwarding headers.
  Browser calls are same-origin; cross-origin CORS access is not enabled.
- Passwords: Argon2id via maintained library, 8–128 Unicode codepoints for registration;
  no composition requirement. Email trim/lowercase, practical ASCII format, unique in DB.
- Sessions: 32 random bytes encoded base64url; store SHA-256 only, fixed 7-day expiry;
  logout deletes the active session. Cookie Path=/, HttpOnly, SameSite=Lax, no Domain.
  `APP_ENV=production` requires HTTPS origins and Secure `__Host-ta_session`;
  development uses `ta_session` on HTTP localhost.
- Keep `/health` process-only and `/ready` a no-store 2-second database ping
  (200 ready / 503 not_ready). Preserve request cancellation.
- JSON slog lifecycle logs omit secrets. Never log request bodies, session tokens,
  hashes, database URLs, or raw database errors that can include credentials/user data.
- PostgreSQL uses `log_error_verbosity=terse` and zero bind-parameter logging so constraint
  errors omit auth row values. Keep equivalent settings in production database logging.
- Go formatting via gofmt; static checks via go vet. Pool max 10, 5-second startup ping,
  graceful HTTP shutdown before pool close.

## Database and migrations
- Versioned Goose SQL lives in `backend/db/migrations` and is embedded by `backend/db`.
  Use sequential filenames and Goose Up/Down annotations; applied migrations are immutable.
- `cmd/migrate up` applies pending versions; `cmd/migrate status` reports state.
  Goose records schema versions and uses a PostgreSQL advisory lock for concurrent runners.
  Do not replace migrations with PostgreSQL init scripts or automatic API-startup DDL.
- Compose runs a one-shot `migrate` service after PostgreSQL is healthy and before API
  startup. After adding migrations to a running stack, explicitly rerun that service.
  Production deployment must run migrations before starting/replacing API instances.
- Users use UUID primary keys, unique normalized email and timestamptz timestamps.
  Sessions reference users by foreign key; unique token digest and lookup/expiry indexes.
- Development PostgreSQL uses a named volume and UTC. Integration database is tmpfs,
  on its own Compose project/network without host ports. Never target development data in tests.
- Root `.env` configures Compose; `.env.example` is the template. Go reads process env.
  `DATABASE_URL` is required; `NEXT_PUBLIC_*` values are public and cannot hold secrets.

## Tests
- Jest uses next/jest, jsdom, Testing Library and role-based assertions. Feature tests
  sit beside source; cross-feature/page tests in `frontend/tests`. Node environment
  may be selected for server-only code and route-handler tests.
- Go uses testing, httptest, table-driven cases and t.Setenv. Unit tests sit beside code.
- `backend/tests/integration` uses build tag `integration` and mandatory
  `TEST_DATABASE_URL`. Apply the same embedded migrations against the isolated real DB.
  Cover auth lifecycle, expiry/revocation, constraints and cross-user isolation alongside readiness.
- Run checks listed in AGENTS.md and report executed results; do not claim unrun checks pass.
