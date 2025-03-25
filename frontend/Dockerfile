FROM node:22-alpine

WORKDIR /app

RUN corepack enable && corepack prepare pnpm@latest --activate

ENV PNPM_STORE_PATH=/app/.pnpm-store

COPY package.json pnpm-lock.yaml ./
RUN pnpm install

COPY . .

CMD ["pnpm", "dev"]