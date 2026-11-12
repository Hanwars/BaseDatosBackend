# ---------- Build ----------
# go.mod exige Go >= 1.26.2; la etiqueta 1.26 resuelve al último parche 1.26.x.
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Dependencias primero para aprovechar la caché de capas.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
      -trimpath -buildvcs=false -ldflags="-s -w" \
      -o /out/api .

# ---------- Runtime ----------
FROM alpine:3 AS runtime

# ca-certificates: TLS hacia Azure SQL. tzdata: zonas horarias. su-exec: bajar privilegios.
RUN apk add --no-cache ca-certificates tzdata su-exec \
    && addgroup -S -g 10001 app \
    && adduser -S -u 10001 -G app -h /app -H app \
    && mkdir -p /app /data/supplies \
    && chown -R app:app /app /data

WORKDIR /app
COPY --from=builder /out/api /app/api

# El script arranca como root solo para arreglar permisos del volumen montado en /data
# y enseguida ejecuta la API como usuario "app" (sin privilegios).
# sed elimina CRLF por si el archivo se editó en Windows.
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN sed -i 's/\r$//' /usr/local/bin/docker-entrypoint.sh \
    && chmod +x /usr/local/bin/docker-entrypoint.sh

ENV GIN_MODE=release \
    SERVER_URL=0.0.0.0:8080

EXPOSE 8080

# Sin HEALTHCHECK de Docker: Northflank usa sus propios health checks
# (HTTP GET /health en el puerto 8080).
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
