# Build stage
FROM golang:1.25.5-alpine3.23 AS builder

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

# Create a default config file to ensure startup success
RUN echo "server:" > config.yaml && \
    echo "  name: \"victorialogs-mcp\"" >> config.yaml && \
    echo "  version: \"1.0.0\"" >> config.yaml && \
    echo "  transport: \"sse\"" >> config.yaml && \
    echo "  tcp_addr: \":8080\"" >> config.yaml && \
    echo "victorialogs:" >> config.yaml && \
    echo "  url: \"http://localhost:9428\"" >> config.yaml && \
    echo "  timeout: \"30s\"" >> config.yaml && \
    echo "  auth:" >> config.yaml && \
    echo "    type: \"none\"" >> config.yaml && \
    echo "logging:" >> config.yaml && \
    echo "  level: \"info\"" >> config.yaml

# Expose port (if TCP transport is used, though stdio is default for MCP)
# EXPOSE 8080

ENTRYPOINT ["./vlmcp"]
