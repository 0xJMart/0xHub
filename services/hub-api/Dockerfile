# syntax=docker/dockerfile:1.6

ARG GO_VERSION=1.24.10
FROM golang:${GO_VERSION}-bullseye AS builder

WORKDIR /workspace

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/workspace/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -o /out/hub-api ./cmd/hub-api

FROM gcr.io/distroless/base-debian12:nonroot AS runtime

ENV PORT=8080
EXPOSE 8080

WORKDIR /app
COPY --from=builder /out/hub-api /app/hub-api

ENTRYPOINT ["/app/hub-api"]

