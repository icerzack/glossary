#!/bin/sh
set -e

BACKEND_PORT=${PORT:-8080}
NGINX_PORT=${NGINX_PORT:-3000}

export BACKEND_PORT
export NGINX_PORT

envsubst '${BACKEND_PORT} ${NGINX_PORT}' < /etc/nginx/nginx.conf > /tmp/nginx.conf
mv /tmp/nginx.conf /etc/nginx/nginx.conf

echo "Starting backend server on port ${BACKEND_PORT}..."
/usr/local/bin/glossary-server --seed &
BACKEND_PID=$!

echo "Waiting for backend to start..."
MAX_ATTEMPTS=30
ATTEMPT=0

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

echo "Starting nginx..."
exec nginx -g "daemon off;"

