# Frontend

Next.js App Router, React, TypeScript, Tailwind CSS 4, dan Bun. Authentication v1
menggunakan shadcn/ui, React Hook Form, dan Zod. Cara menjalankan stack dan migrasi
ada di [README root](../README.md).

## Authentication v1

| Route | Perilaku |
| --- | --- |
| `/` | Mengarahkan guest ke `/login`, user terautentikasi ke `/dashboard` |
| `/register` | Email, password, konfirmasi; registrasi sukses langsung masuk dashboard |
| `/login` | Email/password; kegagalan kredensial ditampilkan dengan pesan generik |
| `/dashboard` | Memeriksa session di server, menampilkan email user dan aksi sign out |
| `/api/auth/register`, `/api/auth/login`, `/api/auth/logout` | Proxy POST khusus autentikasi ke backend Go |
| `/api/auth/me` | Proxy GET current user ke backend Go |

User terautentikasi yang membuka login/register diarahkan ke dashboard. Setiap
halaman memeriksa autentikasi di server; nama folder/layout atau state client
bukan pengaman akses. Backend tetap memvalidasi session pada setiap endpoint
terproteksi. Dashboard hanya shell akun, belum berisi fitur trading.

### Form dan validasi

- Email di-trim dan di-lowercase, lalu diperiksa sebagai alamat ASCII konvensional.
- Registrasi menerima password 8–128 Unicode code points; tidak ada syarat
  uppercase, angka, atau simbol. Password tidak di-trim.
- Konfirmasi password hanya untuk frontend dan tidak dikirim ke backend.
- React Hook Form memakai `onTouched`: error muncul setelah field kehilangan
  fokus atau form disubmit, lalu diperbarui saat diketik. Petunjuk panjang password
  diperbarui langsung; konfirmasi yang sudah disentuh juga diperiksa ulang ketika
  password berubah.
- Login hanya memerlukan bentuk email valid dan password tidak kosong; tidak
  menerapkan ulang aturan panjang/strength registrasi.
- Label, `aria-invalid`, deskripsi error, show/hide password, status submit, dan
  tombol disabled tersedia. Kegagalan API ditampilkan tanpa menyalin error internal.
- Logout menunggu revocation backend sebelum navigasi. Bila gagal, halaman tetap
  terbuka dan user dapat mencoba lagi.

### Request, cookie, dan proteksi

Browser memanggil `/api/auth/*` pada origin frontend yang sama dengan
`credentials: same-origin`. Route handler menerima hanya kombinasi action/metode
yang diizinkan, meneruskan ke `/auth/*` melalui `INTERNAL_API_URL`, membatasi body
16 KiB, dan memberi timeout upstream 10 detik. Client memakai timeout 15 detik.
Tidak ada URL target dari browser atau redirect upstream yang diikuti.

Proxy meneruskan `Origin` asli tanpa membuat/mengganti nilainya, Content-Type,
dan hanya cookie `ta_session` / `__Host-ta_session`. Respons meneruskan status,
body, Content-Type, dan Set-Cookie session; tidak meneruskan header internal lain.
Semua respons auth memakai `Cache-Control: no-store`. Backend memeriksa Origin
mutation terhadap `AUTH_ALLOWED_ORIGINS` dan mewajibkan JSON. Logout mengirim `{}`.
Tidak ada CORS credentialed lintas origin atau token di localStorage/sessionStorage.

Cookie ditetapkan oleh backend: HttpOnly, SameSite=Lax, Path=/, tanpa Domain,
masa berlaku 7 hari. Development memakai `ta_session` pada HTTP localhost;
production memakai Secure `__Host-ta_session` dengan HTTPS. `INTERNAL_API_URL`
hanya berada di server; contoh Compose adalah `http://backend:8080`.

`getCurrentUser()` memanggil backend `/auth/me` secara langsung dari server dengan
cookie session, `cache: no-store`, dan timeout. Cookie yang tidak ada atau respons
401 berarti guest. Kegagalan jaringan/server menjadi halaman error dengan retry,
bukan dianggap logout. Hanya ID, email, dan waktu pembuatan user yang dipakai.
Future protected pages harus memanggil pemeriksaan ini pada akses data/page,
dan tetap memakai pemeriksaan autentikasi backend secara independen.

## Struktur kode

| Lokasi | Tanggung jawab |
| --- | --- |
| `src/app` | Routes, layout, metadata, loading/error, dan komposisi halaman |
| `src/app/api/auth/[action]/route.ts` | Entry point proxy autentikasi GET/POST |
| `src/features/auth` | Form, shell akun, schema, client API, server auth, proxy, dan tests |
| `src/components/ui` | Primitives shadcn/ui Button, Input, Label, Card, Alert |
| `src/lib/utils.ts` | Helper `cn` untuk class Tailwind |
| `tests/home.test.tsx` | Perilaku guest/authenticated pada home, login/register, dan dashboard |
| `public` | Aset statis publik; jangan menyimpan secret |

Folder placeholder lain tetap tersedia untuk kebutuhan berikutnya. Jangan membuat
provider, service wrapper, atau fitur masa depan hanya untuk mengisi struktur.
Logika domain berada di `src/features/<fitur>`; route merangkai halaman. Komponen
hanya dipindah ke folder bersama jika memang dipakai lintas fitur.

### shadcn/ui dan dependency

Integrasi mengikuti [manual installation shadcn/ui](https://ui.shadcn.com/docs/installation/manual):
`components.json`, alias yang sudah ada, CSS variables pada `globals.css`, serta
helper `cn` (clsx + tailwind-merge). Komponen UI disalin dari registry resmi
`new-york-v4`, dengan import `cn` diarahkan ke helper lokal. Komponen dimiliki
repository dan dapat disesuaikan mengikuti pola shadcn; tidak ada design system
paralel. Radix mendukung primitives, class-variance-authority mendukung variant,
lucide-react menyediakan ikon, dan tw-animate-css mendukung utilitas animasi.

React Hook Form + Zod + resolvers dipakai untuk form. `server-only` mencegah
helper yang membutuhkan konfigurasi internal masuk ke bundle client. Dependency
fitur dipasang dengan versi exact dan tercatat di `bun.lock`.

### Konvensi

- Alias `@/*` menunjuk ke `src/*`; TypeScript strict.
- File/folder proyek menggunakan kebab-case; komponen/tipe PascalCase;
  fungsi/variabel camelCase; hook berawalan `use`. Pertahankan filename framework.
- JSX memakai `.tsx`, kode TypeScript lain `.ts`.
- Secret dan konfigurasi internal tidak boleh masuk variabel `NEXT_PUBLIC_*`.
- Bun adalah package manager. Perubahan `package.json` dan `bun.lock` harus
  bersama; jangan commit dependency, hasil build, `.env`, atau secret.
- Ikuti [konvensi bersama](../docs/ai/CONVENTIONS.md) dan perintah di
  [AGENTS.md](../AGENTS.md). Commit/push hanya atas instruksi user.

## Pemeriksaan

Startup container memasang dependency melalui `bun install --frozen-lockfile`.
Jest dan seluruh dependency form/UI sudah tercatat di manifest/lockfile; tidak
perlu mengulang instalasi paket satu per satu.

Dari root repository dengan container frontend aktif:

```powershell
docker compose exec frontend bun run test
docker compose exec frontend bun run lint
docker compose exec frontend bun run typecheck
```

Production build memakai volume `.next` yang sama dengan development. Hentikan
frontend development selama build, lalu jalankan kembali:

```powershell
docker compose stop frontend
docker compose run --rm --no-deps -e NODE_ENV=production frontend bun --bun run build
docker compose start frontend
```

Jalankan kembali frontend juga bila build gagal. Jangan menjalankan build dan
server development bersamaan pada cache yang sama. Hasil pemeriksaan aktual
tercatat di [CURRENT.md](../docs/ai/CURRENT.md) dan [LOG.md](../docs/ai/LOG.md).

### Testing dengan Jest

Gunakan `bun run test`, bukan `bun test`. Bun memasang paket; script menjalankan
Jest melalui Node.js yang tersedia dalam image Docker. Alternatif tanpa server
Next/backend/database aktif, setelah dependency tersedia:

```powershell
docker compose run --rm --no-deps frontend bun run test
docker compose run --rm --no-deps frontend bun run test:watch
docker compose run --rm --no-deps frontend bun run test:ci
```

`test:coverage` juga tersedia. Coverage berada di `frontend/coverage/` dan tidak
ikut Git. Suite kosong harus gagal; jangan menambah `passWithNoTests`.

- Schema tests mencakup normalisasi email, input tidak valid, batas panjang Unicode,
  konfirmasi password, dan perbedaan aturan login/registrasi.
- Form tests memakai Testing Library dengan role/label untuk validasi setelah
  interaksi, visibility password, loading, error backend, navigasi, dan logout.
- Client tests memeriksa kontrak JSON, tidak terkirimnya konfirmasi, credential
  policy, logout `{}`, dan error generik.
- Server/proxy tests memakai environment Node untuk cookie filtering, no-store,
  method/action allowlist, forwarding Origin, body limit, dan kegagalan upstream.
- Page tests menguji hasil fungsi async page dengan batas auth dimock; ini bukan
  pengganti verifikasi Next.js runtime dan lifecycle cookie di browser.

Test fitur berada dekat source; test halaman/lintas fitur di `tests`. Mock batas
API/router, jangan mock fungsi yang sedang diuji. Gunakan assertion perilaku,
bukan snapshot besar atau class CSS. Test UI, lint, typecheck, dan build tetap
perlu dilengkapi pemeriksaan lifecycle auth melalui backend/database nyata.
