# Backend

Go `net/http` API + pgx/v5 + PostgreSQL. Module: `trading-assistant/backend`.

## Struktur

- `cmd/api`: config, startup, HTTP timeouts, JSON lifecycle logs, shutdown.
- `cmd/migrate`: CLI migration `up` dan `status`.
- `internal/app`: komposisi route; `/health` dan `/ready` tetap tersedia.
- `internal/config`: environment, port, allowlist Origin.
- `internal/modules/auth`: handler, password/session primitives, PostgreSQL store, middleware.
- `internal/platform/database`: pool maksimal 10 koneksi, validasi URL dan startup ping.
- `db/migrations`: SQL Goose yang di-embed oleh `db/migrate.go`.
- `tests/integration`: lifecycle auth dan readiness dengan PostgreSQL nyata.

## Menjalankan dan migration

Dari root repository dengan Docker Desktop aktif:

```sh
docker compose up -d --build backend
docker compose logs -f backend
docker compose run --rm migrate go run -mod=readonly ./cmd/migrate status
docker compose run --rm migrate go run -mod=readonly ./cmd/migrate up
```

Compose menunggu PostgreSQL sehat, menjalankan service `migrate` sekali, lalu
menyalakan backend. API sendiri tidak mengubah schema pada startup. Di luar
Compose, set `DATABASE_URL` dan jalankan `go run ./cmd/migrate up` dari `backend/`
sebagai langkah deployment sebelum menjalankan API.

Goose v3 menyimpan versi di `goose_db_version`, menjalankan tiap migration SQL
secara transaksional, dan memakai PostgreSQL advisory lock agar runner paralel
tidak menerapkan migration yang sama bersamaan. `up` yang diulang aman.
Tambahkan migration baru dengan nomor naik, misalnya `00002_description.sql`,
dan anotasi `-- +goose Up` / `-- +goose Down`. Jangan mengubah migration yang sudah
diterapkan. CLI sengaja hanya menyediakan `up` dan `status`; rollback SQL tersedia
untuk review, tetapi perubahan produksi sebaiknya melalui migration koreksi.

`00001_auth.sql` membuat `users` (UUID, email unik/normalized, password hash,
created/updated timestamps) dan `sessions` (UUID, user FK cascade, digest unik
32-byte, waktu dibuat/kedaluwarsa). Index sesi tersedia untuk user dan expiry.

## HTTP API

| Endpoint | Perilaku |
| --- | --- |
| `GET /health` | Liveness 200, tanpa database. |
| `GET /ready` | Ping database dengan timeout 2 detik; 200 / 503. |
| `POST /auth/register` | JSON email/password; 201 user + session cookie langsung. |
| `POST /auth/login` | JSON email/password; 200 user + session cookie baru. |
| `GET /auth/me` | 200 user dari session aktif; 401 jika cookie hilang/invalid/expired/revoked. |
| `POST /auth/logout` | JSON `{}`; hapus session, clear cookie, 204. Idempotent. |

Sukses auth mengembalikan `{ "user": { "id", "email", "created_at" } }`.
Error menggunakan `{ "error": { "code", "message" } }`. Auth selalu `no-store`.
JSON body dibatasi 8 KiB; unknown fields, non-object, dan trailing data ditolak.
Invalid email/password registration: 422; duplicate: 409 `registration_unavailable`
dengan pesan umum untuk mencoba login/email lain. Status ini tetap memungkinkan
inferensi akun terdaftar; v1 memilih feedback yang usable tanpa endpoint lookup email.
Login unknown email dan password salah sama-sama 401 `invalid_credentials`.
Kegagalan penyimpanan mengembalikan 503 tanpa detail database.

Gunakan `authHandler.RequireUser(handler)` untuk route backend yang dilindungi.
Di dalam handler, panggil `auth.UserFromContext(r.Context())`; gunakan `User.ID`
sebagai FK/filter data milik user. Lookup session dibatasi 5 detik, sementara
handler selanjutnya tetap memakai context request asal.

## Password, session, cookie, CSRF

- Email di-trim dan lowercase, memakai bentuk email ASCII praktis; tidak menghapus
  `+tag` atau titik. Batas 254 byte total, local-part 64, domain label 63.
- Password 8–128 Unicode codepoints; tanpa aturan komposisi wajib, tanpa trim.
  Hash memakai `alexedwards/argon2id` / `x/crypto`, Argon2id 64 MiB, 3 iterasi,
  2 lanes, salt acak 16 byte, hasil 32 byte. Unknown email tetap memverifikasi
  dummy hash agar biaya hashing setara dengan password salah. Maksimal dua operasi
  password bersamaan per proses; kelebihan mendapat 429 + `Retry-After: 1`.
- Token session: `crypto/rand` 32 byte, base64url tanpa padding. Database hanya
  menyimpan SHA-256 dari token. Register menyimpan user + session dalam satu
  transaksi; login membuat session baru. Setiap request memeriksa expiry database.
- Masa berlaku tetap 7 hari, tanpa sliding renewal. Logout menghapus record sebelum
  menghapus cookie; jika revocation gagal, respons 503 dan cookie tetap ada.
- Cookie development: `ta_session`; production: `__Host-ta_session`. Keduanya
  `HttpOnly`, `SameSite=Lax`, `Path=/`, tanpa `Domain`, Max-Age 604800 dan Expires
  konsisten. Production wajib `Secure`; penghapusan memakai nama/path/flags sama.
- Browser mengakses route Next.js same-origin `/api/auth/*`; Next meneruskan cookie
  dan **Origin asli** ke Go dan meneruskan `Set-Cookie` ke browser. Tidak ada CORS
  allow headers. Token tidak digunakan di localStorage/sessionStorage.
- Semua POST, termasuk register/login/logout, wajib Origin persis dalam
  `AUTH_ALLOWED_ORIGINS` dan `Content-Type: application/json`. Missing/null/unknown
  Origin ditolak. Ini melindungi cookie auth dan login dari CSRF, ditambah SameSite.
  Server-to-server callers juga wajib mengirim Origin yang dikonfigurasi.
- Jangan log password, hash, token, cookie, atau database credentials.

Expired session tidak dipakai untuk autentikasi. Pembersihan berkala tabel bisa
menjalankan `DELETE FROM sessions WHERE expires_at <= CURRENT_TIMESTAMP` dalam
operasional database; scheduler cleanup bukan bagian v1.

## Environment

- `APP_ENV`: `development` (default), `test`, atau `production`; nilai lain ditolak.
- `HTTP_HOST`: `0.0.0.0`; `HTTP_PORT`: `8080` (port internal container).
- `DATABASE_URL`: wajib; dibaca dari environment proses, bukan file `.env` backend.
- `AUTH_ALLOWED_ORIGINS`: daftar Origin dipisahkan koma. Development default
  `http://localhost:3000`. Production wajib dikonfigurasi eksplisit dan seluruh
  Origin harus HTTPS tanpa path/query/fragment/credentials/wildcard.

Production memerlukan HTTPS pada origin frontend dan APP_ENV=production pada Go.
Next.js meneruskan kedua nama cookie; backend menentukan nama sesuai environment.
Air memantau source Go dan SQL embedded; reload API tidak menjalankan migration.

## Pemeriksaan

Dari root saat backend berjalan:

```sh
docker compose exec backend go test -mod=readonly -count=1 ./...
docker compose exec backend go vet ./...
docker compose exec backend gofmt -l cmd internal tests db
```

Format dengan mengganti `-l` menjadi `-w`. Unit tests mencakup input/JSON/Origin,
hashing, cookie, middleware, failure handling, session, login generik, dan logout.

Integration test harus menggunakan Compose test **secara terpisah**:

```sh
docker compose -p trading-assistant-tests -f compose.test.yml up --build --abort-on-container-exit --exit-code-from backend-test
docker compose -p trading-assistant-tests -f compose.test.yml down
```

Database test memakai tmpfs dan network sendiri, tanpa host port. Test menerapkan
migration yang sama, memeriksa idempotensi, uniqueness concurrent, rollback
registration saat session gagal, persistence hash, expiry/revocation, isolasi user,
protected backend, serta `/health` dan `/ready`. `TEST_DATABASE_URL` wajib dengan
tag `integration`; tidak ada skip diam-diam. Test hanya membersihkan user fixture
miliknya. Tanpa tag tersebut, suite tidak membutuhkan database.
