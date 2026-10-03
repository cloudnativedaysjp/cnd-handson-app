# Build context is the repo root (pnpm workspace): docker build -f frontend/Dockerfile .
FROM node:24.21.0-alpine3.23

WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY frontend/package.json ./frontend/
RUN pnpm install --frozen-lockfile --filter app

COPY frontend ./frontend
WORKDIR /app/frontend

CMD ["pnpm", "dev"]
