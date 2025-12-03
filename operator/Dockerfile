# Build stage
FROM golang:1.24-alpine AS builder

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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager cmd/manager/main.go

# Final stage
FROM alpine:latest

WORKDIR /
COPY --from=builder /workspace/manager .

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

ENTRYPOINT ["/manager"]

