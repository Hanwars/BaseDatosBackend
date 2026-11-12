#!/bin/sh
set -e

# Puerto: SERVER_URL explícito > $PORT (Render/Heroku-style) > 8080.
if [ -z "${SERVER_URL:-}" ]; then
  SERVER_URL="0.0.0.0:${PORT:-8080}"
fi
export SERVER_URL

# Las imágenes de insumos se guardan en /data/supplies (disco/volumen persistente).
# Un disco nuevo suele montarse como root, así que se corrigen los permisos aquí.
mkdir -p /data/supplies 2>/dev/null || true
if [ "$(id -u)" = "0" ]; then
  chown -R app:app /data 2>/dev/null || true
  exec su-exec app:app /app/api "$@"
fi

exec /app/api "$@"
