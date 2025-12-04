# Builder stage
FROM golang:1.23.4-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the manager binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o manager ./cmd/manager/main.go

# Build the migrator binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrator ./cmd/migrator/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install postgresql client for health checks
RUN apk --no-cache add postgresql-client

# Copy binaries from builder
COPY --from=builder /app/manager .
COPY --from=builder /app/migrator .

# Copy migrations
COPY migration ./migration

# Expose port
EXPOSE 8081

# Default command
CMD ["./manager"]
