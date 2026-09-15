# Frontend

Next.js App Router, TypeScript, Tailwind CSS, dan Bun. Cara menjalankan Docker
ada di [README root](../README.md).

## Struktur kode

| Folder | Fungsi | Contoh file saat diperlukan |
| --- | --- | --- |
| `src/app` | Routes, layout, metadata, dan komposisi halaman | `trading-plan/page.tsx` merangkai tampilan fitur |
| `src/components/ui` | Komponen UI dasar yang tidak mengenal domain trading | `button.tsx`, `input.tsx` |
| `src/components/common` | Komponen generik gabungan, dipakai lintas fitur | `empty-state.tsx`, `error-message.tsx` |
| `src/components/layout` | Kerangka tampilan aplikasi | `app-header.tsx`, `app-sidebar.tsx` |
| `src/features` | Kode khusus suatu fitur | `trading-plan/components/trading-plan-form.tsx` |
| `src/hooks` | Hook generik lintas fitur | `use-debounce.ts` |
| `src/lib` | Infrastruktur dan helper generik | `api-client.ts` untuk transport HTTP |
| `src/config` | Konfigurasi dan konstanta aplikasi | `navigation.ts` untuk daftar menu |
| `src/types` | Tipe generik yang benar-benar dipakai bersama | `pagination.ts` |
| `public` | Aset statis yang boleh diakses publik | Gambar dan ikon; jangan simpan secret |
| `tests` | Test halaman atau skenario lintas fitur | `home.test.tsx` untuk halaman awal |
| `documentations` | Penjelasan keputusan/pola teknis yang tidak cukup jelas dari kode | `api-client.md` setelah integrasi API dibuat |

Folder kosong memakai `.gitkeep` agar ikut Git. Tambahkan subfolder fitur
`auth`, `trading-plan`, `risk`, `journal`, dan `market` ketika mulai dikerjakan.
Route groups `(auth)` dan `(protected)`, route `api`, `providers.tsx`, serta
`proxy.ts` dibuat saat implementasi membutuhkannya. Nama `(protected)` sendiri
bukan mekanisme pemeriksaan akses.

Logika bisnis khusus fitur berada di `src/features/<fitur>`, bukan di route.
Komponen/hook/tipe tetap dekat fitur sampai benar-benar dibutuhkan bersama.
Alias `@/*` menunjuk ke `src/*`. Perhitungan domain trading nantinya di backend.
Halaman dan aset saat ini masih bawaan Next.js.

## Cara mengerjakan satu fitur

1. Mulai dari task kecil: tentukan halaman, perilaku yang diinginkan, dan kriteria
   selesai. Jika membutuhkan backend, sepakati request, response, dan error API.
2. Periksa kode yang sudah ada. Buat folder `src/features/<nama-fitur>` hanya
   untuk fitur yang sedang dikerjakan; gunakan kembali komponen yang sesuai.
3. Letakkan komponen, tipe, hook, dan pemanggilan API khusus fitur di folder fitur.
   Untuk fitur sederhana, beberapa file langsung di folder fitur sudah cukup.
   Tambahkan subfolder ketika jumlah file memang membutuhkannya.
4. Tambahkan route di `src/app` untuk menampilkan fitur. Route bertugas merangkai
   halaman; logika khusus fitur tetap berada di folder fitur.
5. Implementasikan kondisi loading, data kosong, error, validasi input, dan
   interaksi yang relevan. Validasi frontend membantu UX; backend tetap menjadi
   sumber validasi dan perhitungan bisnis trading.
6. Jalankan pemeriksaan di bawah, lalu uji alur normal dan error di browser.
   Tambahkan test perilaku yang relevan menggunakan Jest dan React Testing Library.
7. Review perubahan, perbarui dokumentasi jika pola teknis berubah, lalu commit
   hanya file terkait. Jangan commit `.env`, dependency, atau hasil build.

Contoh penempatan untuk task trading plan yang membutuhkan form dan API:

```text
src/
├── app/
│   └── trading-plan/
│       └── page.tsx
└── features/
    └── trading-plan/
        ├── components/
        │   ├── trading-plan-form.tsx
        │   └── trading-plan-form.test.tsx
        ├── api/
        │   └── create-trading-plan.ts
        └── types.ts
```

Ini contoh, bukan daftar file yang wajib dibuat sekarang. Hook seperti
`use-trading-plan.ts` hanya dibuat jika ada perilaku yang perlu dipisahkan.
`src/lib/api-client.ts` menangani transport HTTP umum; file di `features/.../api`
menangani endpoint dan data khusus fitur. Jangan membuat lapisan service/helper
yang hanya meneruskan satu pemanggilan tanpa tanggung jawab tambahan.

## Penamaan dan batas tanggung jawab

- Folder dan file buatan proyek memakai `kebab-case`, misalnya
  `trading-plan-form.tsx`. Pertahankan nama file khusus framework seperti `page.tsx`.
- Nama komponen dan tipe memakai `PascalCase`, misalnya `TradingPlanForm` dan
  `TradingPlan`. Fungsi serta variabel memakai `camelCase`, misalnya `createTradingPlan`.
- Nama hook diawali `use`, misalnya `useTradingPlan` di `use-trading-plan.ts`.
- Gunakan `.tsx` untuk file dengan JSX dan `.ts` untuk file TypeScript lainnya.
- Komponen khusus trading plan tetap di fitur tersebut meskipun digunakan oleh
  beberapa halaman. Pindahkan ke `components` bersama hanya jika sudah generik.
- Jangan duplikasi tipe khusus fitur ke `src/types`. Hindari import detail internal
  fitur lain; pisahkan bagian bersama jika memang dibutuhkan lintas fitur.
- Secret dan kredensial database tidak boleh berada di kode browser atau variabel
  `NEXT_PUBLIC_*`. Akses database ditangani backend Go.
- Ikuti gaya kode yang sudah ada. Belum ada aturan khusus untuk memaksa arrow
  function menjadi deklarasi `function`, dan belum ada formatter otomatis.

Panduan ini adalah konvensi kerja, bukan aturan yang semuanya diperiksa ESLint.
Jangan menambah dependency, provider, auth guard, atau folder fitur masa depan
hanya untuk melengkapi struktur contoh.

## Pemeriksaan

Dari root repository, saat frontend berjalan:

```powershell
docker compose exec frontend bun run lint
docker compose exec frontend bun run typecheck
docker compose run --rm --no-deps -e NODE_ENV=production frontend bun --bun run build
```

Konfigurasi Jest sudah tersedia; lakukan instalasi satu kali di bagian Testing.
Commit `bun.lock` bersama perubahan `package.json`. Gunakan Bun secara konsisten
sebagai package manager.

Lint memeriksa aturan kode, typecheck memeriksa tipe, dan build memeriksa apakah
aplikasi bisa dibangun. Ketiganya tidak menggantikan pengujian perilaku di browser.
Jika pemeriksaan gagal, perbaiki penyebabnya; jangan mematikan aturan hanya agar lolos.

## Yang sudah tersedia dan yang ditunda

Fondasi yang tersedia: scaffold Next.js, TypeScript strict, ESLint Next.js/TypeScript,
Tailwind, Bun lockfile, konfigurasi Docker development, dan folder kerja.
Status lulus lint/typecheck/build perlu dibuktikan dengan menjalankan perintah di atas.

Prettier, Husky, lint-staged, autentikasi, API client, state/query library,
serta UI bisnis belum disiapkan. Tambahkan bertahap sesuai task. README ini cukup
untuk petunjuk teknis awal; PRD dan roadmap tetap di luar repository.

## Testing dengan Jest

Jest + React Testing Library menggunakan konfigurasi `next/jest` dan lingkungan
DOM jsdom. Bun tetap memasang paket. Script test secara eksplisit menjalankan
Jest dengan Node.js, yang disediakan oleh image Docker frontend.
Jalankan `bun run test`, bukan `bun test` (runner Bun yang berbeda).

Instalasi satu kali, dari root repository (oleh pengguna):

```powershell
docker compose run --rm --no-deps --build frontend bun add -d jest@30 jest-environment-jsdom@30 @types/jest@30 @testing-library/react@16 @testing-library/dom@10 @testing-library/jest-dom@6
```

Perintah ini memperbarui `package.json` dan `bun.lock` serta memasang dependency
ke volume Docker. Keduanya perlu di-commit bersama. Konfigurasi saja belum cukup
untuk menjalankan test; dependency belum dipasang oleh asisten.
Pada clone setelah lockfile diperbarui, gunakan `bun install --frozen-lockfile`
di container, bukan mengulang `bun add`.

Jalankan test tanpa memerlukan server Next.js/backend/database aktif:

```powershell
docker compose run --rm --no-deps frontend bun run test
docker compose run --rm --no-deps frontend bun run test:watch
docker compose run --rm --no-deps frontend bun run test:ci
```

`test` berjalan sekali; `test:watch` mengulang saat source berubah (`Ctrl+C` untuk
berhenti). `test:ci` berjalan sekali dengan coverage dan gagal jika test gagal.
Coverage disimpan di `frontend/coverage/`, tidak ikut Git. Tidak ada opsi
`passWithNoTests`: suite kosong harus gagal. `test:coverage` tersedia untuk laporan lokal.

`tests/home.test.tsx` adalah smoke test scaffold halaman awal, bukan bukti fitur
bisnis telah benar. Ganti assertion-nya saat halaman berubah. Untuk kode khusus
fitur, letakkan `*.test.ts(x)` dekat source; skenario lintas komponen di `tests/`.
Semua test menggunakan alias `@/*` yang sama dengan aplikasi.

Test komponen, hook, dan fungsi khusus fitur berada di
`src/features/<nama-fitur>/`, bersebelahan dengan file yang diuji.
Test komponen bersama juga bersebelahan dengan source di `src/components/`.
Folder `tests/` bukan tempat wajib untuk semua test. Jest mencari file
`*.test.ts` dan `*.test.tsx` di kedua lokasi tanpa perubahan konfigurasi.

Uji perilaku yang terlihat pengguna, memakai role/label. Untuk fitur baru,
uji input valid/tidak valid, loading, kosong, sukses, dan error yang relevan.
Mock batas eksternal seperti API; jangan mock fungsi yang justru sedang diuji.
Hindari snapshot besar dan assertion terhadap class CSS untuk membuktikan perilaku.
Coverage adalah petunjuk bagian yang belum diuji, bukan bukti bebas bug.

Layout tidak masuk coverage komponen saat ini. Async Server Components dan alur
browser nyata memerlukan strategi E2E berikutnya; suite ini belum mencakupnya.
Lihat [panduan Jest Next.js](https://nextjs.org/docs/app/guides/testing/jest).
Hasil PASS belum dikonfirmasi sampai perintah dijalankan oleh pengguna.
