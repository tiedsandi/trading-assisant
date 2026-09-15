# Backend

Server HTTP awal menggunakan Go standard library. Modul: `trading-assistant/backend`.
Koneksi database menggunakan pgx/v5 (pgxpool). Instalasi dependency dilakukan
pengguna dengan perintah di bagian PostgreSQL. Router chi, sqlc, dan modul bisnis
belum ditambahkan.

## Struktur

- `cmd/api/main.go`: entry point, logging JSON, HTTP timeout, graceful shutdown.
- `internal/config/config.go`: membaca environment dan memvalidasi port HTTP.
- `internal/app/router.go`: merangkai `/health` dan `/ready`.
- `internal/platform/database`: konfigurasi pool, startup ping, dan test konfigurasi.
- `tests/integration`: test database nyata, diaktifkan dengan build tag `integration`.
- `internal/app/router_test.go`: pengujian endpoint dan metode/path yang tidak sesuai.
- `../docker/backend`: Dockerfile dan pengaturan Air.

Folder `internal/modules/<fitur>` dan `db/` baru dibuat
saat implementasinya dimulai. Logika bisnis nantinya berada di modul, bukan `main.go`.
SQL/migration/sqlc ditempatkan di `db/` saat integrasi database dibuat.

## Menjalankan (oleh pengguna)

Dari root repository dengan Docker Desktop aktif:

```powershell
docker compose up -d --build backend
docker compose logs -f backend
```

Compose ikut menjalankan PostgreSQL dan menunggu healthcheck database.
`Ctrl+C` hanya keluar dari tampilan log. Cek API:

```powershell
Invoke-RestMethod http://localhost:8080/health
```

Hasil: `status` bernilai `ok`. `/health` adalah liveness: proses HTTP berjalan.
Endpoint ini tidak memeriksa koneksi database. `/ready` menjalankan ping database:
HTTP 200 dengan `status: ready`, atau HTTP 503 dengan `status: not_ready`.
`DATABASE_URL` dari Compose sekarang digunakan; CORS belum diimplementasikan.

Konfigurasi yang digunakan: `APP_ENV` (default `development`), `HTTP_HOST`
(default `0.0.0.0`), dan `HTTP_PORT` (default `8080`). Port host ditentukan oleh
`BACKEND_PORT` di `.env` root. HTTP_PORT adalah port di dalam container.

Air membangun `./cmd/api` dan melakukan reload saat source Go berubah.
Untuk mencoba, ubah respons health sementara, simpan, cek log rebuild dan
respons endpoint, lalu kembalikan perubahan. Jika mengubah kontrak permanen,
perbarui test-nya juga.

## Pemeriksaan

Saat backend berjalan, jalankan dari root:

```powershell
docker compose exec backend gofmt -w cmd internal
docker compose exec backend go test ./...
docker compose exec backend go vet ./...
```

Test tidak membutuhkan database. Hasil pengujian belum dikonfirmasi sampai
perintah dijalankan. Tidak perlu menjalankan `go mod init` lagi.

## Testing bawaan Go

Backend menggunakan paket `testing` dan `net/http/httptest` bawaan Go, tanpa
Jest atau dependency test tambahan. Nama test berakhiran `_test.go` dan berada
di package yang sama dengan kode yang diuji.

Suite saat ini mencakup respons JSON `/health`, penolakan metode/path yang tidak
sesuai, konfigurasi default/override, port tidak valid, batas port, dan alamat IPv6.
Test environment memakai `t.Setenv` agar perubahan dipulihkan setelah test.

Dari root repository, bisa dijalankan tanpa server atau database aktif:

```powershell
docker compose run --rm --no-deps backend go test -count=1 -v ./...
docker compose run --rm --no-deps backend go test -count=1 -coverprofile=coverage.out ./...
docker compose run --rm --no-deps backend go tool cover -func=coverage.out
docker compose run --rm --no-deps backend go vet ./...
```

`-count=1` menghindari hasil test dari cache. `coverage.out` tidak ikut Git.
Air melakukan rebuild aplikasi, bukan menjalankan test otomatis; ulangi test
setelah perubahan kode. Entry point dan graceful shutdown belum dicakup suite ini.
Test handler belum menguji koneksi jaringan Docker maupun database PostgreSQL.

Saat koneksi database dibuat, tambahkan integration test terhadap database test
terpisah. Saat perhitungan trading dibuat, uji batas nilai, rounding, dan input
tidak valid dengan hasil yang ditentukan independen dari implementasi.

## PostgreSQL dan readiness

Dari root repository, pasang dependency sekali (oleh pengguna):

```powershell
docker compose run --rm --no-deps backend go get github.com/jackc/pgx/v5/pgxpool
docker compose run --rm --no-deps backend go mod tidy
docker compose run --rm --no-deps backend gofmt -w cmd internal tests
docker compose run --rm --no-deps backend go test -count=1 -v ./...
docker compose run --rm --no-deps backend go vet ./...
docker compose up -d backend
Invoke-RestMethod http://localhost:8080/ready
```

Commit `go.mod` dan `go.sum` hasil instalasi bersama kode. Sebelum dependency
dipasang, kode baru belum dapat dikompilasi. Tidak perlu mengulang `go mod init`.
Jika Air belum membangun ulang setelah instalasi, jalankan `docker compose restart backend`.

API mewajibkan `DATABASE_URL`; Compose sudah menyediakannya dengan hostname
`postgres` dan port internal `5432`. Untuk Go yang dijalankan langsung di host,
gunakan `localhost` dan port publish PostgreSQL di `.env` root (misalnya `15432`).
Aplikasi membaca environment proses, tidak memuat file `.env` backend otomatis.
Jangan mencetak URL database atau kredensial ke log.

Startup memvalidasi URL dan ping dalam batas 5 detik; kegagalan membuat API keluar.
Pool dibatasi 10 koneksi dan ditutup setelah server berhenti. Readiness memakai
context request dengan batas 2 detik. Jika database terputus setelah startup,
`/health` tetap 200 dan `/ready` menjadi 503; setelah koneksi pulih, readiness
dapat kembali 200. Belum ada migration atau tabel bisnis pada tahap ini.

## Integration test terisolasi

Gunakan file Compose test secara mandiri, bukan digabung dengan compose development:

```powershell
docker compose -p trading-assistant-tests -f compose.test.yml up --build --abort-on-container-exit --exit-code-from backend-test
docker compose -p trading-assistant-tests -f compose.test.yml down
```

Jalankan setelah `go get` dan `go mod tidy` di atas. PostgreSQL test tidak
memublikasikan port, memakai network terpisah dan penyimpanan sementara (tmpfs).
Perintah `down` hanya membersihkan container/network proyek test; data development
tidak disentuh. Cache Go test tetap tersimpan. Dependency test dapat diunduh saat run.

Suite integration memeriksa koneksi nyata, `SELECT 1`, readiness sukses, readiness
503 setelah pool ditutup, dan liveness yang tetap sukses. Tidak menulis tabel.
Ini belum mensimulasikan putus jaringan atau recovery database. Tanpa tag
`integration`, unit test tidak membutuhkan database. Dengan tag tersebut,
`TEST_DATABASE_URL` wajib ada; tidak ada skip diam-diam jika konfigurasi hilang.
Hasil test kode baru belum diverifikasi sampai pengguna menjalankan perintahnya.
