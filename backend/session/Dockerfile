FROM --platform=${BUILDPLATFORM} golang:1.24.1-alpine3.21 AS base
WORKDIR /src
ENV CGO_ENABLED=0
ENV GOFLAGS="-mod=readonly"
COPY go.mod .
COPY go.sum . 
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    go mod download

FROM base AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/session-service ./cmd/server/main.go
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/session-client ./cmd/client/main.go

FROM base AS unit-test
WORKDIR /src
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    mkdir /out && go test -v -coverprofile=/out/cover.out ./...

FROM golangci/golangci-lint:v1.55-alpine AS lint-base

FROM base AS lint
WORKDIR /src
ENV GOLANGCI_LINT_CACHE="/root/.cache/golangci-lint"
RUN --mount=target=. \
    --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=${GOLANGCI_LINT_CACHE},sharing=locked \
    golangci-lint run --timeout 10m0s ./...

FROM scratch AS unit-test-coverage
COPY --from=unit-test /out/cover.out /cover.out

FROM scratch AS bin
COPY --from=build /out/session-service /

FROM scratch AS client
COPY --from=build /out/session-client /
