#!/usr/bin/env bash

set -ex -u -o pipefail

if [ -f .env ]; then
    set -a
    # shellcheck source=/dev/null
    source .env
    set +a
fi

curl https://api.anthropic.com/v1/models \
    --header "x-api-key: ${ANTHROPIC_API_KEY}" \
    --header "anthropic-version: 2023-06-01" | jq -r '.data[].id'
