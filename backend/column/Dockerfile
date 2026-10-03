FROM python:3.12-slim AS lint

WORKDIR /app
COPY . /app

COPY requirements-dev.txt /tmp/requirements-dev.txt
RUN pip install --no-cache-dir -r /tmp/requirements-dev.txt
RUN black . --check --exclude="(venv|proto|gen)"
RUN flake8 . --exclude=venv,proto,gen --ignore=E501

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

FROM python:3.12-slim AS runtime

WORKDIR /app
COPY --from=probe /grpc_health_probe /grpc_health_probe
COPY . /app

# 本番用の依存関係をインストール
RUN pip install --no-cache-dir -r requirements.txt

ENV PYTHONUNBUFFERED=1
ENV PYTHONPATH=/app:/app/gen
USER 65532:65532

CMD ["python", "./cmd/server/main.py", "server"]
