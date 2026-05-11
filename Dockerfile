# ── Stage 1: Next.js deps ────────────────────────────────────────────────────
FROM node:20-alpine AS frontend-deps
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci

# ── Stage 2: Next.js build ───────────────────────────────────────────────────
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY --from=frontend-deps /app/node_modules ./node_modules
COPY frontend/ .
RUN mkdir -p public
ARG NEXT_PUBLIC_API_URL=/api/v1
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
RUN npm run build

# ── Stage 3: Go build ────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS backend-builder
RUN apk add --no-cache git
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o pharmasense-api ./cmd/api

# ── Stage 4: Runtime ─────────────────────────────────────────────────────────
FROM node:20-alpine
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

# Go binary
COPY --from=backend-builder /app/pharmasense-api ./pharmasense-api

# Next.js standalone bundle
COPY --from=frontend-builder /app/.next/standalone/ ./frontend/
COPY --from=frontend-builder /app/.next/static      ./frontend/.next/static
COPY --from=frontend-builder /app/public            ./frontend/public

COPY start.sh ./start.sh
RUN chmod +x ./start.sh

ENV NODE_ENV=production
ENV PORT=8080
ENV NEXT_PORT=3001
ENV HOSTNAME=0.0.0.0

EXPOSE 8080

CMD ["./start.sh"]
