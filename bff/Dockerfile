# Build context is the repo root (Go module + pnpm workspace): docker build -f bff/Dockerfile --build-arg VARIANT=legacy .
FROM --platform=${BUILDPLATFORM} node:24.21.0-alpine3.23 AS web
WORKDIR /src
RUN corepack enable
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY frontend/package.json ./frontend/
RUN pnpm install --frozen-lockfile --filter app
COPY frontend ./frontend
ARG VARIANT
RUN VITE_VARIANT=${VARIANT} pnpm --filter app build

FROM --platform=${BUILDPLATFORM} golang:1.25.14-alpine3.23 AS base
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

FROM golangci/golangci-lint:v2.14.0-alpine AS lint-base

FROM base AS lint
ENV GOLANGCI_LINT_CACHE="/root/.cache/golangci-lint"
RUN --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=${GOLANGCI_LINT_CACHE},sharing=locked \
    golangci-lint run --timeout 10m0s ./bff/...

FROM scratch AS bin
ARG VARIANT
ENV VARIANT=${VARIANT}
COPY --from=build /out/handson-entry /
COPY --from=web /src/frontend/dist /web
USER 65532:65532
ENTRYPOINT ["/handson-entry"]
