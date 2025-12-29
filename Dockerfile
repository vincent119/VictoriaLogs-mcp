# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -trimpath -o vlmcp ./cmd/vlmcp

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/vlmcp .
COPY --from=builder /app/configs ./configs

# Expose port (if TCP transport is used, though stdio is default for MCP)
# EXPOSE 8080

ENTRYPOINT ["./vlmcp"]
