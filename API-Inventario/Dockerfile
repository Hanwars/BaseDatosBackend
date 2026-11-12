# syntax=docker/dockerfile:1.7

# Build: cachea módulos antes de copiar el código para acelerar reconstrucciones.
FROM golang:1.27.4-alpine AS build
WORKDIR /src

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -buildvcs=false \
    -ldflags="-s -w" \
    -o /out/api ./main.go

# Imagen de ejecución pequeña, sin compilador ni código fuente.
FROM alpine:3.22 AS runtime
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S -g 10001 app \
    && adduser -S -D -H -u 10001 -G app app \
    && mkdir -p /data/supplies \
    && chown -R app:app /data

WORKDIR /app
COPY --from=build --chown=app:app /out/api ./api

ENV GIN_MODE=release \
    SERVER_URL=0.0.0.0:8080 \
    BACKUP_DIR=

USER 10001:10001
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=15s --retries=3 \
  CMD wget -q -O /dev/null http://127.0.0.1:8080/health || exit 1

ENTRYPOINT ["/app/api"]
