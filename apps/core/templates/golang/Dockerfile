# ==========================
# Builder Stage
# ==========================
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Download dependencies trước để tận dụng Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o main .

# ==========================
# Runtime Stage
# ==========================
FROM alpine:3.22

WORKDIR /app

# Cài CA certificates nếu service gọi HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary từ builder
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]