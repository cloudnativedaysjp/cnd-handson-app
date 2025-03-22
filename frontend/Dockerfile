FROM node:22-alpine AS build


WORKDIR /frontend

RUN corepack enable && corepack prepare pnpm@latest --activate

COPY package.json pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

COPY . .

CMD ["/bin/sh"]
