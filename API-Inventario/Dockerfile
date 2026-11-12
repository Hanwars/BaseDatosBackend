# Build stage: usa una versión de Go compatible con go.mod (go 1.26.2).
FROM golang:1.27.4 AS builder

WORKDIR /src

# Descargar dependencias por separado para aprovechar la caché de Docker.
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código y compilar un binario estático.
COPY . .
RUN mkdir -p /out \
    && CGO_ENABLED=0 go build \
    -trimpath \
    -buildvcs=false \
    -ldflags="-s -w" \
    -o /out/api \
    .

# Runtime: no incluye Go ni el código fuente.
FROM debian:stable-slim AS runtime

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl tzdata \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system --gid 10001 app \
    && useradd --system --uid 10001 --gid app --home-dir /app --no-create-home app \
    && mkdir -p /data/supplies \
    && chown -R app:app /data

WORKDIR /app
COPY --from=builder --chown=app:app /out/api ./api

ENV GIN_MODE=release \
    SERVER_URL=0.0.0.0:8080 \
    BACKUP_DIR=

USER 10001:10001
EXPOSE 8080

# /health comprueba también la conexión con SQL Server.
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
    CMD ["curl", "-fsS", "http://127.0.0.1:8080/health"]

ENTRYPOINT ["/app/api"]
