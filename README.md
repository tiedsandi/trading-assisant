# Trading Assistant

Panduan teknis development. PRD dan perencanaan produk dikelola di luar repository.

## Status dan struktur

Frontend sudah diinisialisasi dengan Next.js 16, React 19, TypeScript, Tailwind CSS 4,
dan Bun 1.4.2. Backend masih kosong; konfigurasi Go + Air sudah tersedia.
PostgreSQL 18 dijalankan melalui Docker Compose. Setup ini khusus development.

```text
trading-assistant/
├── AGENTS.md
├── CLAUDE.md
├── README.md
├── .env.example
├── compose.yml
├── frontend/                 # Lihat frontend/README.md
├── backend/                  # Belum diinisialisasi
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
| Backend setelah diimplementasikan | http://localhost:8080 |

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
docker compose exec frontend bun --bun run build
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

Backend belum dapat dijalankan. Task berikutnya: inisialisasi modul Go, buat
`backend/cmd/api/main.go`, baca environment, dan sediakan endpoint `/health`.
Server harus mendengarkan `0.0.0.0:8080`. Setelah backend tersedia, jalankan semua
service dengan `docker compose up -d --build`.

Belum ada test runner frontend, autentikasi, fitur trading, Redis, atau worker.
