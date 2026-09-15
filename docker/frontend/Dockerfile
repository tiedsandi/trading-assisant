ARG BUN_VERSION=1.4.2
FROM oven/bun:${BUN_VERSION}-debian

WORKDIR /app

ENV NODE_ENV=development
ENV NEXT_TELEMETRY_DISABLED=1

EXPOSE 3000

CMD ["sh", "-ec", "if [ ! -f package.json ] || [ ! -f bun.lock ]; then echo 'frontend/package.json atau frontend/bun.lock belum ada. Inisialisasi proyek Next.js dengan Bun terlebih dahulu.' >&2; exit 1; fi; bun install --frozen-lockfile; exec bun --bun run dev --hostname 0.0.0.0"]
