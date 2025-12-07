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

# Final stage
FROM alpine:latest

WORKDIR /

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

COPY --from=builder /workspace/server .

EXPOSE 8080

ENTRYPOINT ["/server"]

