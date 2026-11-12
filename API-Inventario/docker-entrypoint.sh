
#!/bin/sh
set -e

mkdir -p /data/supplies

if [ "$(id -u)" = "0" ]; then
    chown -R app:app /data 2>/dev/null || true
    exec su-exec app:app /app/api "$@"
fi

exec /app/api "$@"
