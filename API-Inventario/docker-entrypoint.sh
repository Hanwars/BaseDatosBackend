#!/bin/sh
set -e

# Las imágenes de insumos se guardan en /data/supplies (volumen persistente en Northflank).
# Un volumen nuevo suele montarse como root, así que se corrigen los permisos aquí.
mkdir -p /data/supplies
if [ "$(id -u)" = "0" ]; then
  chown -R app:app /data 2>/dev/null || true
  exec su-exec app:app /app/api "$@"
fi

exec /app/api "$@"
