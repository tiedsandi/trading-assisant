# Handoff Report

Tanggal verifikasi akhir: 2026-09-16 (Asia/Jakarta).

## 1. Status

**DONE** — seluruh acceptance criteria yang berlaku terpenuhi.

## 2. What Was Implemented

- Register email/password/confirmation, langsung masuk dashboard setelah akun dibuat.
- Login, current user, sesi server 7 hari, logout dengan revocation dan cookie deletion.
- Dashboard minimal dengan email aman; guest/authenticated route guards di server.
- shadcn/ui + React Hook Form + Zod, live validation, show/hide, loading/error state,
  layout desktop/mobile, fokus field invalid dan semantics aksesibilitas.
- Backend net/http, PostgreSQL user/session store, migration berulang, dan tests.

## 3. Files Changed

- `.env.example` — Wiring origin, migration startup, dan konfigurasi development/logging.
- `AGENTS.md` — Dokumentasi workflow, arsitektur, migration, dan hasil verifikasi.
- `backend/cmd/api/main.go` — Wiring modul auth ke server HTTP.
- `backend/cmd/migrate/main.go` — CLI migration up/status dengan error aman.
- `backend/db/migrate.go` — Embedded Goose provider dengan PostgreSQL advisory lock.
- `backend/db/migrations/00001_auth.sql` — Schema users/sessions, constraints, foreign key, dan indeks.
- `backend/go.mod` — Dependency fitur auth beserta lock/checksum.
- `backend/go.sum` — Dependency fitur auth beserta lock/checksum.
- `backend/internal/app/router.go` — Registrasi endpoint auth; health/readiness tetap dipertahankan.
- `backend/internal/config/auth_test.go` — Regression test auth, validasi, security, atau integrasi database.
- `backend/internal/config/config.go` — Validasi environment dan allowlist Origin production/development.
- `backend/internal/modules/auth/handler_test.go` — Regression test auth, validasi, security, atau integrasi database.
- `backend/internal/modules/auth/handler.go` — HTTP lifecycle auth, validasi request, cookie, CSRF, dan middleware.
- `backend/internal/modules/auth/security.go` — Argon2id policy, random token/digest, normalisasi email, cookie.
- `backend/internal/modules/auth/store_integration_test.go` — Regression test auth, validasi, security, atau integrasi database.
- `backend/internal/modules/auth/store.go` — Query pgx, transaksi registrasi, lookup dan revocation sesi.
- `backend/README.md` — Dokumentasi workflow, arsitektur, migration, dan hasil verifikasi.
- `backend/tests/integration/auth_test.go` — Regression test auth, validasi, security, atau integrasi database.
- `compose.test.yml` — Wiring origin, migration startup, dan konfigurasi development/logging.
- `compose.yml` — Wiring origin, migration startup, dan konfigurasi development/logging.
- `docker/backend/air.toml` — Wiring origin, migration startup, dan konfigurasi development/logging.
- `docs/ai/CONVENTIONS.md` — Dokumentasi workflow, arsitektur, migration, dan hasil verifikasi.
- `docs/ai/CURRENT.md` — Dokumentasi workflow, arsitektur, migration, dan hasil verifikasi.
- `docs/ai/LOG.md` — Dokumentasi workflow, arsitektur, migration, dan hasil verifikasi.
- `frontend/bun.lock` — Dependency fitur auth beserta lock/checksum.
- `frontend/components.json` — Konfigurasi integrasi shadcn/ui Tailwind 4.
- `frontend/package.json` — Dependency fitur auth beserta lock/checksum.
- `frontend/README.md` — Dokumentasi workflow, arsitektur, migration, dan hasil verifikasi.
- `frontend/src/app/api/auth/[action]/route.ts` — Route/UI auth, guard server, loading/error state, metadata dan CSS.
- `frontend/src/app/dashboard/page.tsx` — Route/UI auth, guard server, loading/error state, metadata dan CSS.
- `frontend/src/app/error.tsx` — Route/UI auth, guard server, loading/error state, metadata dan CSS.
- `frontend/src/app/globals.css` — Route/UI auth, guard server, loading/error state, metadata dan CSS.
- `frontend/src/app/layout.tsx` — Route/UI auth, guard server, loading/error state, metadata dan CSS.
- `frontend/src/app/loading.tsx` — Route/UI auth, guard server, loading/error state, metadata dan CSS.
- `frontend/src/app/login/page.tsx` — Route/UI auth, guard server, loading/error state, metadata dan CSS.
- `frontend/src/app/page.tsx` — Route/UI auth, guard server, loading/error state, metadata dan CSS.
- `frontend/src/app/register/page.tsx` — Route/UI auth, guard server, loading/error state, metadata dan CSS.
- `frontend/src/components/ui/alert.tsx` — Primitif shadcn/ui yang dipakai auth.
- `frontend/src/components/ui/button.tsx` — Primitif shadcn/ui yang dipakai auth.
- `frontend/src/components/ui/card.tsx` — Primitif shadcn/ui yang dipakai auth.
- `frontend/src/components/ui/input.tsx` — Primitif shadcn/ui yang dipakai auth.
- `frontend/src/components/ui/label.tsx` — Primitif shadcn/ui yang dipakai auth.
- `frontend/src/features/auth/auth-shell.tsx` — Form, validasi, transport same-origin, state/session dan guard auth.
- `frontend/src/features/auth/client.test.ts` — Regression test auth, validasi, security, atau integrasi database.
- `frontend/src/features/auth/client.ts` — Form, validasi, transport same-origin, state/session dan guard auth.
- `frontend/src/features/auth/forms.test.tsx` — Regression test auth, validasi, security, atau integrasi database.
- `frontend/src/features/auth/login-form.tsx` — Form, validasi, transport same-origin, state/session dan guard auth.
- `frontend/src/features/auth/logout-button.tsx` — Form, validasi, transport same-origin, state/session dan guard auth.
- `frontend/src/features/auth/password-input.tsx` — Form, validasi, transport same-origin, state/session dan guard auth.
- `frontend/src/features/auth/proxy.test.ts` — Regression test auth, validasi, security, atau integrasi database.
- `frontend/src/features/auth/proxy.ts` — Form, validasi, transport same-origin, state/session dan guard auth.
- `frontend/src/features/auth/register-form.tsx` — Form, validasi, transport same-origin, state/session dan guard auth.
- `frontend/src/features/auth/schemas.test.ts` — Regression test auth, validasi, security, atau integrasi database.
- `frontend/src/features/auth/schemas.ts` — Form, validasi, transport same-origin, state/session dan guard auth.
- `frontend/src/features/auth/server.test.ts` — Regression test auth, validasi, security, atau integrasi database.
- `frontend/src/features/auth/server.ts` — Form, validasi, transport same-origin, state/session dan guard auth.
- `frontend/src/features/auth/session-cookie.ts` — Form, validasi, transport same-origin, state/session dan guard auth.
- `frontend/src/lib/utils.ts` — Helper cn standar shadcn untuk penggabungan class.
- `frontend/tests/home.test.tsx` — Regression test auth, validasi, security, atau integrasi database.
- `README.md` — Dokumentasi workflow, arsitektur, migration, dan hasil verifikasi.
- `docs/ai/AUTH-V1-HANDOFF.md` — Laporan lengkap dan checklist penerimaan.

## 4. Database Changes

Migration `backend/db/migrations/00001_auth.sql`, dijalankan oleh Goose v3.28.0.
SQL di-embed di binary; provider memakai PostgreSQL advisory lock untuk serialisasi runner.
`cmd/migrate` menyediakan `up` dan `status`. Compose menjalankan service `migrate`
setelah PostgreSQL sehat dan sebelum backend startup. PostgreSQL init scripts tidak dipakai
sebagai migration. Pengulangan migration terverifikasi menerapkan 0 migration tambahan.

- `users`: UUID PK, email normalized unique, password_hash, created_at/updated_at timestamptz.
  CHECK memastikan normalisasi/panjang email; constraint unik menangani race registrasi.
- `sessions`: UUID PK, user_id FK (ON DELETE CASCADE), token_hash bytea unique (32 byte),
  expires_at dan created_at timestamptz; indeks user_id dan expires_at.
- User dan sesi registrasi dibuat dalam satu transaksi. Logout menghapus sesi aktif.
- `goose_db_version` melacak migration yang diterapkan. Tidak ada schema trading/exchange.

Perintah dari root:

```sh
docker compose run --rm migrate
docker compose run --rm migrate go run -mod=readonly ./cmd/migrate status
```

## 5. API Changes

| Endpoint Go | Perilaku |
| --- | --- |
| `POST /auth/register` | `{email,password}`; 201 + safe user + cookie sesi; invalid 422; duplicate 409 |
| `POST /auth/login` | `{email,password}`; 200 + safe user + cookie baru; kredensial salah 401 generik |
| `GET /auth/me` | 200 + safe user; missing/invalid/expired/revoked 401 |
| `POST /auth/logout` | JSON `{}`; delete sesi aktif, clear cookie, 204; idempotent |

Browser memakai padanannya di `/api/auth/*` melalui proxy Next.js.
Safe response: `{user:{id,email,created_at}}`. Error: `{error:{code,message}}`.
Semua mutasi wajib Origin allowlist dan `application/json`; unknown fields, JSON malformed,
non-object dan multiple values ditolak. Body Go maksimal 8 KiB; proxy 16 KiB. Auth no-store.
Overload pekerjaan hashing memberi 429; backend unavailable 503, kegagalan proxy 502 generik.
`RequireUser` dan `UserFromContext` tersedia untuk endpoint protected berikutnya.
`/health` dan `/ready` tetap memiliki perilaku sebelumnya.

## 6. Verification

- `docker compose config --quiet`: PASS.
- `docker compose up -d --build backend` setelah PATH/port lokal diperbaiki: PASS.
- Migration up, status, reapply tanpa perubahan: PASS.
- `docker compose exec frontend bun run test`: PASS — 6 suites, 66 tests.
- `docker compose exec frontend bun run lint`: PASS.
- `docker compose exec frontend bun run typecheck`: PASS.
- `docker compose run --rm --no-deps -e NODE_ENV=production frontend bun --bun run build`: PASS.
  Dev frontend dihentikan sementara untuk mencegah konflik `.next`, lalu dijalankan kembali.
- `docker compose exec backend go test -mod=readonly -count=1 ./...`: PASS.
- `docker compose exec backend go vet ./...`: PASS.
- `docker compose exec backend gofmt -l cmd internal db tests`: PASS — output kosong.
- Isolated PostgreSQL Compose workflow dari AGENTS.md: PASS — exit 0; container/network dibersihkan.
  Termasuk registration rollback, duplicate race, hashing/digest persistence, expiry,
  revocation, middleware, isolasi dua pengguna dan health/readiness.
- HTTP smoke lewat Next.js→Go: PASS — origin rejection, register/login/me/logout,
  cookie attributes/expiry, SSR navigation, distinct users, revoked-cookie replay ditolak.
- Browser smoke localhost: PASS — live validation, mismatch, visibility, register otomatis
  masuk dashboard, error login generik, login/logout, protected/guest redirects dan fokus invalid.
- Visual desktop 1280×720 dan mobile 390×844: PASS — register/dashboard serta login diperiksa.
- Live `/health` dan `/ready`: PASS — keduanya 200 setelah startup terakhir.
- PostgreSQL constraint-test logging: PASS — error tetap tercatat, row detail/bind values dihilangkan.
- `git diff --check`: PASS.

Masalah awal yang sudah diperbaiki: executable/credential-helper Docker belum ada di PATH,
`.env` hilang dan port 5432 bentrok, tipe body proxy, typings locator test, dan fokus field
saat validasi async. HTTP harness disesuaikan dengan redirect HTML streaming Next.js;
redirect juga diverifikasi langsung di browser. Tidak ada kegagalan pemeriksaan tersisa.

Verifikasi browser memakai HTTP lokal; production Secure-cookie flags/config diuji otomatis.
Deployment HTTPS publik tidak dilakukan karena tidak termasuk scope.

## 7. Acceptance Criteria

Checklist berikut menyalin semua acceptance criteria dari instruksi. Seluruhnya terverifikasi;
indikator kualitas password tambahan bersifat opsional dan tidak ditambahkan.
### Registration

- [x] User can register with email, password, and frontend password confirmation.
- [x] Email format is validated in frontend and backend.
- [x] Email uniqueness is enforced reliably by the database.
- [x] Password minimum 8 characters is enforced by the backend.
- [x] Password is never persisted as plaintext.
- [x] Successful registration immediately creates an authenticated 7-day session.
- [x] Successful registration redirects the user to the authenticated dashboard.

### Frontend UX

- [x] Auth UI uses shadcn/ui for the new form/interface components.
- [x] Forms use React Hook Form + Zod.
- [x] Register validation provides live/user-friendly feedback.
- [x] Password UI visibly reports the 8-character requirement.
- [x] Optional uppercase/lowercase/number/special-character indicators are hints unless explicitly enforced.
- [x] Confirm-password mismatch is validated client-side.
- [x] Login has appropriate email/non-empty-password validation.
- [x] Password visibility can be toggled where appropriate.
- [x] Submission/loading state prevents accidental repeated submission.
- [x] Backend errors are presented cleanly.
- [x] Auth pages are usable on mobile and desktop.
- [x] Forms use accessible labels and validation/error semantics.

### Login / Session

- [x] Valid credentials can login.
- [x] Invalid credentials return a generic error.
- [x] Login creates a cryptographically secure server-side session.
- [x] Session lifetime is 7 days.
- [x] Raw session tokens are not stored in PostgreSQL.
- [x] Browser authentication uses an HttpOnly cookie.
- [x] Production cookie behavior uses Secure appropriately.
- [x] Authentication token is not stored in localStorage/sessionStorage.
- [x] Current authenticated user can be retrieved from the backend.
- [x] Missing, invalid, expired, or revoked sessions are rejected.
- [x] User A's session cannot authenticate as User B.

### Navigation / Protection

- [x] Guest users cannot use the authenticated area.
- [x] Guest access to protected UI routes leads to login.
- [x] Authenticated users visiting login/register are directed to the authenticated area.
- [x] Authenticated dashboard/shell shows safe current-user information.
- [x] Dashboard remains minimal and contains no trading feature implementation.

### Logout

- [x] Logout revokes/deletes the active server-side session.
- [x] Logout clears the authentication cookie.
- [x] Revoked session cannot be reused.
- [x] User is returned to login after logout.

### Database / Architecture

- [x] Users and sessions are created through reproducible migrations.
- [x] Migration workflow is documented.
- [x] Appropriate constraints/indexes exist for email uniqueness and session lookup.
- [x] Auth handling is reusable for future protected backend endpoints.
- [x] Future user-owned trading data can reliably reference the authenticated user.

### Security / Integration

- [x] Backend independently validates all security-relevant input.
- [x] Password hashing uses an established implementation.
- [x] Secrets/passwords/raw session tokens are not exposed in logs or responses.
- [x] Cookie/CORS/credentials behavior works with the actual frontend/backend architecture.
- [x] CSRF implications are evaluated and the chosen protection is implemented/documented where required.
- [x] Existing health/readiness behavior remains functional.

### Quality

- [x] Relevant backend unit tests exist and pass.
- [x] Relevant backend integration tests exist and pass.
- [x] Relevant frontend auth tests exist and pass.
- [x] Frontend lint passes.
- [x] Frontend typecheck passes.
- [x] Frontend production build passes.
- [x] Go tests pass.
- [x] Go vet passes.
- [x] Go formatting check passes.
- [x] No unrelated feature/refactor is included.

## 8. Decisions / Assumptions

- Password: `github.com/alexedwards/argon2id` di atas Go x/crypto; Argon2id 64 MiB,
  3 iterations, 2 lanes, random salt 16 byte, key 32 byte. Min 8/max 128 Unicode codepoints;
  tanpa composition rules atau pemotongan password. Unknown email tetap menjalankan
  pembandingan dummy hash; maksimal 2 pekerjaan hashing paralel untuk membatasi memori.
  Dasar primitive: [Go Argon2](https://pkg.go.dev/golang.org/x/crypto/argon2).
- Email: trim + lowercase seluruh alamat dan validasi ASCII praktis yang konsisten
  frontend/backend; tidak menghapus titik atau plus-addressing khusus provider.
- Migration: Goose embedded SQL, version table dan advisory lock; explicit deployment
  step, tanpa DDL otomatis dalam proses API. [Goose Provider](https://pressly.github.io/goose/documentation/provider/).
- Identity: UUID pengguna menjadi referensi data milik pengguna; email bukan foreign key.
  Auth aplikasi terpisah dari kredensial exchange; tidak ada abstraksi OAuth spekulatif.
- Sessions: crypto/rand 32 byte → base64url token; hanya SHA-256 token disimpan.
  Fixed 7 hari tanpa sliding/refresh token, UUID session, lookup database setiap request,
  delete aktif saat logout. Sesi perangkat lain tidak ikut dicabut oleh logout satu sesi.
- Cookie: `ta_session` untuk development HTTP; `__Host-ta_session` untuk production.
  Path=/, tanpa Domain, HttpOnly, SameSite=Lax, lifetime 604800 detik; Secure di production.
  Logout memakai nama/path/flags sama dengan Max-Age=0 dan expiration di masa lalu.
- CORS/CSRF: browser same-origin ke Next proxy, jadi CORS credentialed lintas-origin
  tidak diperlukan. Proxy meneruskan Origin asli dan hanya cookie auth. Go menolak
  Origin absent/null/tidak allowlisted pada seluruh mutasi, termasuk login/register/logout,
  serta mewajibkan JSON. Production hanya menerima konfigurasi origin HTTPS eksplisit.
  [OWASP CSRF prevention](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html).
- Duplicate registration: 409 dengan pesan umum yang mengarahkan ke sign-in atau email lain.
  Respons tidak memuat detail akun; status konflik masih bisa mengindikasikan alamat sudah
  terdaftar. Ini tradeoff UX v1 yang sengaja didokumentasikan. Login selalu 401/pesan sama.
- Logging: Go tidak mencatat body/cookie/hash/DB detail. Setelah test constraint menunjukkan
  PostgreSQL mencetak row detail, Compose dev/test menggunakan terse error verbosity dan
  menonaktifkan logging bind parameters. Error utama tetap tersedia.
  [PostgreSQL logging](https://www.postgresql.org/docs/current/runtime-config-logging.html).
- Deployment publik memerlukan HTTPS dan konfigurasi ingress sesuai kapasitas, termasuk
  pembatasan percobaan login. Admission cap hashing bukan rate limit per pengguna/IP.
- Lingkungan lokal: `.env` dibuat dari template karena hilang, hanya port DB disetel 15432
  agar tidak bentrok dengan PostgreSQL host. Tiga akun QA task dihapus setelah smoke test.

## 9. Issues / Remaining Work

None dalam scope Authentication v1. Tidak ada trading/exchange, reset password, email
verification, OAuth, 2FA atau roles yang ditambahkan. Stack development ditinggalkan berjalan.

## 10. Git

- Commit: NOT CREATED
- Push: NOT PERFORMED

## 11. Suggested Review

- `backend/internal/modules/auth/{handler,security,store}.go`: hash, transaksi, sesi,
  generic errors, Origin checks, cookie clearing, context middleware.
- `backend/db/migrations/00001_auth.sql`, `backend/db/migrate.go`, `backend/cmd/migrate/main.go`:
  schema/constraints/indexes, embedded migrations, advisory locking dan deployment workflow.
- `backend/internal/config/config.go`, `compose.yml`, `compose.test.yml`, `.env.example`:
  fail-closed production configuration dan privacy logging.
- `frontend/src/features/auth/{proxy,server,session-cookie,client}.ts` dan route guards:
  cookie forwarding, original Origin, no-store, backend outage behavior dan revocation.
- `frontend/src/features/auth/*form.tsx` dan tests: validasi live, confirmation,
  focus/read-only behavior, loading state dan accessible feedback.
