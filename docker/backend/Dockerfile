ARG GO_VERSION=1.27.1
FROM golang:${GO_VERSION}-alpine

ARG AIR_VERSION=v1.67.3

RUN apk add --no-cache ca-certificates git \
    && go install github.com/air-verse/air@${AIR_VERSION}

WORKDIR /app

COPY air.toml /etc/air/air.toml

EXPOSE 8080

CMD ["sh", "-c", "if [ ! -f go.mod ]; then echo 'backend/go.mod belum ada. Inisialisasi proyek Go terlebih dahulu.' >&2; exit 1; fi; exec air -c /etc/air/air.toml"]
