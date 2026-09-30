FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o project-service ./cmd/server/main.go

FROM alpine:3.21 AS probe
ARG TARGETARCH
RUN wget -qO /grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.35/grpc_health_probe-linux-${TARGETARCH} \
    && chmod +x /grpc_health_probe

FROM alpine:3.18

WORKDIR /app

RUN apk --no-cache add ca-certificates
COPY --from=builder /app/project-service .
COPY --from=probe /grpc_health_probe /grpc_health_probe

EXPOSE 50053

CMD ["./project-service", "server"]