# Build context is the repo root (single Go module): docker build -f bff/Dockerfile --build-arg VARIANT=legacy .
FROM --platform=${BUILDPLATFORM} golang:1.24.3-alpine3.21 AS base
WORKDIR /src
ENV CGO_ENABLED=0
ENV GOFLAGS="-mod=readonly"
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    go mod download
COPY gen/go ./gen/go
COPY pkg ./pkg
COPY bff ./bff

FROM base AS build
ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/handson-entry ./bff

FROM golangci/golangci-lint:v2.1.6-alpine AS lint-base

FROM base AS lint
ENV GOLANGCI_LINT_CACHE="/root/.cache/golangci-lint"
RUN --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=${GOLANGCI_LINT_CACHE},sharing=locked \
    golangci-lint run --timeout 10m0s ./bff/...

FROM scratch AS bin
ARG VARIANT
ENV VARIANT=${VARIANT}
COPY --from=build /out/handson-entry /
USER 65532:65532
ENTRYPOINT ["/handson-entry"]
