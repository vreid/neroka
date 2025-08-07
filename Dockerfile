FROM golang:1.24-alpine AS build

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY internal/ internal/
COPY main.go .

ARG CGO_ENABLED=0
RUN go build -ldflags "-s -w"

FROM scratch

WORKDIR /app

COPY --from=build /build/neroka .

ENTRYPOINT [ "/app/neroka" ]