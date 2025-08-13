FROM oven/bun:1-alpine AS build

WORKDIR /build

COPY bun.lock .
COPY package.json .

RUN bun install \
    --frozen-lockfile

COPY tsconfig.json .
COPY src/ src/

RUN bun build \
    --target=bun \
    --production \
    --outfile=neroka.js \
    src/index.ts

FROM oven/bun:1-alpine

WORKDIR /app

COPY bun.lock .
COPY package.json .

RUN bun install \
    --production \
    --frozen-lockfile

COPY --from=build /build/neroka.js .

ENTRYPOINT [ "bun", "run", "/app/neroka.js" ]