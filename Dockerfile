# Multi-stage Dockerfile for Glossary Application
# Stage 1: Build Go backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app/backend

# Install build dependencies with cache mount
RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod files and download dependencies with cache
COPY backend/go.mod backend/go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Copy backend source
COPY backend/ ./

# Generate Swagger documentation
RUN swag init -g main.go --output docs

# Build the application with cache
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /app/glossary-server main.go

# Stage 2: Build React frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files
COPY frontend/package*.json ./

# Install dependencies with cache
RUN --mount=type=cache,target=/root/.npm \
    npm install

# Copy frontend source
COPY frontend/ ./

# Build the application
RUN npm run build

# Stage 3: Final image with nginx
FROM nginx:alpine

# Install sqlite libs with cache mount
RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache sqlite-libs

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Copy frontend build
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html

# Copy backend binary
COPY --from=backend-builder /app/glossary-server /usr/local/bin/glossary-server

# Create directory for database
RUN mkdir -p /data

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose port
EXPOSE 63000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:63001/health || exit 1

# Start services
ENTRYPOINT ["/entrypoint.sh"]

