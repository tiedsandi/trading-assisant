# Trading Assistant

Panduan teknis development. PRD dan perencanaan produk dikelola di luar repository.

## Status dan struktur

Frontend sudah diinisialisasi dengan Next.js 16, React 19, TypeScript, Tailwind CSS 4,
dan Bun 1.4.2. Backend memiliki server HTTP minimal dan endpoint `/health`, dengan Go + Air.
PostgreSQL 18 dijalankan melalui Docker Compose. Setup ini khusus development.

```text
trading-assistant/
├── AGENTS.md
├── CLAUDE.md
├── README.md
├── .env.example
├── compose.yml
├── frontend/                 # Lihat frontend/README.md
├── backend/                  # Server HTTP; lihat backend/README.md
└── docker/
    ├── frontend/Dockerfile   # Bun
    ├── backend/
    │   ├── Dockerfile        # Go + Air
    │   └── air.toml          # Build ./cmd/api dan hot reload
    └── postgres/init/        # Script init hanya berjalan pada volume baru
```

Dockerfile tetap di `docker/`. Source aplikasi di-bind mount dari komputer.
Dependency, cache build, dan data database disimpan dalam named volume Docker.

## Menjalankan

Aktifkan Docker Desktop dengan Linux containers. Jalankan dari root repository
di PowerShell. Pada clone baru, buat `.env` jika belum tersedia:

```powershell
if (!(Test-Path .env)) { Copy-Item .env.example .env }
docker compose config --quiet
docker compose up -d --build postgres frontend
docker compose ps
docker compose logs -f frontend
```

Frontend memasang dependency sesuai `bun.lock` saat startup. `Ctrl+C` pada tampilan
log tidak menghentikan container. Jangan jalankan `create-next-app` lagi.

| Akses | Alamat |
| --- | --- |
| Frontend di browser | http://localhost:3000 |
| Database dari komputer | localhost:15432 |
| Database dari backend container | postgres:5432 |
| Backend (setelah container dijalankan) | http://localhost:8080/health |

`.env` root dibaca Compose, lalu variabel yang tercantum pada `environment`
diteruskan ke container. Tidak perlu menyalin `.env` ke frontend/backend.
`NEXT_PUBLIC_API_URL` untuk browser; `INTERNAL_API_URL` untuk server Next.js di Docker.
Kode integrasi API belum dibuat. Jika port aplikasi diubah, sesuaikan URL dan CORS.
Password database contoh hanya untuk lokal. Gunakan karakter URL-safe karena
password dirangkai ke `DATABASE_URL`. Perubahan password di `.env` tidak mengubah
password yang sudah tersimpan di database.

## Development

```powershell
docker compose exec frontend bun run lint
docker compose exec frontend bun run typecheck
docker compose run --rm --no-deps -e NODE_ENV=production frontend bun --bun run build
docker compose logs -f postgres
docker compose down
```

Tambahkan paket melalui `docker compose exec frontend bun add NAMA_PAKET`.
Commit `package.json` dan `bun.lock` bersama; gunakan Bun secara konsisten.
Dependency berada di volume Linux; editor Windows mungkin belum dapat membaca
tipe dari `node_modules` lokal.

Next.js menangani reload frontend, Air menangani backend. `WATCHPACK_POLLING`
hanya berlaku untuk Webpack, bukan Turbopack bawaan Next.js. Jika perubahan file
Windows tidak terdeteksi, coba Webpack setelah frontend berhasil diinstal:

```powershell
docker compose stop frontend
docker compose run --rm --service-ports --no-deps frontend bun --bun run dev --webpack --hostname 0.0.0.0
```

`docker compose down` mempertahankan data. `down -v` menghapus volume, termasuk
database. Volume tidak ikut Git; pindah laptop perlu backup/restore bila data
lokal ingin dibawa.

## Berikutnya

Backend sudah diinisialisasi dan kode server mendengarkan `0.0.0.0:8080` secara default.
Jalankan backend dengan `docker compose up -d --build backend`, lalu periksa
`http://localhost:8080/health`. Panduan test dan hot reload ada di [README backend](backend/README.md).
Pengguna sudah mengonfirmasi `/health`, test awal Go, Jest, lint, typecheck, dan
build frontend berhasil. `/health` hanya memeriksa proses HTTP, bukan database.
Kode koneksi pgx dan `/ready` kini tersedia; pasang dependency dan jalankan test
baru mengikuti [README backend](backend/README.md#postgresql-dan-readiness).
Integration test menggunakan `compose.test.yml` dengan database sementara terpisah.
Hasil test tahap database belum dikonfirmasi. Router chi ditambahkan ketika penataan API dimulai.
Untuk menjalankan semua service gunakan `docker compose up -d --build`.

Konfigurasi Jest frontend dan test bawaan Go sudah disiapkan. Instalasi dependency
Jest serta perintah test ada di [README frontend](frontend/README.md#testing-dengan-jest)
dan [README backend](backend/README.md#testing-bawaan-go). Jest dijalankan dengan Node.js
di image frontend; Bun tetap digunakan untuk instalasi paket dan Next.js.
Untuk clone baru, pasang dependency sesuai lockfile sebelum menjalankan pemeriksaan.
Belum ada autentikasi, fitur trading, Redis, atau worker.
