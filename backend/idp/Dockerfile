# Build context is the repo root (single Go module): docker build -f backend/idp/Dockerfile .
FROM --platform=${BUILDPLATFORM} golang:1.24.3-alpine3.21 AS base
WORKDIR /src
ENV CGO_ENABLED=0
ENV GOFLAGS="-mod=readonly"
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    go mod download
COPY gen/go ./gen/go
COPY pkg ./pkg
COPY backend/idp ./backend/idp

FROM base AS build
ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/idp-service ./backend/idp/cmd/server

FROM base AS unit-test
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    mkdir /out && go test -v -coverprofile=/out/cover.out ./backend/idp/...

FROM golangci/golangci-lint:v2.1.6-alpine AS lint-base

FROM base AS lint
ENV GOLANGCI_LINT_CACHE="/root/.cache/golangci-lint"
RUN --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=${GOLANGCI_LINT_CACHE},sharing=locked \
    golangci-lint run --timeout 10m0s ./backend/idp/...

FROM scratch AS unit-test-coverage
COPY --from=unit-test /out/cover.out /cover.out

FROM --platform=${BUILDPLATFORM} alpine:3.21 AS probe
ARG TARGETARCH
RUN case "${TARGETARCH}" in \
      amd64) SHA256=7e98a938983343b00c09a04faf13be0a2cb77fb476410ebca99721598fa2499f ;; \
      arm64) SHA256=43b381aabf5e1f6b2636558b78dfa169ba5b0edb18b031ebbe71558ac82e8789 ;; \
      *) echo "unsupported arch: ${TARGETARCH}" >&2; exit 1 ;; \
    esac \
    && wget -qO /grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.35/grpc_health_probe-linux-${TARGETARCH} \
    && echo "${SHA256}  /grpc_health_probe" | sha256sum -c - \
    && chmod +x /grpc_health_probe

FROM scratch AS bin
COPY --from=probe /grpc_health_probe /
COPY --from=build /out/idp-service /
USER 65532:65532
