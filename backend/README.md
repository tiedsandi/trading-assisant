# Backend

Server HTTP awal menggunakan Go standard library. Modul: `trading-assistant/backend`.
Belum ada dependency aplikasi tambahan. Ini tahap bootstrap sebelum router chi,
pgx/sqlc, dan modul bisnis ditambahkan sesuai kebutuhan implementasi.

## Struktur

- `cmd/api/main.go`: entry point, logging JSON, HTTP timeout, graceful shutdown.
- `internal/config/config.go`: membaca environment dan memvalidasi port HTTP.
- `internal/app/router.go`: merangkai routes; sementara hanya endpoint `/health`.
- `internal/app/router_test.go`: pengujian endpoint dan metode/path yang tidak sesuai.
- `../docker/backend`: Dockerfile dan pengaturan Air.

Folder `internal/modules/<fitur>`, `internal/platform/database`, dan `db/` baru dibuat
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
Endpoint ini belum memeriksa koneksi database. `DATABASE_URL` dan CORS sudah
disediakan Compose, tetapi belum digunakan aplikasi. Integrasi PostgreSQL
dan readiness database adalah tahap berikutnya.

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
