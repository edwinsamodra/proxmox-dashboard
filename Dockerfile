# ─── Stage 1: Build Go backend ───────────────────────────────────────────────
FROM golang:1.21-alpine AS backend-builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /pmo-server \
    ./cmd/server

# ─── Stage 2: Build Vue frontend ─────────────────────────────────────────────
FROM node:20-alpine AS frontend-builder

WORKDIR /build/frontend
COPY frontend/package*.json ./
RUN npm ci --prefer-offline

COPY frontend/ ./
RUN npm run build

# ─── Stage 3: Final minimal image ────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S pmo && \
    adduser -S pmo -G pmo

# Copy Go binary
COPY --from=backend-builder /pmo-server /usr/local/bin/pmo-server

# Copy built frontend into a static directory served by the backend
COPY --from=frontend-builder /build/frontend/dist /app/static

# Default workspace directory for OpenTofu
RUN mkdir -p /app/workspaces && chown pmo:pmo /app /app/workspaces

USER pmo
WORKDIR /app

# Expose HTTP port
EXPOSE 8080

ENV PMO_PORT=8080 \
    PMO_TOFU_DIR=/app/workspaces

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/api/health || exit 1

ENTRYPOINT ["/usr/local/bin/pmo-server"]
