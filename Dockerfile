FROM oven/bun:1-debian AS build

WORKDIR /build

COPY bun.lock .
COPY package.json .

RUN bun install \
    --frozen-lockfile

COPY tsconfig.json .
COPY src/ src/

RUN bun run build

FROM oven/bun:1-debian

WORKDIR /app

COPY bun.lock .
COPY package.json .

RUN apt-get update &&\
    apt-get install -y \
    ca-certificates

RUN bun install \
    --production \
    --frozen-lockfile

COPY --from=build /build/neroka.js .

ENTRYPOINT [ "bun", "run", "/app/neroka.js" ]