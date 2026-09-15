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
| `tests` | Pengujian lintas fitur atau integrasi | Skenario alur membuat trading plan setelah runner tersedia |
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
   Tambahkan pengujian otomatis yang sesuai ketika test runner sudah tersedia.
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
        │   └── trading-plan-form.tsx
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
docker compose exec frontend bun --bun run build
```

Belum ada test runner atau pengujian fitur. Commit `bun.lock` bersama perubahan
`package.json`. Gunakan Bun secara konsisten sebagai package manager.

Lint memeriksa aturan kode, typecheck memeriksa tipe, dan build memeriksa apakah
aplikasi bisa dibangun. Ketiganya tidak menggantikan pengujian perilaku di browser.
Jika pemeriksaan gagal, perbaiki penyebabnya; jangan mematikan aturan hanya agar lolos.

## Yang sudah tersedia dan yang ditunda

Fondasi yang tersedia: scaffold Next.js, TypeScript strict, ESLint Next.js/TypeScript,
Tailwind, Bun lockfile, konfigurasi Docker development, dan folder kerja.
Status lulus lint/typecheck/build perlu dibuktikan dengan menjalankan perintah di atas.

Prettier, Husky, lint-staged, test runner, autentikasi, API client, state/query library,
serta UI bisnis belum disiapkan. Tambahkan bertahap sesuai task. README ini cukup
untuk petunjuk teknis awal; PRD dan roadmap tetap di luar repository.
