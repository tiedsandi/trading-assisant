# Agent Workflow

Trading Assistant is a development scaffold: Next.js frontend, Go HTTP API,
and PostgreSQL in Docker. Product scope is not yet implemented.

## Start here
- Before coding, read `AGENTS.md`, `docs/ai/CURRENT.md`, and `docs/ai/CONVENTIONS.md`.
- Inspect `git status` and relevant source/config; the repository is the source of truth.
- Read `docs/ai/LOG.md` only when historical decisions are needed.

## Repository map
- `frontend/`: App Router in `src/app`, shared/feature placeholders, `tests`, Bun manifests.
- `backend/`: `cmd/api`, `internal/app`, `internal/config`, `internal/platform/database`, `tests/integration`.
- `docker/`: frontend/backend images, Air config, PostgreSQL init directory.
- `compose.yml`: development services; `compose.test.yml`: isolated database tests.
- `.env.example`: root Compose configuration template; `docs/ai/`: shared context.

## Commands
Run from repository root with Docker Desktop/Linux containers. On a fresh clone,
copy `.env.example` to `.env` only if missing. Frontend startup installs locked dependencies.

```sh
docker compose config --quiet
docker compose up -d --build
# Or start individually (backend also starts PostgreSQL):
docker compose up -d --build frontend
docker compose up -d --build backend
docker compose logs -f backend
```

Checks below use running containers with dependencies available:

```sh
docker compose exec frontend bun run test
docker compose exec frontend bun run lint
docker compose exec frontend bun run typecheck
docker compose run --rm --no-deps -e NODE_ENV=production frontend bun --bun run build
docker compose exec backend go test -mod=readonly -count=1 ./...
docker compose exec backend go vet ./...
docker compose exec backend gofmt -l cmd internal tests
```

Frontend also has `test:watch`, `test:coverage`, and `test:ci` scripts. Use `bun run test`,
not the Bun test runner. To apply Go formatting, replace `gofmt -l` with `gofmt -w`.
No frontend formatter script is configured.

Integration tests use a separate Compose project (never merge with development Compose):

```sh
docker compose -p trading-assistant-tests -f compose.test.yml up --build --abort-on-container-exit --exit-code-from backend-test
docker compose -p trading-assistant-tests -f compose.test.yml down
```

## Working rules
- Keep the requested scope; explain any necessary scope change. Avoid unrelated refactors.
- Preserve user changes. Never remove tests or weaken checks just to make them pass.
- Follow existing conventions. Database schema changes must use migrations; tooling is not yet established.
- Run relevant tests after changes; report actual results or explicitly state what was not run.
- Do not commit secrets, dependencies, or generated output. Never commit/push unless the user asks.
- At completion, update `docs/ai/CURRENT.md` to the latest snapshot. Append a dated,
  concise `docs/ai/LOG.md` entry for meaningful changes/decisions and verification.
- Keep context short: workflow here, snapshot in CURRENT, history in LOG, code patterns in CONVENTIONS.
