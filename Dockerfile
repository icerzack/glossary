FROM golang:1.21-alpine AS backend-builder

WORKDIR /app/backend

RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache gcc musl-dev sqlite-dev

COPY backend/go.mod backend/go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && \
    go install github.com/swaggo/swag/cmd/swag@latest

COPY backend/ ./

RUN swag init -g main.go --output docs

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=1 GOOS=linux \
    CGO_CFLAGS="-D_LARGEFILE64_SOURCE" \
    go build -a -installsuffix cgo -o /app/glossary-server main.go

FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package*.json ./

RUN --mount=type=cache,target=/root/.npm \
    npm install

COPY frontend/ ./
RUN npm run build

FROM nginx:alpine

RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache sqlite-libs gettext

COPY nginx.conf /etc/nginx/nginx.conf

COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html

COPY --from=backend-builder /app/glossary-server /usr/local/bin/glossary-server

RUN mkdir -p /data

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${NGINX_PORT:-3000}/health || exit 1

ENTRYPOINT ["/entrypoint.sh"]

