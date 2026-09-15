# Code Conventions

Based on current code/config and existing frontend/backend READMEs. Distinguish
documented placement rules from implemented features; evolve this file with the codebase.

## Frontend
- App Router routes/layout/global CSS live in `src/app`; static assets in `public`.
- Existing README assigns domain code to `src/features/<feature>` and reusable code
  to `components/{ui,common,layout}`, `hooks`, `lib`, `config`, `types`; these are placeholders today.
- Routes compose pages; feature components/hooks/types stay near their feature until shared.
- TypeScript strict; `@/*` maps to `src/*`; JSX uses `.tsx`, other TypeScript `.ts`.
- Documented naming: kebab-case files/folders, PascalCase components/types,
  camelCase functions/variables, `use`-prefixed hooks. Preserve framework filenames.
- Tailwind through PostCSS; ESLint Next core-web-vitals/TypeScript. No Prettier or formatter configured.
- Bun + `bun.lock`; keep manifest/lockfile changes together. Jest runs on Node.

## Backend and API
- `cmd/api`: wiring/lifecycle; `internal/config`: process environment;
  `internal/app`: HTTP handlers; `internal/platform/database`: pgxpool.
- Standard `net/http` ServeMux; no chi or sqlc dependency. Existing README reserves
  `internal/modules/<feature>` and `db/` for future implementation; neither exists yet.
- Health/readiness return JSON `status`; readiness has `Cache-Control: no-store`,
  a 2-second request-derived timeout, and generic 503 on database failure.
- Unknown routes/methods use ServeMux defaults; no general business API/error format established.
- Errors return to startup caller; JSON slog logs lifecycle. Database errors omit credentials.
- Pool max 10, startup ping/connect timeout 5 seconds; shutdown precedes pool close.
- Go formatting via gofmt; static checks via go vet. No additional Go linter configured.

## Tests
- Jest uses `next/jest`, jsdom, jest-dom, Testing Library, and role-based assertions.
- `tests/home.test.tsx` covers starter UI. README places feature/shared tests beside source,
  with page/cross-feature tests in `tests`; Jest discovers `*.test.ts(x)` in both.
- Go uses `testing`, `httptest`, table-driven cases, and `t.Setenv`; unit tests sit beside source.
- `tests/integration` requires tag `integration` and `TEST_DATABASE_URL`; missing URL fails.
  It checks SELECT 1, readiness, and closed-pool failure; it does not simulate network recovery.

## Database and environment
- No business SQL, migrations, or migration runner exists. Migration naming/versioning is not established.
- `docker/postgres/init` contains only README; init scripts run only for new volumes, not as a migration system.
- Development PostgreSQL uses a named volume and UTC; test PostgreSQL uses isolated tmpfs without host ports.
- Root `.env` configures Compose; `.env.example` is the tracked template. Go reads process env,
  not an env file. `DATABASE_URL` is required at startup; never log its credentials.
- `NEXT_PUBLIC_*` is browser-visible; keep secrets out. `INTERNAL_API_URL` points to backend in Docker.
- API URL and CORS variables are provisioned, but frontend API integration and backend CORS are not implemented.
