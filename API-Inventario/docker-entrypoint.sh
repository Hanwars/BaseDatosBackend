#!/bin/sh
set -eu

# Las imágenes se guardan en /data/supplies.
mkdir -p /data/supplies

# Render puede montar un disco en /data. Como el volumen puede pertenecer a root,
# se ajustan permisos antes de bajar privilegios para ejecutar la API.
if [ "$(id -u)" -eq 0 ]; then
    chown -R app:app /data
    exec su-exec app:app /app/api "$@"
fi

exec /app/api "$@"
