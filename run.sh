#!/usr/bin/env bash

set -e -u -o pipefail

docker build -t "vreid/neroka:latest" .
docker run \
    --env-file .env \
    "vreid/neroka:latest" \
    -h
