# Stage 1 — Frontend build
FROM node:22-alpine AS frontend-builder

ARG VITE_BASE_URL
ENV VITE_BASE_URL=$VITE_BASE_URL

WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# Stage 2 — Go build
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/service ./cmd/api/main.go

# Stage 3 — Runtime
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder /app/bin/service /app/service
COPY --from=frontend-builder /app/frontend/dist /app/frontend/dist

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["/app/service"]
