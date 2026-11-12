#!/bin/sh
set -eu

# The API runs without root privileges. Images are stored in Supabase Storage.
exec su-exec app:app /app/api "$@"
