FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o project-service ./cmd/server/main.go

FROM alpine:3.21 AS probe
ARG TARGETARCH
RUN case "${TARGETARCH}" in \
      amd64) SHA256=7e98a938983343b00c09a04faf13be0a2cb77fb476410ebca99721598fa2499f ;; \
      arm64) SHA256=43b381aabf5e1f6b2636558b78dfa169ba5b0edb18b031ebbe71558ac82e8789 ;; \
      *) echo "unsupported arch: ${TARGETARCH}" >&2; exit 1 ;; \
    esac \
    && wget -qO /grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.35/grpc_health_probe-linux-${TARGETARCH} \
    && echo "${SHA256}  /grpc_health_probe" | sha256sum -c - \
    && chmod +x /grpc_health_probe

FROM alpine:3.18

WORKDIR /app

RUN apk --no-cache add ca-certificates
COPY --from=builder /app/project-service .
COPY --from=probe /grpc_health_probe /grpc_health_probe

EXPOSE 50053

CMD ["./project-service", "server"]