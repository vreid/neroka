#!/usr/bin/env bash

set -ex -u -o pipefail

if [ -f .env ]; then
    set -a
    # shellcheck source=/dev/null
    source .env
    set +a
fi

curl https://openrouter.ai/api/v1/models |
    jq -r '.data[].id'

curl -X POST https://openrouter.ai/api/v1/chat/completions \
    -H "Authorization: Bearer ${OPENROUTER_API_KEY}" \
    -H "Content-Type: application/json" \
    -d '{
  "model": "anthropic/claude-3.5-haiku-20241022",
  "messages": [
    {
      "role": "user",
      "content": "Say \"This is a test.\""
    }
  ]
}' | jq
