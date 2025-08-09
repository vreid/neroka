#!/usr/bin/env bash

set -ex -u -o pipefail

if [ -f .env ]; then
  set -a
  # shellcheck source=/dev/null
  source .env
  set +a
fi

curl https://api.openai.com/v1/models \
  -H "Authorization: Bearer ${OPENAI_API_KEY}" |
  jq -r '.data[].id'
