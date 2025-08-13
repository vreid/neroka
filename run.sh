#!/usr/bin/env bash

docker build -t "vreid/neroka:latest" .
docker run \
    --env-file .env \
    "vreid/neroka:latest"
