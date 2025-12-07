# Build stage
FROM golang:1.24-alpine AS builder

# Declare build args for multi-platform support
ARG TARGETARCH
ARG TARGETOS=linux

WORKDIR /workspace

# Copy go mod files
COPY go.mod go.mod
COPY go.sum go.sum

# Download dependencies
RUN go mod download

# Copy source code
COPY cmd/ cmd/
COPY internal/ internal/

# Build
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -o server cmd/server/main.go

# Final stage - use distroless or minimal alpine
FROM alpine:latest

WORKDIR /

# Copy ca-certificates from builder stage (already has them from golang base image)
# This avoids apk install issues during ARM64 emulation
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /workspace/server .

EXPOSE 8080

ENTRYPOINT ["/server"]

