# ---------- Build ----------
# go.mod declara Go 1.26.2; Go 1.27.4 es compatible y contiene las correcciones actuales.
FROM golang:1.27.4-alpine3.23 AS builder

WORKDIR /src

# Dependencias primero para aprovechar la caché de capas de Render/BuildKit.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir -p /out \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
       -trimpath -buildvcs=false -ldflags="-s -w" \
       -o /out/api .

# ---------- Runtime ----------
FROM alpine:3.23 AS runtime

RUN apk add --no-cache ca-certificates tzdata su-exec \
    && addgroup -S -g 10001 app \
    && adduser -S -D -H -u 10001 -G app app \
    && mkdir -p /app /data/supplies \
    && chown -R app:app /app /data

WORKDIR /app
COPY --from=builder /out/api /app/api
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

# Normaliza finales de línea por si el archivo se editó desde Windows.
RUN sed -i 's/\r$//' /usr/local/bin/docker-entrypoint.sh \
    && chmod 755 /usr/local/bin/docker-entrypoint.sh \
    && chmod 755 /app/api

ENV GIN_MODE=release \
    PORT=10000 \
    BACKUP_DIR=

# Render usa PORT en ejecución (10000 por defecto).
EXPOSE 10000

# El entrypoint corrige permisos de /data y luego ejecuta la API sin privilegios.
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
