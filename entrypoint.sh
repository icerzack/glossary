#!/bin/sh
set -e

# Start backend in background
echo "Starting backend server..."
/usr/local/bin/glossary-server --seed &
BACKEND_PID=$!

# Wait for backend to be ready with healthcheck
echo "Waiting for backend to start..."
MAX_ATTEMPTS=30
ATTEMPT=0

BACKEND_PORT=${PORT:-8080}

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
  if wget --spider --quiet "http://localhost:${BACKEND_PORT}/health" 2>/dev/null; then
    echo "Backend is ready!"
    break
  fi
  ATTEMPT=$((ATTEMPT + 1))
  echo "Attempt $ATTEMPT/$MAX_ATTEMPTS: Backend not ready yet..."
  sleep 1
done

if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
  echo "ERROR: Backend failed to start within ${MAX_ATTEMPTS} seconds"
  exit 1
fi

# Start nginx in foreground
echo "Starting nginx..."
exec nginx -g "daemon off;"

