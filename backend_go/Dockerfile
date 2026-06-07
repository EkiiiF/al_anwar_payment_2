# Stage 1: Build
FROM golang:1.26-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Optimasi binary: hapus debug info (-s -w) dan static build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/server/main.go

# Stage 2: Run
FROM alpine:3.20
RUN apk add --no-cache tzdata ca-certificates
WORKDIR /app
# Keamanan: User non-root
RUN adduser -D -u 1001 appuser
USER appuser
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
# Port sesuai konfigurasi Fiber
EXPOSE 8080
ENTRYPOINT ["./main"]