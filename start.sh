#!/bin/sh
set -e

# Start Next.js on the internal port (not exposed outside the container).
PORT=${NEXT_PORT:-3000} HOSTNAME=127.0.0.1 node /app/frontend/server.js &

# Start the Go API in the foreground on $PORT.
# Go reverse-proxies all non-/api/v1 requests to the Next.js server above.
exec /app/pharmasense-api
