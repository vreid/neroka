#!/usr/bin/env bash

set -ex -u -o pipefail

if [ -f .env ]; then
    set -a
    # shellcheck source=/dev/null
    source .env
    set +a
fi

go build -ldflags "-s -w"

./neroka test-anthropic
./neroka test-openai
./neroka test-openrouter
