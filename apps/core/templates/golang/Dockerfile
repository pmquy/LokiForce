# ==========================
# Builder Stage
# ==========================
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy source and resolve dependencies
COPY . .
RUN go mod tidy && go mod download

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