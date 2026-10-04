# syntax=docker/dockerfile:1.7
# Builds one image with the three Habarchy binaries; pick one with the
# container command: api (default) | worker | scheduler.
FROM golang:1.24-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates tzdata
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
ARG VERSION=dev
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOFLAGS=-trimpath go build -ldflags="-s -w -X main.version=${VERSION}" -o /out/api ./cmd/api && \
    CGO_ENABLED=0 GOFLAGS=-trimpath go build -ldflags="-s -w" -o /out/worker ./cmd/worker && \
    CGO_ENABLED=0 GOFLAGS=-trimpath go build -ldflags="-s -w" -o /out/scheduler ./cmd/scheduler

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata curl && adduser -D -u 10001 habarchy
COPY --from=build /out/api /out/worker /out/scheduler /usr/local/bin/
USER habarchy
ENV HABARCHY_HTTP_ADDR=:8080 HABARCHY_LOG_PRETTY=false
EXPOSE 8080 9090
HEALTHCHECK --interval=15s --timeout=3s --start-period=10s CMD curl -fsS http://127.0.0.1:8080/healthz || exit 1
ENTRYPOINT []
CMD ["api"]
