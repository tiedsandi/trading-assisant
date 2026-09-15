# Trading Assistant

## Overview
Trading Assistant adalah aplikasi pribadi untuk membantu proses trading
Futures secara lebih terstruktur, disiplin, dan berbasis trading plan.

Aplikasi tidak bertujuan menggantikan keputusan trader.
Aplikasi membantu mengumpulkan data pasar, melakukan screening,
menyusun analisis, memvalidasi setup melalui checklist, mengatur risiko,
dan mencatat hasil trading.

## Problem

Trading sering dilakukan dengan proses yang tidak konsisten:

- pemilihan coin tidak terstruktur
- analisis berbeda-beda setiap trade
- entry terlambat atau terlalu cepat
- risk management tidak konsisten
- keputusan trading sulit dievaluasi kembali
- histori kesalahan dan keberhasilan tidak terdokumentasi dengan baik

Trading Assistant dibuat untuk membuat proses tersebut repeatable
dan measurable.

## Core Workflow

Market
→ Screening
→ Analysis
→ Trading Checklist
→ Trading Plan
→ Entry Decision
→ Position
→ Exit
→ Journal
→ Evaluation

Setiap tahap nantinya dapat dikembangkan secara bertahap.

## Product Principles

### Decision Support, Not Decision Replacement
Aplikasi membantu memberikan informasi dan struktur.
Keputusan entry/exit tetap berada pada user.

### Risk First
Risk management harus menjadi bagian dari trading plan,
bukan keputusan setelah entry.

### Evidence Over Feeling
Keputusan trading sebisa mungkin mempunyai alasan yang dapat dicatat
dan dievaluasi.

### Build Incrementally
Tidak mencoba membangun seluruh trading platform sekaligus.
Setiap feature harus usable sebelum tahap berikutnya dibuat.

## Product Stages

### Stage 1 — Trading Foundation
Membangun workflow trading dasar.

Target:
- market data
- coin screening
- watchlist
- trading checklist
- trading plan
- risk calculation
- trading journal

Tujuan:
Aplikasi sudah berguna untuk membantu trading sehari-hari tanpa AI.

### Stage 2 — AI Trading Coach
AI menggunakan data dan histori trading untuk membantu evaluasi.

Contoh:
- review trading plan
- mendeteksi pelanggaran rule
- memberikan feedback terhadap setup
- menganalisis pola kesalahan
- post-trade review

AI berfungsi sebagai coach, bukan autonomous trader.

### Stage 3 — News & Market Narrative
Menambahkan konteks fundamental dan narrative pasar.

Contoh:
- market news aggregation
- coin-specific news
- narrative/trend detection
- event awareness
- menghubungkan technical setup dengan market context

## Long-Term Direction

Trading Assistant berkembang menjadi personal trading workspace yang
menghubungkan:

Market Data
+
Technical Analysis
+
Risk Management
+
Trading Journal
+
AI Coaching
+
Market Context

Tujuan akhirnya bukan menghasilkan sinyal sebanyak mungkin,
tetapi membantu user mengambil keputusan trading yang lebih
konsisten, terukur, dan dapat dievaluasi.

## Current Scope

Pengembangan dilakukan bertahap.

Detail feature yang sedang dikerjakan tidak disimpan di dokumen ini.
Lihat:

- `docs/ai/CURRENT.md` untuk kondisi project saat ini.
- `docs/ai/LOG.md` untuk histori perubahan penting.
- `WORKFLOW.md` untuk development workflow.

Dokumen ini hanya menyimpan arah produk yang relatif stabil.