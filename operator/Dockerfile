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
COPY api/ api/
COPY controllers/ controllers/
COPY internal/ internal/
COPY cmd/ cmd/

# Build
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -o manager cmd/manager/main.go

# Final stage
FROM alpine:latest

WORKDIR /

# Install ca-certificates for HTTPS
# Use --no-cache and update index first for better multi-platform compatibility
RUN apk update && apk --no-cache add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /workspace/manager .

ENTRYPOINT ["/manager"]

