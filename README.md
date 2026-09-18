# Trading Assistant

Next.js + Bun frontend, Go `net/http` API, PostgreSQL 18. Authentication v1
menyediakan register, login, sesi 7 hari, dashboard minimal, dan logout.
Fitur trading belum diimplementasikan. Compose ini khusus development lokal.

## Menjalankan

Aktifkan Docker Desktop dengan Linux containers, lalu dari root repository:

```powershell
if (!(Test-Path .env)) { Copy-Item .env.example .env }
docker compose config --quiet
docker compose up -d --build
docker compose ps
```

Frontend memasang dependency dari `bun.lock`. Service `migrate` menjalankan
migration Goose sebelum backend dimulai; perubahan schema tidak bergantung pada
inisialisasi volume PostgreSQL. Air menjalankan reload backend.

| Akses | Default |
| --- | --- |
| Aplikasi / register / login | http://localhost:3000 |
| Liveness backend | http://localhost:8080/health |
| Readiness database | http://localhost:8080/ready |
| PostgreSQL dari host | localhost:5432; `.env` lokal dapat memakai 15432 |

`.env` root hanya dibaca Compose. `AUTH_ALLOWED_ORIGINS` berisi origin browser
persis (mis. `http://localhost:3000`), dipisah koma jika lebih dari satu.
Jika port frontend diubah, ubah origin tersebut juga. `INTERNAL_API_URL` adalah
alamat Go dari server Next.js (`http://backend:8080` di Docker). Browser memakai
`/api/auth/*` pada origin aplikasi. `NEXT_PUBLIC_API_URL` dan
`CORS_ALLOWED_ORIGINS` lama tidak lagi digunakan.

## Migration dan pemeriksaan

```powershell
docker compose run --rm migrate
docker compose run --rm migrate go run -mod=readonly ./cmd/migrate status
docker compose exec frontend bun run test
docker compose exec frontend bun run lint
docker compose exec frontend bun run typecheck
docker compose run --rm --no-deps -e NODE_ENV=production frontend bun --bun run build
docker compose exec backend go test -mod=readonly -count=1 ./...
docker compose exec backend go vet ./...
docker compose exec backend gofmt -l cmd internal db tests
docker compose -p trading-assistant-tests -f compose.test.yml up --build --abort-on-container-exit --exit-code-from backend-test
docker compose -p trading-assistant-tests -f compose.test.yml down
```

Service `migrate` dapat dijalankan ulang: versi yang sudah diterapkan tidak
berjalan lagi. Setelah menambah migration saat container aktif, jalankan kembali
migration secara eksplisit. Integration test memakai database tmpfs dan network
terpisah, tanpa port host. Jangan menggabungkan kedua file Compose.

Detail ada di [backend](backend/README.md), [frontend](frontend/README.md), dan
[konvensi](docs/ai/CONVENTIONS.md). Perintah rutin agen ada di [AGENTS.md](AGENTS.md).

## Deployment authentication

Gunakan HTTPS, `APP_ENV=production`, dan `AUTH_ALLOWED_ORIGINS` HTTPS yang eksplisit.
Jalankan migration sebagai langkah deployment sebelum API menerima traffic.
Cookie production bernama `__Host-ta_session`: Secure, HttpOnly, SameSite=Lax,
Path=/, tanpa Domain. Development HTTP memakai `ta_session`. Next.js meneruskan
Origin asli ke Go; jangan menggantinya dengan origin tepercaya pada reverse proxy.

Go adalah otoritas akses: endpoint protected harus memakai middleware auth dan
mengambil ID pengguna dari context. Setiap query data milik pengguna di masa
mendatang harus dibatasi oleh ID itu. Kredensial exchange terpisah dari login app.

PostgreSQL Compose memakai error verbosity `terse` dan menonaktifkan logging bind parameters
agar kegagalan constraint tidak mencetak baris auth. Terapkan pengaturan setara di deployment.

Password database contoh hanya untuk lokal. Simpan secret di environment yang
terkelola dan jangan log request body/cookie. Untuk deployment publik, tetapkan
pembatasan percobaan login pada ingress sesuai kapasitas dan topologi deployment.
`docker compose down` mempertahankan data; `down -v` menghapus volume database.
